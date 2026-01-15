package user

import (
	"github.com/dbunt1tled/go-api/pkg/http/dto"
	"github.com/dbunt1tled/go-api/pkg/storage"
)

func NewUserResource(u *User) *dto.Resource {
	resource := dto.NewResource("user", u.ID.String())
	resource.MarshalAttributes(u)
	return resource
}

func NewUserResponse(u *User) *dto.Document {
	return dto.NewResponse().SetData(NewUserResource(u)).Build()
}

func NewUserListResponse(u *storage.Paginator[*User]) *dto.Document {
	resources := make([]*dto.Resource, len(u.Items))
	for i, user := range u.Items {
		resources[i] = NewUserResource(user)
	}
	return dto.NewResponse().SetData(resources).SetMetaPagination(u).Build()
}
