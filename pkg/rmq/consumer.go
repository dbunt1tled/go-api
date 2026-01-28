package rmq

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dbunt1tled/go-api/internal/jobs"
	"github.com/dbunt1tled/go-api/pkg/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const Delay = 100 * time.Millisecond

type ResilientConnection struct {
	url           string
	conn          *amqp.Connection
	channel       *amqp.Channel
	mu            sync.RWMutex
	reconnecting  sync.Mutex
	done          chan struct{}
	reconnectDone chan struct{}
}

func NewResilientConnection(url string) (*ResilientConnection, error) {
	rc := &ResilientConnection{
		url:           url,
		done:          make(chan struct{}),
		reconnectDone: make(chan struct{}),
	}

	if err := rc.connect(); err != nil {
		return nil, err
	}

	close(rc.reconnectDone)
	return rc, nil
}

func (rc *ResilientConnection) connect() error {
	conn, err := amqp.Dial(rc.url)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	if err := ch.Qos(10, 0, false); err != nil { //nolint:nolintlint,mnd
		_ = ch.Close()
		_ = conn.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	rc.mu.Lock()
	rc.conn = conn
	rc.channel = ch
	rc.mu.Unlock()

	go rc.monitorConnection()

	log.Logger().Info("Consumer connection established")
	return nil
}

func (rc *ResilientConnection) monitorConnection() {
	rc.mu.RLock()
	conn := rc.conn
	rc.mu.RUnlock()

	errChan := make(chan *amqp.Error)
	conn.NotifyClose(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			log.Logger().Infof("Connection closed: %v. Reconnecting...", err)
			rc.reconnect()
		}
	case <-rc.done:
		return
	}
}

func (rc *ResilientConnection) reconnect() {
	rc.reconnecting.Lock()
	defer rc.reconnecting.Unlock()

	rc.reconnectDone = make(chan struct{})
	defer close(rc.reconnectDone)

	rc.mu.Lock()
	if rc.channel != nil {
		_ = rc.channel.Close()
	}
	if rc.conn != nil {
		_ = rc.conn.Close()
	}
	rc.conn = nil
	rc.channel = nil
	rc.mu.Unlock()

	backoff := time.Second
	maxBackoff := 30 * time.Second //nolint:nolintlint,mnd

	for {
		select {
		case <-rc.done:
			log.Logger().Info("Reconnection cancelled due to shutdown")
			return
		case <-time.After(backoff):
			log.Logger().Infof("Reconnection attempt...")

			if err := rc.connect(); err != nil {
				log.Logger().Infof("Reconnection failed: %v. Retrying in %v", err, backoff)
				backoff = min(backoff*2, maxBackoff) //nolint:nolintlint,mnd
				continue
			}

			log.Logger().Info("Reconnection successful")
			return
		}
	}
}

func (rc *ResilientConnection) GetChannel() (*amqp.Channel, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if rc.channel == nil || rc.channel.IsClosed() {
		return nil, fmt.Errorf("channel not available")
	}

	return rc.channel, nil
}

func (rc *ResilientConnection) WaitForConnection(ctx context.Context) error {
	select {
	case <-rc.reconnectDone:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-rc.done:
		return fmt.Errorf("connection closed")
	}
}

