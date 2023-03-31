package utils

import "github.com/google/uuid"

type UUIDRepository interface {
	NewString() string
	NewUUID() uuid.UUID
}

var _ UUIDRepository = (*uuidGen)(nil)

type uuidGen struct {}

func NewUUIDGenerator() UUIDRepository {
	return &uuidGen{}
}

func (g *uuidGen) NewString() string {
	return uuid.NewString()
}

func (g *uuidGen) NewUUID() uuid.UUID {
	return uuid.New()
}