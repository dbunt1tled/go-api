package user

import (
	"database/sql"
	"go_echo/internal/util/type/roles"
	"go_echo/internal/util/type/timestamp"
	"strings"
)

type User struct {
	ID          int64               `json:"id" jsonapi:"primary,url"`
	FirstName   string              `json:"firstName" jsonapi:"attr,firstName"`
	SecondName  string              `json:"secondName" jsonapi:"attr,secondName"`
	Email       string              `json:"email" jsonapi:"attr,email"`
	PhoneNumber string              `json:"phoneNumber" jsonapi:"attr,phoneNumber"`
	Password    string              `json:"password" jsonapi:"attr,password"`
	Status      int                 `json:"status" jsonapi:"attr,status"`
	Hash        sql.NullString      `json:"hash" jsonapi:"attr,hash"`
	Roles       roles.Roles         `json:"roles" jsonapi:"attr,roles"`
	ConfirmedAt sql.NullTime        `json:"confirmed_at" jsonapi:"attr,confirmedAt"`
	CreatedAt   timestamp.Timestamp `json:"created_at" jsonapi:"attr,createdAt"`
	UpdatedAt   timestamp.Timestamp `json:"updated_at" jsonapi:"attr,updatedAt"`
}

func (us *User) Sanitize() {
	us.FirstName = strings.TrimSpace(us.FirstName)
	us.Email = strings.TrimSpace(us.Email)
	us.PhoneNumber = strings.TrimSpace(us.PhoneNumber)
	us.Password = strings.TrimSpace(us.Password)
	if us.Hash.Valid {
		us.Hash.String = strings.TrimSpace(us.Hash.String)
	}
}
