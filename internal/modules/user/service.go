package user

import "github.com/dbunt1tled/go-api/pkg/storage"

type UserService struct {
	userRepository *storage.Repository[User]
}
