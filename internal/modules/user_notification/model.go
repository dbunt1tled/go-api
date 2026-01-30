package user_notification

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type UserNotification struct {
	bun.BaseModel `bun:"table:user_notifications"`

	ID        uuid.UUID      `bun:"id,pk,type:uuid,default:uuid_v7()" json:"id"`
	UserID    uuid.UUID      `bun:"user_id,notnull,type:uuid"         json:"userId"`
	Data      map[string]any `bun:"data,type:jsonb"                   json:"data"`
	Status    Status         `bun:"status,notnull,default:1"          json:"status"`
	CreatedAt time.Time      `bun:"created_at,notnull,default:now()"  json:"createdAt"`
	UpdatedAt time.Time      `bun:"updated_at,notnull,default:now()"  json:"updatedAt"`
}

func (u *UserNotification) TableName() string  { return "users" }
func (u *UserNotification) GetID() uuid.UUID   { return u.ID }
func (u *UserNotification) SetID(id uuid.UUID) { u.ID = id }
func (u *UserNotification) NextID() *UserNotification {
	var err error
	u.ID, err = uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("failed to generate uuid: %w", err))
	}
	return u
}
