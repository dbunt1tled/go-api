package json

import (
	"database/sql/driver"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type JsonField map[string]interface{}

func (m *JsonField) Scan(src interface{}) error {
	var source []byte
	_m := make(map[string]interface{})
	switch src.(type) {
	case []uint8:
		source = []byte(src.([]uint8))
	case nil:
		return nil
	default:
		return errors.New("incompatible type for JsonField")
	}
	err := sonic.Unmarshal(source, &_m)
	if err != nil {
		return err
	}
	*m = JsonField(_m)
	return nil
}

func (m JsonField) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil //nolint:nilnil
	}
	j, err := sonic.Marshal(m)
	if err != nil {
		return nil, err
	}
	return driver.Value([]byte(j)), nil
}
