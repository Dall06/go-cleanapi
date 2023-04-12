// +build !coverage

package utils

type jwtMockRepository struct {}

func NewMockJWTRepository() JWTRepository {
	return jwtMockRepository{}
}

func (j jwtMockRepository) CreateUserJWT(uid string) (string, error) {
	return "", nil
}

func (j jwtMockRepository) CheckUserJwt(requestToken string) (bool, error) {
	return true, nil
}

func (j jwtMockRepository) CreateApiJWT() (string, error) {
	return "", nil
}

func (j jwtMockRepository) CheckApiJwt(requestToken string) (bool, error) {
	return true, nil
}