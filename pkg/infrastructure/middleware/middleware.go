//go:build !coverage
// +build !coverage

// individual packages from middleware for fiber are already tested in their own github repository

package middleware

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/keyauth/v2"
)

type MiddlewareRepository interface {
	CORS() fiber.Handler
	Helmet() fiber.Handler
	Compress() fiber.Handler
	EncryptCookie() fiber.Handler
	ETag() fiber.Handler
	Recover() fiber.Handler
	JwtWare() fiber.Handler
	KeyAuth() fiber.Handler
	CRSF() fiber.Handler
}

var _ MiddlewareRepository = (*middleware)(nil)

type middleware struct {
	jwt utils.JWTRepository
}

func NewMiddleware(jr utils.JWTRepository) MiddlewareRepository {
	return middleware{
		jwt: jr,
	}
}

func (middleware) CORS() fiber.Handler {
	cfg := &cors.Config{
		AllowOrigins:  "*",
		AllowHeaders:  "Origin,Content-Type,Accept,X-Session-Token,X-Application-Key",
		AllowMethods:  "GET,POST,PUT,DELETE",
		ExposeHeaders: "Content-Length,Authorization",
		MaxAge:        5600,
	}
	return cors.New(*cfg)
}

func (middleware) Helmet() fiber.Handler {
	cfg := helmet.Config{
		CSPReportOnly: true,
	}

	return helmet.New(cfg)
}

func (middleware) Compress() fiber.Handler {
	cfg := compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Level: compress.LevelBestSpeed, // 1
	}
	return compress.New(cfg)
}

func (middleware) EncryptCookie() fiber.Handler {
	cfg := encryptcookie.Config{
		Key: config.CookieSecret,
	}
	return encryptcookie.New(cfg)
}

func (middleware) ETag() fiber.Handler {
	cfg := etag.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Weak: true,
	}
	return etag.New(cfg)
}

func (middleware) Recover() fiber.Handler {
	return recover.New()
}

func (middleware) JwtWare() fiber.Handler {
	cfg := jwtware.Config{
		SigningKey:  config.JwtSecret,
		TokenLookup: "cookie:x-session-token",
		Filter: func(c *fiber.Ctx) bool {
			if c.Path() == config.ApiBasePath {
				return true
			}

			if c.Path() == config.ApiBasePath+"/v1/swagger/*" {
				return true
			}

			if c.Path() == config.ApiBasePath+"/v1/user/hello" {
				return true
			}

			if c.Path() == config.ApiBasePath+"/v1/welcome" {
				return true
			}

			return false
		},
	}

	return jwtware.New(cfg)
}

func (m middleware) KeyAuth() fiber.Handler {
	cfg := keyauth.Config{
		KeyLookup: "header:x-access-token",
		Validator: func(c *fiber.Ctx, jwts string) (bool, error) {
			return m.jwt.CheckApiJwt(jwts)
		},
		Filter: func(c *fiber.Ctx) bool {
			if c.Path() == "http://localhost:8080/swagger/doc.json" {
				return true
			}

			if c.Path() == config.ApiBasePath+"/v1/user/hello" {
				return true
			}

			if c.Path() == config.ApiBasePath+"/v1/welcome" {
				return true
			}

			if c.Path() == config.ApiBasePath+"/v1/swagger/*" {
				return true
			}

			return false
		},
	}
	return keyauth.New(cfg)
}

func (middleware) CRSF() fiber.Handler {
	cfg := csrf.Config{
		Expiration: 15 * time.Minute,
	}
	// Or extend your config for customization
	return csrf.New(cfg)
}
