package user_notification

import "fmt"

type Status int

const (
	New  Status = 1
	Read Status = 2
)

func (s *Status) String() string {
	switch *s {
	case New:
		return "New"
	case Read:
		return "Read"
	default:
		panic(fmt.Errorf("unknown user status: %d", s))
	}
}
