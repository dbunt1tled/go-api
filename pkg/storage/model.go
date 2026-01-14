package storage

import "github.com/google/uuid"

type Model interface {
	TableName() string
	GetID() uuid.UUID
	SetID(uuid.UUID)
}
