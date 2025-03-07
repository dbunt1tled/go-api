package roles

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

type Roles []string

func (c Roles) Value() (driver.Value, error) {
	if len(c) == 0 {
		return nil, nil
	}
	j, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return driver.Value([]byte(j)), nil
}

// func (c *Roles) Scan(src interface{}) (err error) {
// 	var rl []string
// 	switch src.(type) {
// 	case string:
// 		err = json.Unmarshal([]byte(src.(string)), &rl)
// 	case []byte:
// 		err = json.Unmarshal(src.([]byte), &rl)
// 	// case []uint8:
// 	// 	err = json.Unmarshal([]byte(src.([]uint8)), &rl)
// 	default:
// 		return errors.New("unsupported type")
// 	}
// 	if err != nil {
// 		return
// 	}
// 	*c = rl
// 	return nil
// }

func (c *Roles) Scan(v interface{}) error {
	var (
		_rl Roles
		err error
	)
	switch tv := v.(type) {
	case []byte:
		err = sonic.ConfigFastest.Unmarshal(tv, &_rl)
		// case []uint8:
		// 	err = sonic.ConfigFastest.Unmarshal([]byte(tv), &_rl)
	}
	if err != nil {
		errors.Wrap(err, "Error roles")
	}
	*c = _rl
	return nil
}
