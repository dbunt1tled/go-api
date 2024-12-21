package rmqmail

import (
	"go_echo/app/jobs"
	"go_echo/app/jobs/rmqmail/handlers"
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
