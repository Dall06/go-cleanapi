// +build !coverage

package utils

import "github.com/google/uuid"

type mockUUIDRepository struct{}

func NewMockUUIDRepository() UUIDRepository {
	return &mockUUIDRepository{}
}

func (r *mockUUIDRepository) NewUUID() uuid.UUID {
	return uuid.New()
}

func (r *mockUUIDRepository) NewString() string {
	return uuid.NewString()
}