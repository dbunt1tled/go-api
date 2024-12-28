package dbtype

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

type NilString sql.NullString //nolint:recvcheck

func (ns NilString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NilString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}

func (ns *NilString) Scan(value any) error {
	return (*sql.NullString)(ns).Scan(value)
}

func (ns NilString) Value() (driver.Value, error) {
	return (*sql.NullString)(&ns).Value()
}
