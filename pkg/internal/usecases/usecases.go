package usecases

import (
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/utils"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type UseCases interface {
	RegisterUser(req interface{}) error
	IndexByID(req interface{}) (*internal.User, error)
	IndexAll() (internal.Users, error)
	ModifyUser(req interface{}) error
	DestroyUser(req interface{}) error
}

var _ UseCases = (*cases)(nil)

type cases struct {
	repository Repository
	uuid       utils.UUIDRepository
}

func NewUseCases(r Repository, uid utils.UUIDRepository) UseCases {
	return &cases{
		repository: r,
		uuid:       uid,
	}
}

func (s *cases) RegisterUser(req interface{}) error {
	user := &internal.User{}

	if req == nil {
		return fmt.Errorf("empty request")
	}

	err := mapstructure.Decode(req, &user)
	if err != nil {
		return fmt.Errorf("failed to decode user details: %v", err)
	}

	// add uuidGenerator to the user
	user.ID = s.uuid.NewString()
	err = s.repository.Create(user)
	if err != nil {
		return fmt.Errorf("failed to fetch user details: %v", err)
	}

	return nil
}

func (s *cases) IndexByID(req interface{}) (*internal.User, error) {
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

func (s *cases) IndexAll() (internal.Users, error) {
	users, err := s.repository.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user details: %v", err)
	}

	return users, nil
}

func (s *cases) ModifyUser(req interface{}) error {
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

func (s *cases) DestroyUser(req interface{}) error {
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
