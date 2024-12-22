package user

import (
	"database/sql/driver"
	"fmt"

	"github.com/pkg/errors"
)

type Status int

const (
	Active  Status = 1
	Delete  Status = 2
	Pending Status = 3
)

func (s *Status) String() string {
	switch *s {
	case Active:
		return "Active"
	case Delete:
		return "Inactive"
	case Pending:
		return "Pending"
	default:
		panic(fmt.Errorf("unknown user status: %d", s))
	}
}

func (s Status) Val() int {
	return FromStatus(s)
}

func FromStatus(s Status) int {
	switch s {
	case Active:
		return 1 //nolint:nolintlint
	case Delete:
		return 2 //nolint:nolintlint
	case Pending:
		return 3 //nolint:nolintlint
	default:
		panic(fmt.Errorf("unknown user status: %d", s))
	}
}

func FromInt(value int) (Status, error) {
	switch value {
	case 1:
		return Active, nil
	case 2:
		return Delete, nil
	case 3:
		return Pending, nil
	default:
		panic(fmt.Errorf("unknown user status: %d", value))
	}
}

func (s *Status) Scan(value interface{}) error {
	val, ok := value.(int64)
	if !ok {
		return errors.New("failed to scan status: not an int64")
	}

	status, err := FromInt(int(val))
	if err != nil {
		return err
	}

	*s = status
	return nil
}

func (s *Status) Value() (driver.Value, error) {
	switch *s {
	case Active:
		return 1, nil //nolint:nolintlint
	case Delete:
		return 2, nil //nolint:nolintlint
	case Pending:
		return 3, nil //nolint:nolintlint
	default:
		return nil, errors.New("invalid status value")
	}
}
