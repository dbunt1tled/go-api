package rmqmail

import (
	"time"

	"github.com/dbunt1tled/go-api/internal/jobs"
	"github.com/dbunt1tled/go-api/internal/jobs/rmqmail/handlers"
)

const (
	MailExchange   = "mail-go-exchange"
	MailQueue      = "mail-go-queue"
	MailRoutingKey = "mail"
	MaxRetries     = 3
	RetryDelay     = 30 * time.Second
)

type RMQJobMailResolver struct {
	resolver                 *jobs.RMQJobResolver
	UserConfirmationEmailJob *handlers.UserConfirmationEmailJob
}

func NewRMQJobMailResolver(
	userConfirmationEmailJob *handlers.UserConfirmationEmailJob,
) *RMQJobMailResolver {
	r := &RMQJobMailResolver{
		resolver:                 jobs.NewRMQJobResolver(),
		UserConfirmationEmailJob: userConfirmationEmailJob,
	}

	r.resolver.RegisterHandler(r.UserConfirmationEmailJob.Action(), r.UserConfirmationEmailJob)

	return r
}

func (r *RMQJobMailResolver) Resolve(jobName string) (*jobs.RMQJobHandler, error) {
	return r.resolver.Resolve(jobName)
}
