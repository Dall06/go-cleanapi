// Package server runs the server configuration and initialization
//
//go:generate swag init -g server.go
package server

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/api/controller"
	"dall06/go-cleanapi/pkg/api/routes"
	"dall06/go-cleanapi/pkg/infrastructure/database"
	"dall06/go-cleanapi/pkg/infrastructure/middleware"
	"dall06/go-cleanapi/pkg/internal/repository"
	"dall06/go-cleanapi/pkg/internal/usecases"
	"dall06/go-cleanapi/utils"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

// Server is an interface for server
type Server interface {
	Start() error
}

type server struct {
	logger     utils.Logger
	jwt        utils.JWT
	uids       utils.UUID
	validation validator.Validate
}

var _ Server = (*server)(nil)

// NewServer is a constructor for server
func NewServer(
	l utils.Logger,
	j utils.JWT,
	u utils.UUID,
	v validator.Validate) Server {
	return server{
		logger:     l,
		jwt:        j,
		uids:       u,
		validation: v,
	}
}

func (s server) Start() error {
	// init database
	dbConn := database.NewDBConn(s.logger)
	conn, err := dbConn.Open()
	if err != nil {
		s.logger.Error("Failed to open database connection", err)
		return err
	}

	// generate caches, depending on the needs of each dependency
	ctrlCache := cache.New(5*time.Minute, 10*time.Minute)

	// generate internal controllers
	// user
	repo := repository.NewRepository(conn)
	usecases := usecases.NewUseCases(repo, s.uids)
	ctrl := controller.NewController(usecases, s.validation, s.logger, s.jwt, *ctrlCache)

	// init server
	cfg := fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		ServerHeader:  "go-cleanapi",
		AppName:       "go-cleanapi_v1.0.0",
	}

	app := fiber.New(cfg)
	// init middleware
	mw := middleware.NewMiddleware(s.jwt)
	app.Use(mw.CORS())
	app.Use(mw.Compress())
	app.Use(mw.Helmet())
	app.Use(mw.EncryptCookie())
	app.Use(mw.ETag())
	app.Use(mw.Recover())
	app.Use(mw.JwtWare())
	app.Use(mw.KeyAuth())
	app.Use(mw.CRSF())
	app.Use(mw.Idempotency())

	// generate routing
	rtsV1 := routes.NewRoutes(app, ctrl)
	rtsV1.SetV1()

	// run gracefully
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", config.APIPort)); err != nil {
			s.logger.Error("Failed to listen on port", err)
		}
		s.logger.Info("Running server...")
	}()

	// Gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Clean up tasks
	s.logger.Info("Shutting down server...")
	err = app.Shutdown()
	if err != nil {
		s.logger.Error("Failed to shutdown", err)
		return err
	}
	err = dbConn.Close(conn)
	if err != nil {
		s.logger.Error("Failed to close db connection")
		return err
	}

	// Before close
	s.logger.Info("Successfully shutdow nof the server")
	fmt.Println("shuted down...")
	return nil
}
