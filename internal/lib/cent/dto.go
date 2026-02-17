package cent

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	UserID     uuid.UUID `json:"userId"`
	FromUserID uuid.UUID `json:"fromUserId"`
	Message    string    `json:"message"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt"`
	Status     string    `json:"status"`
}
