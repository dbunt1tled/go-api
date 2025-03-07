package dbtype

import (
	"database/sql"
	"database/sql/driver"

	"github.com/bytedance/sonic"
)

type NilTime sql.NullTime //nolint:recvcheck

func (nt NilTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return sonic.ConfigFastest.Marshal(nt.Time)
}

func (nt *NilTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}
	nt.Valid = true
	return sonic.ConfigFastest.Unmarshal(data, &nt.Time)
}

func (nt *NilTime) Scan(value any) error {
	return (*sql.NullTime)(nt).Scan(value)
}

func (nt NilTime) Value() (driver.Value, error) {
	return (*sql.NullTime)(&nt).Value()
}
