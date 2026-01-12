package dto

type SetDefaults interface {
	SetDefaults()
}

type PaginationQuery struct {
	Page *PageQuery `json:"page" query:"page" validate:"omitempty"`
	Sort *Sorting   `json:"sort" query:"sort" validate:"omitempty"`
}

type PageQuery struct {
	Page  int `json:"page"  query:"page"  validate:"required,numeric"`
	Limit int `json:"limit" query:"limit" validate:"required,numeric"`
}

type Sorting struct {
	Field string `json:"field" query:"field" validate:"required,min=2"`
	Order string `json:"order" query:"order" validate:"required,oneof=asc desc"`
}

func (r *PaginationQuery) SetDefaults() {
	if r.Page == nil {
		r.Page = &PageQuery{
			Page:  1,
			Limit: 20, //nolint:mnd // default limit
		}
	} else {
		if r.Page.Page == 0 {
			r.Page.Page = 1
		}
		if r.Page.Limit == 0 {
			r.Page.Limit = 20 //nolint:nolintlint,mnd // default limit
		}
	}

	if r.Sort == nil {
		field := "id"
		order := "desc"
		r.Sort = &Sorting{
			Field: field,
			Order: order,
		}
	} else {
		if r.Sort.Field == "" {
			f := "id"
			r.Sort.Field = f
		}
		if r.Sort.Order == "" {
			o := "desc"
			r.Sort.Order = o
		}
	}
}
