package roles

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type Roles []string

func (c Roles) Value() (driver.Value, error) {
	if len(c) == 0 {
		return "[]", nil
	}
	return json.Marshal(c) // return json marshalled value
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
	switch tv := v.(type) {
	case []byte:
		return json.Unmarshal(tv, &c) // unmarshal
		// case []uint8:
		// 	return json.Unmarshal([]byte(tv), &c) // can't remember the specifics, but this may be needed
	}
	return errors.New("unsupported type")
}
