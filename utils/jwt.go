package utils

import (
	"crypto/sha512"
	"dall06/go-cleanapi/config"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwtExpirationTime = 72 * time.Hour

type userClaims struct {
	UID string `json:"uid"`
	jwt.RegisteredClaims
}

type apiClaims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

type JWTRepository interface {
	CreateUserJWT(uid string) (string, error)
	CheckUserJwt(requestToken string) (bool, error)
	CreateApiJWT() (string, error)
	CheckApiJwt(requestToken string) (bool, error)
}

var _ JWTRepository = (*myJwt)(nil)

type myJwt struct {
	secret []byte
	apiKey string
}

// NewJWT returns a pointer to a JwtUtil struct.
func NewJWT() JWTRepository {
	return &myJwt{
		secret: config.JwtSecret,
		apiKey: config.ApiKey,
	}
}

func (ju *myJwt) CreateUserJWT(id string) (string, error) {
	if id == "" {
		return "", errors.New("id cannot be empty")
	}

	expiresAt := time.Now().Add(jwtExpirationTime)
	userClaims := userClaims{
		UID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaims)

	// token -> string. Only server knows the secret.
	signedToken, err := token.SignedString(ju.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (ju *myJwt) CheckUserJwt(requestToken string) (bool, error) {
	if requestToken == "" {
		return false, errors.New("token cannot be empty")
	}

	claims := &userClaims{}
	token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (interface{}, error) {
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

func (ju *myJwt) CreateApiJWT() (string, error) {
	apiKey := config.ApiKey
	if apiKey == "" {
		return "", errors.New("id cannot be empty")
	}

	sha := sha512.Sum512_256([]byte(apiKey))
	hexString := hex.EncodeToString(sha[:])

	expiresAt := time.Now().Add(jwtExpirationTime)

	apiClaims := apiClaims{
		Hash: hexString,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, apiClaims)

	// token -> string. Only server knows the secret.
	signedToken, err := token.SignedString(ju.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (ju *myJwt) CheckApiJwt(requestToken string) (bool, error) {
	apiKeyHash := config.ApiKeyHash
	if apiKeyHash == "" {
		return false, errors.New("id cannot be empty")
	}

	if requestToken == "" {
		return false, errors.New("token cannot be empty")
	}

	claims := &apiClaims{}
	token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return ju.secret, nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return false, errors.New("invalid token")
	}

	fmt.Println("///////")
	fmt.Println(token.Claims)

	c, ok := token.Claims.(*apiClaims)
	if !ok {
		return false, fmt.Errorf("unexpected claims error: %v", token.Claims)
	}
	hashToCompare := c.Hash

	if apiKeyHash != hashToCompare {
		return false, fmt.Errorf("different hash error: %v", c)
	}

	return true, nil
}
