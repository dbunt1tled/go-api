package rmqmail

import (
	"github.com/dbunt1tled/go-api/app/jobs"
	"github.com/dbunt1tled/go-api/app/jobs/rmqmail/handlers"
)

type RMQJobMailResolver struct {
	Resolver *jobs.RMQJobResolver
}

func NewRMQJobMailResolver() *RMQJobMailResolver {
	resolver := &RMQJobMailResolver{
		Resolver: jobs.NewRMQJobResolver(),
	}

	resolver.Resolver.RegisterHandler(handlers.ConfirmSubject, handlers.UserConfirmationEmail{})

	return resolver
}
