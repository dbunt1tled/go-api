package user_notification

import "github.com/dbunt1tled/go-api/pkg/http/dto"

type ListRequest struct {
	dto.PaginationQuery

	Status *Status `query:"status" json:"status" validate:"omitempty" example:"0,1"`
}
