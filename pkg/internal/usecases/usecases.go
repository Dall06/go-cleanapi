package usecases

import (
	"dall06/go-cleanapi/pkg/internal"
	"fmt"

	"github.com/mitchellh/mapstructure"
)
type UseCases struct {
	repository Repository
}

func NewUseCases(r Repository) *UseCases {
	return &UseCases{
		repository: r,
	}
}

func (s *UseCases) RegisterUser(req interface{}) error {
    user := &internal.User{}
    
    if req == nil {
        return fmt.Errorf("empty request")
    }

    err := mapstructure.Decode(req, &user)
    if err != nil {
        return fmt.Errorf("failed to decode user details: %v", err)
    }

    err = s.repository.Create(user)
    if err != nil {
        return fmt.Errorf("failed to fetch user details: %v", err)
    }

	return nil
}

func (s *UseCases) IndexByID(req interface{}) (*internal.User, error) {
	user := &internal.User{}

    if req == nil {
        return nil, fmt.Errorf("empty request")
    }

    err := mapstructure.Decode(req, &user)
    if err != nil {
        return nil, fmt.Errorf("failed to decode user details: %v", err)
    }

    res, err := s.repository.Read(user)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch user details: %v", err)
    }

	return res, nil
}

func (s *UseCases) IndexAll() (internal.Users, error) {
    users, err := s.repository.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch user details: %v", err)
    }

	return users, nil
}

func (s *UseCases) ModifyUser(req interface{}) error {
    user := &internal.User{}

    if req == nil {
        return fmt.Errorf("empty request")
    }

    err := mapstructure.Decode(req, &user)
    if err != nil {
        return fmt.Errorf("failed to decode user details: %v", err)
    }

    err = s.repository.Update(user)
    if err != nil {
        return fmt.Errorf("failed to fetch user details: %v", err)
    }

	return nil
}

func (s *UseCases) DestroyUser(req interface{}) error {
    user := &internal.User{}

    if req == nil {
        return fmt.Errorf("empty request")
    }

    err := mapstructure.Decode(req, &user)
    if err != nil {
        return fmt.Errorf("failed to decode user details: %v", err)
    }

    err = s.repository.Delete(user)
    if err != nil {
        return fmt.Errorf("failed to fetch user details: %v", err)
    }

	return nil
}