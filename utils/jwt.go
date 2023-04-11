package utils

import (
	"dall06/go-cleanapi/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtExpirationTime = 72 * time.Hour

type claims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

type JWTRepository interface {
	Create(uid string) (string, error)
	Check(requestToken string) (bool, error)
}

var _ JWTRepository = (*myJwt)(nil)

type myJwt struct {
	secret []byte
}

// NewJWT returns a pointer to a JwtUtil struct.
func NewJWT() JWTRepository {
	return &myJwt{
		secret: []byte(config.JwtSecret),
	}
}

func (ju *myJwt) Create(id string) (string, error) {
	if id == "" {
		return "", errors.New("id cannot be empty")
	}

	expiresAt := time.Now().Add(jwtExpirationTime)
	claims := claims{
		UID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// token -> string. Only server knows the secret.
	signedToken, err := token.SignedString(ju.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (ju *myJwt) Check(requestToken string) (bool, error) {
	if requestToken == "" {
		return false, errors.New("token cannot be empty")
	}

	userClaims := &claims{}
	token, err := jwt.ParseWithClaims(requestToken, userClaims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return ju.secret, nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	return true, nil
}

