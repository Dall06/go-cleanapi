package middleware

import (
	"dall06/go-cleanapi/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
}

var _ MiddlewareRepository = (*Middleware)(nil)

type Middleware struct {}

func NewMiddleware() Middleware {
	return Middleware{}
}

func (Middleware) CORS() fiber.Handler {
    cfg := cors.Config{
        AllowOrigins:  "http://localhost:8080",
        AllowHeaders:  "Origin, Content-Type, Accept, X-Session-Token, X-Application-Key",
        AllowMethods:  "GET, POST, PUT, DELETE",
        ExposeHeaders: "Content-Length, Authorization",
        MaxAge:        5600,
    }
	return cors.New(cfg)
}

func (Middleware) Helmet() fiber.Handler {
	cfg := &helmet.Config{
		CSPReportOnly: true,
	}
	
	return helmet.New(*cfg)
}

func (Middleware) Compress() fiber.Handler {
	cfg := compress.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
		Level: compress.LevelBestSpeed, // 1
	}
	return compress.New(cfg)
}

func (Middleware) EncryptCookie() fiber.Handler {
	cfg := encryptcookie.Config{
		Key: config.EncryptCookie,
	}
	return encryptcookie.New(cfg)
}

func (Middleware) ETag() fiber.Handler {
	cfg := etag.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != "GET"
		},
		Weak: true,
	}
	return etag.New(cfg)
}

func (Middleware) Recover() fiber.Handler {
	return recover.New()
}

func (Middleware) JwtWare() fiber.Handler {
	cfg := jwtware.Config{
		SigningKey: []byte(config.JWTSecret),
		Filter: func(c *fiber.Ctx) bool {
			if c.Path() == config.BasePath {
				return true
			}
			
			if c.Path() == config.BasePath+"/user/hello" {
				return true
			}

			return false
		},
	}
	return jwtware.New(cfg)
}

func (Middleware) KeyAuth() fiber.Handler {
	cfg := keyauth.Config{
		KeyLookup: "header:x-access-token",
	}
	return keyauth.New(cfg)
}
