// +build !coverage

package utils

type jwtMockRepository struct {}

func NewMockJWTRepository() JWTRepository {
	return jwtMockRepository{}
}

func (j jwtMockRepository) Create(uid string) (string, error) {
	return "", nil
}

func (j jwtMockRepository) Check(requestToken string) (bool, error) {
	return true, nil
}