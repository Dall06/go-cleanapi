// Package utils is a package that provides general method for the api usage
package utils

import "github.com/google/uuid"

// UUID is an interface for uuidGen
type UUID interface {
	NewString() string
	NewUUID() uuid.UUID
}

var _ UUID = (*uuidGen)(nil)

type uuidGen struct{}

// NewUUIDGenerator is a constructir for uidGen
func NewUUIDGenerator() UUID {
	return &uuidGen{}
}

func (g *uuidGen) NewString() string {
	return uuid.NewString()
}

func (g *uuidGen) NewUUID() uuid.UUID {
	return uuid.New()
}
