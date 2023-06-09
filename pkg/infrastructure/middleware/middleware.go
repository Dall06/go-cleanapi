//go:build !coverage
// +build !coverage

// Package middleware Contains middleware implementation
// individual packages from middleware for fiber are already tested in their own github repository
package middleware

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/keyauth/v2"
)

// Middleware is an interface that extends middleware
type Middleware interface {
	CORS() fiber.Handler
	Helmet() fiber.Handler
	Compress() fiber.Handler
	EncryptCookie() fiber.Handler
	ETag() fiber.Handler
	Recover() fiber.Handler
	JwtWare() fiber.Handler
	KeyAuth() fiber.Handler
	CRSF() fiber.Handler
	Idempotency() fiber.Handler
}

var _ Middleware = (*middleware)(nil)

type middleware struct {
	jwt    utils.JWT
	config config.Vars
}

// NewMiddleware is a constructor for middleware
func NewMiddleware(vars config.Vars, jr utils.JWT) Middleware {
	return &middleware{
		jwt:    jr,
		config: vars,
	}
}

func (*middleware) CORS() fiber.Handler {
	cfg := &cors.Config{
		AllowOrigins:  "*",
		AllowHeaders:  "Origin,Content-Type,Accept,X-Session-Token,X-Application-Key",
		AllowMethods:  "GET,POST,PUT,DELETE",
		ExposeHeaders: "Content-Length,Authorization",
		MaxAge:        5600,
	}
	return cors.New(*cfg)
}

func (*middleware) Helmet() fiber.Handler {
	cfg := helmet.Config{
		CSPReportOnly: true,
	}

	return helmet.New(cfg)
}

func (*middleware) Compress() fiber.Handler {
	cfg := compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Level: compress.LevelBestSpeed, // 1
	}
	return compress.New(cfg)
}

func (m *middleware) EncryptCookie() fiber.Handler {
	cfg := encryptcookie.Config{
		Key: m.config.CookieSecret,
	}
	return encryptcookie.New(cfg)
}

func (*middleware) ETag() fiber.Handler {
	cfg := etag.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Weak: true,
	}
	return etag.New(cfg)
}

func (*middleware) Recover() fiber.Handler {
	return recover.New()
}

func (m *middleware) JwtWare() fiber.Handler {
	cfg := jwtware.Config{
		SigningKey:  m.config.JWTSecret,
		TokenLookup: "cookie:session_id",
		Filter: func(c *fiber.Ctx) bool {
			basePath := m.config.APIBasePath

			swaggerPath := fmt.Sprintf("%s/swagger/", basePath)

			usersPath := fmt.Sprintf("%s/users", basePath)
			authPath := fmt.Sprintf("%s/auth", usersPath)
			signupPath := fmt.Sprintf("%s/signup", usersPath)

			if c.Path() == swaggerPath {
				return true
			}
			if c.Path() == authPath {
				return true
			}
			if c.Path() == signupPath {
				return true
			}
			// Exclude all subroutes of /swagger
			if strings.HasPrefix(c.Path(), swaggerPath) {
				return true
			}

			return false
		},
	}

	return jwtware.New(cfg)
}

func (m *middleware) KeyAuth() fiber.Handler {
	cfg := keyauth.Config{
		KeyLookup: "header:x-access-token",
		Validator: func(c *fiber.Ctx, jwts string) (bool, error) {
			return m.jwt.CheckAPIJWT(jwts)
		},
		Filter: func(c *fiber.Ctx) bool {
			basePath := m.config.APIBasePath

			swaggerPath := fmt.Sprintf("%s/swagger/", basePath)

			usersPath := fmt.Sprintf("%s/users", basePath)
			authPath := fmt.Sprintf("%s/auth", usersPath)
			signupPath := fmt.Sprintf("%s/signup", usersPath)

			if c.Path() == swaggerPath {
				return true
			}
			if c.Path() == authPath {
				return true
			}
			if c.Path() == signupPath {
				return true
			}
			// Exclude all subroutes of /swagger
			if strings.HasPrefix(c.Path(), swaggerPath) {
				return true
			}

			return false
		},
	}
	return keyauth.New(cfg)
}

func (*middleware) CRSF() fiber.Handler {
	cfg := csrf.Config{
		Expiration: 15 * time.Minute,
	}
	// Or extend your config for customization
	return csrf.New(cfg)
}

func (*middleware) Idempotency() fiber.Handler {
	return idempotency.New()
}
