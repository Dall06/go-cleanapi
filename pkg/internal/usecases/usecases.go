// Package usecases contains bussines logig cases interface
package usecases

import (
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/pkg/internal/repository"
	"dall06/go-cleanapi/utils"
	"database/sql"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// UseCases is an interface that extend the cases
type UseCases interface {
	RegisterUser(req interface{}) error
	AuthUser(req interface{}) (*internal.User, error)
	IndexUserByID(req interface{}) (*internal.User, error)
	IndexUsers() (internal.Users, error)
	ModifyUser(req interface{}) error
	DestroyUser(req interface{}) error
}

var _ UseCases = (*cases)(nil)

type cases struct {
	repository repository.Repository
	uuid       utils.UUID
}

// NewUseCases is a construcotr for the cases
func NewUseCases(r repository.Repository, uid utils.UUID) UseCases {
	return &cases{
		repository: r,
		uuid:       uid,
	}
}

func (s *cases) AuthUser(req interface{}) (*internal.User, error) {
	user := &internal.User{}

	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	err := mapstructure.Decode(req, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user details: %v", err)
	}

	// add uuidGenerator to the user
	res, err := s.repository.Login(user)
	if err == sql.ErrNoRows {
		empty := &internal.User{}
		return empty, fmt.Errorf("failed to fetch user details: %v", "user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to auth user details: %v", err)
	}

	return res, nil
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

func (s *cases) IndexUserByID(req interface{}) (*internal.User, error) {
	user := &internal.User{}

	if req == nil {
		return nil, fmt.Errorf("empty request")
	}

	err := mapstructure.Decode(req, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user details: %v", err)
	}

	res, err := s.repository.Read(user)
	if err == sql.ErrNoRows {
		empty := &internal.User{}
		return empty, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user details: %v", err)
	}

	return res, nil
}

func (s *cases) IndexUsers() (internal.Users, error) {
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
		return fmt.Errorf("failed to update user: %v", err)
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
