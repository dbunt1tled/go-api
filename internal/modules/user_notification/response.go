package user_notification

import (
	"github.com/dbunt1tled/go-api/pkg/http/dto"
	"github.com/dbunt1tled/go-api/pkg/storage"
)

func NewUserNotificationResource(u *UserNotification) *dto.Resource {
	resource := dto.NewResource("userNotification", u.ID.String())
	resource.MarshalAttributes(u)
	return resource
}

func NewUserNotificationResponse(u *UserNotification) *dto.Document {
	return dto.NewResponse().SetData(NewUserNotificationResource(u)).Build()
}

func NewUserNotificationListResponse(u *storage.Paginator[*UserNotification]) *dto.Document {
	resources := make([]*dto.Resource, len(u.Items))
	for i, user := range u.Items {
		resources[i] = NewUserNotificationResource(user)
	}
	return dto.NewResponse().SetData(resources).SetMetaPagination(u).Build()
}
