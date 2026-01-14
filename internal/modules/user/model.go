package user

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID          uuid.UUID      `bun:"id,pk,type:uuid,default:uuid_v7()"      json:"id"`
	FirstName   string         `bun:"first_name,notnull"                     json:"firstName"`
	SecondName  string         `bun:"second_name,notnull"                    json:"secondName"`
	Email       string         `bun:"email,notnull,unique"                   json:"email"`
	PhoneNumber string         `bun:"phone_number,notnull,unique"            json:"phoneNumber"`
	Status      Status         `bun:"status,notnull,default:0"               json:"status"`
	Password    string         `bun:"password,notnull"                       json:"-"`
	Roles       []string       `bun:"roles,type:text[],notnull,default:'{}'" json:"roles"`
	Address     map[string]any `bun:"address,type:jsonb"                     json:"address"`
	ConfirmedAt *time.Time     `bun:"confirmed_at,nullzero"                  json:"confirmedAt,omitempty"`
	CreatedAt   time.Time      `bun:"created_at,notnull,default:now()"       json:"createdAt"`
	UpdatedAt   time.Time      `bun:"updated_at,notnull,default:now()"       json:"updatedAt"`
}
func (u *User) TableName() string  { return "users" }
func (u *User) GetID() uuid.UUID   { return u.ID }
func (u *User) SetID(id uuid.UUID) { u.ID = id }
func (u *User) NextID() *User {
	var err error
	u.ID, err = uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("failed to generate uuid: %w", err))
	}
	return u
}

func (u *User) WithPassword(password string) *User {
	u.Password = password
	return u
}

func (u *User) Sanitize() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.SecondName = strings.TrimSpace(u.SecondName)
	u.Email = strings.TrimSpace(u.Email)
	u.PhoneNumber = strings.TrimSpace(u.PhoneNumber)
	u.Password = strings.TrimSpace(u.Password)
}

