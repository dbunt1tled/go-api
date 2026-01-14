package user

type Status int

const (
	Pending Status = 0
	Active  Status = 1
	Blocked Status = 2
)

func (s Status) String() string {
	switch s {
	case Active:
		return "active"
	case Blocked:
		return "blocked"
	case Pending:
		return "pending"
	default:
		return "inactive"
	}
}
