package user

import (
	"context"

	"github.com/dbunt1tled/go-api/pkg/storage"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Service struct {
	userRepository *storage.Repository[*User]
}

func NewUserService(db *bun.DB) *Service {
	return &Service{
		userRepository: storage.NewRepository[*User](db),
	}
}

func (s *Service) ByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.userRepository.ByID(ctx, id)
}

func (s *Service) One(ctx context.Context, opts ...storage.QueryOption) (*User, error) {
	return s.userRepository.One(ctx, opts...)
}

func (s *Service) List(ctx context.Context, opts ...storage.QueryOption) ([]*User, error) {
	return s.userRepository.List(ctx, opts...)
}

func (s *Service) Paginate(
	ctx context.Context,
	page int,
	perPage int,
	opts ...storage.QueryOption,
) (*storage.Paginator[*User], error) {
	return s.userRepository.Paginate(ctx, page, perPage, opts...)
}

func (s *Service) Update(ctx context.Context, user *User) (*User, error) {
	return s.userRepository.Update(ctx, user)
}

func (s *Service) Create(ctx context.Context, user *User) (*User, error) {
	return s.userRepository.Create(ctx, user)
}