func (rc *ResilientConnection) Close() error {
	close(rc.done)

	rc.mu.Lock()
	defer rc.mu.Unlock()

	var errs []error
	if rc.channel != nil {
		if err := rc.channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if rc.conn != nil {
		if err := rc.conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connection: %v", errs)
	}

	log.Logger().Info("Consumer connection closed")
	return nil
}

type MessageHandler func(context.Context, []byte) error

type Consumer struct {
	rc         *ResilientConnection
	queueName  string
	retryQueue string
	dlqQueue   string
	maxRetries int
	retryDelay time.Duration
	done       chan struct{}
	wg         sync.WaitGroup
}

type ConsumerConfig struct {
	QueueName    string
	MaxRetries   int
	RetryDelay   time.Duration
	QueueType    string // "quorum", "classic", or "" to use existing
	Durable      bool
	AutoDeclare  bool // Set to false to skip queue declaration
	ExchangeName string
	ExchangeType string
	RoutingKeys  []string
}

func NewConsumer(url string, config ConsumerConfig) (*Consumer, error) {
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay <= 0 {
		config.RetryDelay = 30 * time.Second //nolint:nolintlint,mnd
	}
	if config.QueueType == "" {
		config.QueueType = "classic"
	}

	rc, err := NewResilientConnection(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	consumer := &Consumer{
		rc:         rc,
		queueName:  config.QueueName,
		retryQueue: config.QueueName + "_retry",
		dlqQueue:   config.QueueName + "_dlq",
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryDelay,
		done:       make(chan struct{}),
	}

	// Only declare queues if AutoDeclare is true (default behavior)
	if config.AutoDeclare {
		if err := consumer.setupTopology(config); err != nil {
			_ = rc.Close()
			return nil, err
		}
	}

	log.Logger().Infof("Consumer initialized for queue: %s (max retries: %d, type: %s)",
		config.QueueName, config.MaxRetries, config.QueueType)
	return consumer, nil
}

func (c *Consumer) setupTopology(config ConsumerConfig) error {
	ch, err := c.rc.GetChannel()
	if err != nil {
		return err
	}

	args := amqp.Table{}
	if config.QueueType == "quorum" {
		args["x-queue-type"] = "quorum"
	}

	// Declare exchange if configured
	if config.ExchangeName != "" {
		if config.ExchangeType == "" {
			config.ExchangeType = "direct"
		}

		err = ch.ExchangeDeclare(
			config.ExchangeName,
			config.ExchangeType,
			config.Durable,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange: %w", err)
		}
		log.Logger().Infof("Exchange declared: %s (type: %s)", config.ExchangeName, config.ExchangeType)
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(
		c.dlqQueue,
		config.Durable,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Declare retry queue with TTL
	retryArgs := amqp.Table{
		"x-message-ttl":             int32(c.retryDelay.Milliseconds()),
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": c.queueName,
	}
	if config.QueueType == "quorum" {
		retryArgs["x-queue-type"] = "quorum"
	}

	_, err = ch.QueueDeclare(
		c.retryQueue,
		config.Durable,
		false,
		false,
		false,
		retryArgs,
	)
	if err != nil {
		return fmt.Errorf("failed to declare retry queue: %w", err)
	}

	// Declare main queue (or use existing)
	_, err = ch.QueueDeclare(
		c.queueName,
		config.Durable,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		// If queue exists with different parameters, try passive declare
		_, err = ch.QueueDeclarePassive(
			c.queueName,
			config.Durable,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare main queue: %w", err)
		}
		log.Logger().Infof("Using existing queue: %s", c.queueName)
	}

	// Bind queue to exchange if configured
	if config.ExchangeName != "" {
		if len(config.RoutingKeys) == 0 {
			config.RoutingKeys = []string{""} // Empty routing key for fanout
		}

		for _, bindingKey := range config.RoutingKeys {
			err = ch.QueueBind(
				c.queueName,
				bindingKey,
				config.ExchangeName,
				false,
				nil,
			)
			if err != nil {
				return fmt.Errorf("failed to bind queue to exchange: %w", err)
			}
			log.Logger().Infof("Queue bound: %s -> %s (routing key: %s)", c.queueName, config.ExchangeName, bindingKey)
		}
	}

	log.Logger().
		Infof("Topology setup complete: %s, %s, %s (type: %s)", c.queueName, c.retryQueue, c.dlqQueue, config.QueueType)
	return nil
}

func (c *Consumer) Start(
	ctx context.Context,
	resolver jobs.RMQResolver,
	workers int,
	sleep *time.Duration,
) error {
	if workers <= 0 {
		workers = 1
	}

	delay := Delay
	if sleep != nil {
		delay = *sleep
	}

	log.Logger().Infof("Starting %d consumer workers", workers)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 1; i <= workers; i++ {
		c.wg.Add(1)
		go c.worker(workerCtx, resolver, i, delay)
	}

	select {
	case sig := <-sigChan:
		log.Logger().Infof("Received signal: %v. Shutting down gracefully...", sig)
	case <-ctx.Done():
		log.Logger().Infof("Context cancelled. Shutting down gracefully...")
	}

	cancel()
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Logger().Warn("All consumer workers stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Logger().Warn("Graceful shutdown timeout exceeded")
		return fmt.Errorf("shutdown timeout exceeded")
	}

	return nil
}

func (c *Consumer) worker(ctx context.Context, resolver jobs.RMQResolver, workerID int, delay time.Duration) {
	defer c.wg.Done()

	log.Logger().Infof("Worker %d: Started", workerID)

	for {
		select {
		case <-ctx.Done():
			log.Logger().Infof("Worker %d: Context cancelled", workerID)
			return
		case <-c.done:
			log.Logger().Infof("Worker %d: Shutdown requested", workerID)
			return
		default:
			if err := c.consume(ctx, resolver, workerID, delay); err != nil {
				log.Logger().Infof("Worker %d: Consumer error: %v. Restarting...", workerID, err)
				time.Sleep(time.Second)
			}
		}
	}
}

func (c *Consumer) consume(ctx context.Context, resolver jobs.RMQResolver, workerID int, delay time.Duration) error {
	if err := c.rc.WaitForConnection(ctx); err != nil {
		return err
	}

	ch, err := c.rc.GetChannel()
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queueName,
		fmt.Sprintf("worker-%d-%d", workerID, time.Now().Unix()),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	closeChan := make(chan *amqp.Error)
	ch.NotifyClose(closeChan)

	log.Logger().Infof("Worker %d: Consuming messages from %s", workerID, c.queueName)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.done:
			return nil
		case err := <-closeChan:
			return fmt.Errorf("channel closed: %w", err)
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("message channel closed")
			}
			c.handleMessage(ctx, msg, resolver, workerID, delay)
		}
	}
}

func (c *Consumer) handleMessage(
	ctx context.Context,
	msg amqp.Delivery,
	resolver jobs.RMQResolver,
	workerID int,
	delay time.Duration,
) {
	time.Sleep(delay)
	retryCount := c.getRetryCount(msg.Headers)

	log.Logger().Infof("Worker %d: Processing message (attempt %d/%d): %s",
		workerID, retryCount+1, c.maxRetries+1, string(msg.Body))
	handler, err := resolver.Resolve(msg.Type)
	if err != nil {
		log.Logger().Error(
			fmt.Sprintf("Worker %d: Error processing message: %s", workerID, err.Error()),
			err,
			slog.Any("msg", msg),
		)
		return
	}
	if handler == nil {
		log.Logger().
			Warn(fmt.Sprintf("Worker %d: Error processing message: handler is empty", workerID), slog.Any("msg", msg))
		return
	}
	if err = (*handler).Handle(ctx, msg.Body); err != nil { //nolint:nolintlint,nestif
		log.Logger().
			Error(fmt.Sprintf("Worker %d: Error handling message: %v", workerID, err.Error()), err, slog.Any("msg", msg))

		if retryCount < c.maxRetries {
			c.retryMessage(msg, retryCount+1)
			err = msg.Ack(false)
			if err != nil {
				log.Logger().
					Error(fmt.Sprintf("Worker %d: Failed to ack message: %v", workerID, err.Error()), err, slog.Any("msg", msg))
			}
			log.Logger().Infof("Worker %d: Message sent to retry queue (attempt %d/%d)",
				workerID, retryCount+1, c.maxRetries)
		} else {
			log.Logger().Infof("Worker %d: Max retries exceeded, sending to DLQ", workerID)
			c.sendToDLQ(msg, err)
			err := msg.Ack(false)
			if err != nil {
				log.Logger().Error(fmt.Sprintf("Worker %d: Failed to ack message max retries exceeded: %v", workerID, err.Error()), err, slog.Any("msg", msg))
			}
		}
		return
	}

	if err = msg.Ack(false); err != nil {
		log.Logger().
			Error(fmt.Sprintf("Worker %d: Failed to ack message: %v", workerID, err.Error()), err, slog.Any("msg", msg))
	} else {
		log.Logger().Infof("Worker %d: Message processed successfully", workerID)
	}
}

func (c *Consumer) retryMessage(msg amqp.Delivery, retryCount int) {
	ch, err := c.rc.GetChannel()
	if err != nil {
		log.Logger().Infof("Failed to get channel for retry: %v", err)
		return
	}

	headers := make(amqp.Table)
	if msg.Headers != nil {
		headers = msg.Headers
	}
	headers["x-retry-count"] = int32(retryCount)
	headers["x-first-death-time"] = time.Now().Format(time.RFC3339)
	headers["ex-message-id"] = msg.MessageId

	err = ch.PublishWithContext(
		context.Background(),
		"",
		c.retryQueue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			Headers:      headers,
			Type:         msg.Type,
			UserId:       msg.UserId,
			AppId:        msg.AppId,
			MessageId:    uuid.New().String(),
		},
	)
	if err != nil {
		log.Logger().Infof("Failed to send message to retry queue: %v", err)
	}
}

func (c *Consumer) sendToDLQ(msg amqp.Delivery, processingErr error) {
	ch, err := c.rc.GetChannel()
	if err != nil {
		log.Logger().Infof("Failed to get channel for DLQ: %v", err)
		return
	}

	headers := make(amqp.Table)
	if msg.Headers != nil {
		headers = msg.Headers
	}
	headers["x-error"] = processingErr.Error()
	headers["x-failed-at"] = time.Now().Format(time.RFC3339)
	headers["x-original-queue"] = c.queueName
	headers["ex-message-id"] = msg.MessageId

	err = ch.PublishWithContext(
		context.Background(),
		"",
		c.dlqQueue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			Headers:      headers,
			Type:         msg.Type,
			UserId:       msg.UserId,
			AppId:        msg.AppId,
			MessageId:    uuid.New().String(),
		},
	)
	if err != nil {
		log.Logger().Infof("Failed to send message to DLQ: %v", err)
	}
}

func (c *Consumer) getRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}
	if count, ok := headers["x-retry-count"].(int32); ok {
		return int(count)
	}
	return 0
}

func (c *Consumer) Close() error {
	close(c.done)
	c.wg.Wait()
	return c.rc.Close()
}
