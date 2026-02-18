package cent

import (
	"time"

	"github.com/google/uuid"
)

type UserMessage struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
