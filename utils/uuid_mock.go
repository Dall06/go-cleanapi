//go:build !coverage
// +build !coverage

// Package utils is a package that provides general method for the api usage
package utils

import "github.com/google/uuid"

type uuidMock struct{}

// NewUUIDMock is a contructor for a mock UUIDRepository
func NewUUIDMock() UUID {
	return &uuidMock{}
}

func (r *uuidMock) NewUUID() uuid.UUID {
	return uuid.New()
}

func (r *uuidMock) NewString() string {
	return uuid.NewString()
}
