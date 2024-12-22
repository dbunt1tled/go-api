package usernotification

import (
	"go_echo/internal/util/type/json"
	"go_echo/internal/util/type/timestamp"
)

type UserNotification struct {
	ID        int64               `json:"id" jsonapi:"primary,userNotification"`
	UserId    int64               `json:"userId" jsonapi:"attr,userId"`
	Data      json.JsonField      `json:"data" jsonapi:"attr,data"`
	Status    Status              `json:"status" jsonapi:"attr,status"`
	CreatedAt timestamp.Timestamp `json:"created_at" jsonapi:"attr,createdAt"`
	UpdatedAt timestamp.Timestamp `json:"updated_at" jsonapi:"attr,updatedAt"`
}

func (us *UserNotification) Sanitize() {
	// Sanitize the user notification
}
