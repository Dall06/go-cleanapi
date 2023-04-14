//go:build !coverage
// +build !coverage

// Package utils is a package that provides general method for the api usage
package utils

type jwtMock struct{}

// NewJWTMock is a mock for jwt
func NewJWTMock() JWT {
	return &jwtMock{}
}

func (j *jwtMock) CreateUserJWT(uid string) (string, error) {
	return uid, nil
}

func (j *jwtMock) CheckUserJwt(_ string) (bool, error) {
	return true, nil
}

func (j *jwtMock) CreateAPIJWT() (string, error) {
	return "", nil
}

func (j *jwtMock) CheckAPIJwt(_ string) (bool, error) {
	return true, nil
}
