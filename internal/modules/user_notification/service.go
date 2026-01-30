package user_notification

import (
	"context"

	"github.com/dbunt1tled/go-api/pkg/storage"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Service struct {
	userNotificationRepository *storage.Repository[UserNotification]
}

func NewUserNotificationService(db *bun.DB) *Service {
	return &Service{
		userNotificationRepository: storage.NewRepository[UserNotification](db),
	}
}

func (s *Service) ByID(ctx context.Context, id uuid.UUID) (*UserNotification, error) {
	return s.userNotificationRepository.ByID(ctx, id)
}

func (s *Service) One(ctx context.Context, opts ...storage.QueryOption) (*UserNotification, error) {
	return s.userNotificationRepository.One(ctx, opts...)
}

func (s *Service) List(ctx context.Context, opts ...storage.QueryOption) ([]*UserNotification, error) {
	return s.userNotificationRepository.List(ctx, opts...)
}

func (s *Service) Paginate(
	ctx context.Context,
	page int,
	perPage int,
	opts ...storage.QueryOption,
) (*storage.Paginator[*UserNotification], error) {
	return s.userNotificationRepository.Paginate(ctx, page, perPage, opts...)
}

func (s *Service) Update(ctx context.Context, user *UserNotification) (*UserNotification, error) {
	return s.userNotificationRepository.Update(ctx, user)
}

func (s *Service) Create(ctx context.Context, user *UserNotification) (*UserNotification, error) {
	return s.userNotificationRepository.Create(ctx, user)
}
