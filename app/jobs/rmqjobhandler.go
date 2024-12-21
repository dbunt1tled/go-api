package jobs

import "fmt"

type RMQJobHandler interface {
	Handle(body []byte) error
}

type RMQJobResolver struct {
	handlers map[string]*RMQJobHandler
}

func NewRMQJobResolver() *RMQJobResolver {
	return &RMQJobResolver{
		handlers: make(map[string]*RMQJobHandler),
	}
}

func (r *RMQJobResolver) RegisterHandler(jobName string, handler RMQJobHandler) {
	r.handlers[jobName] = &handler
}

func (r *RMQJobResolver) Resolve(jobName string) (*RMQJobHandler, error) {
	handler, exists := r.handlers[jobName]
	if !exists {
		return nil, fmt.Errorf("job handler for %s not found", jobName)
	}
	return handler, nil
}
