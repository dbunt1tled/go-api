package rmq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const JsonHeader = "application/json"

type Producer struct {
	url           string
	conn          *amqp.Connection
	channelPool   chan *amqp.Channel
	maxChannels   int
	mu            sync.RWMutex
	reconnecting  sync.Mutex
	done          chan struct{}
	reconnectDone chan struct{}
}

type PublishOpt struct {
	UserId  string
	AppId   string
	Headers amqp.Table
}

func getDefaultPublisherOptions() *PublishOpt {
	return &PublishOpt{
		UserId:  "",
		AppId:   "",
		Headers: amqp.Table{},
	}
}

func NewProducer(url string, maxChannels int) (*Producer, error) {
	if maxChannels <= 0 {
		maxChannels = 5
	}

	p := &Producer{
		url:           url,
		maxChannels:   maxChannels,
		channelPool:   make(chan *amqp.Channel, maxChannels),
		done:          make(chan struct{}),
		reconnectDone: make(chan struct{}),
	}

	if err := p.connect(); err != nil {
		return nil, err
	}

	close(p.reconnectDone)
	log.Logger().Info("Producer initialized successfully")

	return p, nil
}

func (p *Producer) connect() error {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	p.mu.Lock()
	p.conn = conn
	p.mu.Unlock()

	// Create channel pool
	for i := 0; i < p.maxChannels; i++ {
		ch, err := p.createChannel()
		if err != nil {
			p.closeAllChannels()
			_ = conn.Close()
			return fmt.Errorf("failed to create channel %d: %w", i, err)
		}
		p.channelPool <- ch
	}

	go p.monitorConnection()

	log.Logger().Info("Connection established")
	return nil
}

func (p *Producer) createChannel() (*amqp.Channel, error) {
	p.mu.RLock()
	conn := p.conn
	p.mu.RUnlock()

	if conn == nil {
		return nil, fmt.Errorf("no connection available")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("failed to enable confirms: %w", err)
	}

	return ch, nil
}

func (p *Producer) monitorConnection() {
	p.mu.RLock()
	conn := p.conn
	p.mu.RUnlock()

	errChan := make(chan *amqp.Error)
	conn.NotifyClose(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			log.Logger().Infof("Connection closed: %v. Initiating reconnection...", err)
			p.reconnect()
		}
	case <-p.done:
		return
	}
}

func (p *Producer) reconnect() {
	p.reconnecting.Lock()
	defer p.reconnecting.Unlock()

	p.reconnectDone = make(chan struct{})
	defer close(p.reconnectDone)

	p.closeAllChannels()

	p.mu.Lock()
	if p.conn != nil {
		_ = p.conn.Close()
		p.conn = nil
	}
	p.mu.Unlock()

	backoff := time.Second
	maxBackoff := 30 * time.Second //nolint:nolintlint,mnd

	for {
		select {
		case <-p.done:
			log.Logger().Info("Reconnection cancelled due to shutdown")
			return
		case <-time.After(backoff):
			log.Logger().Info("Reconnection attempt...")

			if err := p.connect(); err != nil {
				log.Logger().Info("Reconnection failed: %v. Retrying in %v", err, backoff)
				backoff = minDuration(backoff*2, maxBackoff) //nolint:nolintlint,mnd
				continue
			}

			log.Logger().Info("Reconnection successful")
			return
		}
	}
}

func (p *Producer) closeAllChannels() {
	for {
		select {
		case ch := <-p.channelPool:
			if ch != nil && !ch.IsClosed() {
				_ = ch.Close()
			}
		default:
			return
		}
	}
}

func (p *Producer) getChannel(ctx context.Context) (*amqp.Channel, error) {
	// Wait for reconnection if in progress
	select {
	case <-p.reconnectDone:
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.done:
		return nil, fmt.Errorf("producer closed")
	}

	select {
	case ch := <-p.channelPool:
		if ch != nil && !ch.IsClosed() {
			return ch, nil
		}
		return p.getChannel(ctx)
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.done:
		return nil, fmt.Errorf("producer closed")
	}
}

func (p *Producer) returnChannel(ch *amqp.Channel) {
	if ch != nil && !ch.IsClosed() {
		select {
		case p.channelPool <- ch:
		case <-p.done:
			_ = ch.Close()
		default:
			_ = ch.Close()
		}
	}
}

func (p *Producer) DeclareQueue(ctx context.Context, queueName string, durable bool, queueType string) error {
	ch, err := p.getChannel(ctx)
	if err != nil {
		return err
	}
	defer p.returnChannel(ch)

	args := amqp.Table{}
	if queueType == "quorum" {
		args["x-queue-type"] = "quorum"
	}

	_, err = ch.QueueDeclare(
		queueName,
		durable,
		false,
		false,
		false,
		args,
	)

	// If queue exists with different parameters, try passive declare
	if err != nil {
		_, passiveErr := ch.QueueDeclarePassive(
			queueName,
			durable,
			false,
			false,
			false,
			nil,
		)
		if passiveErr != nil {
			return fmt.Errorf("failed to declare queue: %w (original: %v)", passiveErr, err)
		}
		log.Logger().Infof("Using existing queue: %s", queueName)
		return nil
	}

	log.Logger().Infof("Queue declared: %s (type: %s)", queueName, queueType)
	return nil
}

func (p *Producer) Publish(
	ctx context.Context,
	exchange string,
	routingKey string,
	action string,
	body []byte,
	optionFuncs ...func(*PublishOpt),
) error {
	ch, err := p.getChannel(ctx)
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}
	defer p.returnChannel(ch)
	options := getDefaultPublisherOptions()
	for _, optionFunc := range optionFuncs {
		optionFunc(options)
	}

	err = ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  JsonHeader,
			Body:         body,
			Type:         action,
			Timestamp:    time.Now(),
			MessageId:    uuid.New().String(),
			UserId:       options.UserId,
			AppId:        options.AppId,
			Headers:      options.Headers,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	close(p.done)

	p.closeAllChannels()

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conn != nil {
		err := p.conn.Close()
		log.Logger().Info("Producer closed")
		return err
	}

	return nil
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
func WithPublishAppID(appId string) func(*PublishOpt) {
	return func(options *PublishOpt) {
		options.AppId = appId
	}
}
func WithPublishUserId(userId string) func(*PublishOpt) {
	return func(options *PublishOpt) {
		options.UserId = userId
	}
}
