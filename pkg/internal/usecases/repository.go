package usecases

import "dall06/go-cleanapi/pkg/internal"

type Repository interface {
	Create(user *internal.User) error
	Read(user *internal.User) (*internal.User, error)
	ReadAll() (internal.Users, error)
	Update(user *internal.User) error
	Delete(user *internal.User) error
}