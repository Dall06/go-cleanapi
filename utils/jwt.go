package utils

import (
	"dall06/go-cleanapi/config"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type claims struct {
	UID      string `json:"uid"`
	jwt.RegisteredClaims
}

type JWTUtilRespository interface {
	Create(uid string) (string, error)
	Check(requestToken string) (*jwt.Token, error)
}

var _ JWTUtilRespository = (*JwtUtil)(nil)

type JwtUtil struct{
	secret []byte
}

// NewLogger return a Logger.
func NewJWT() JWTUtilRespository {
	return JwtUtil{
		secret: []byte(config.JWTSecret),
	}
}

func(ju JwtUtil) Create(id string) (string, error) {
	if len(id) == 0 {
		return "", errors.New("empty id")
	}

	expiresAt := time.Now().Add(72 * time.Hour)
	claims := claims {
		UID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		// token -> string. Only server knows the secret.
	new, err := token.SignedString(ju.secret)
	if err != nil {
		return "", err
	}
	
	return new, nil
}

func (ju JwtUtil) Check(requestToken string) (*jwt.Token, error) {
	if len(requestToken) == 0 {
		return nil, errors.New("empty id")
	}

	userClaims := &claims{}
	token, err := jwt.ParseWithClaims(requestToken, userClaims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["tkn"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return ju.secret, nil })
	if err != nil {
		return nil, err
	}
	
	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", token.Header["alg"])
	}
	
	return token, nil
}