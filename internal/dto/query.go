package dto

type PaginationQuery struct {
	Page  *int    `query:"page"`
	Limit *int    `query:"limit"`
	Sort  Sorting `json:"sort" validate:"required"`
}

type Sorting struct {
	Field *string `query:"field" validate:"omitempty"`
	Order *string `query:"order" validate:"omitempty,oneof=asc desc"`
}
