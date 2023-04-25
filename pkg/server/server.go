// Package server runs the server configuration and initialization
//
//go:generate swag init -g server.go
package server

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/adapter/controller"
	"dall06/go-cleanapi/pkg/adapter/routes"
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
	config      config.Vars
	logger      utils.Logger
	jwt         utils.JWT
	uids        utils.UUID
	validations utils.Validations
	validation  validator.Validate
}

var _ Server = (*server)(nil)

// NewServer is a constructor for server
func NewServer(
	vars config.Vars,
	l utils.Logger,
	j utils.JWT,
	u utils.UUID,
	vs utils.Validations,
	v validator.Validate) Server {
	return server{
		config:      vars,
		logger:      l,
		jwt:         j,
		uids:        u,
		validations: vs,
		validation:  v,
	}
}

func (s server) Start() error {
	// init database
	dbConn := database.NewDBConn(s.logger, s.config)
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
	ctrl := controller.NewController(usecases, s.validation, s.logger, s.jwt, s.validations, *ctrlCache)

	// init server
	cfg := fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		ServerHeader:  "go-cleanapi",
		AppName:       s.config.AppName,
	}

	app := fiber.New(cfg)
	// init middleware
	mw := middleware.NewMiddleware(s.config, s.jwt)
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
	rts := routes.NewRoutes(app, s.config, ctrl)
	rts.Set()

	// run gracefully
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", s.config.APIPort)); err != nil {
			s.logger.Error("Failed to listen on port", err)
		}
	}()

	s.logger.Info("Running api server version %s in port %s, with base path %s",
		s.config.APIVersion, s.config.APIPort, s.config.APIBasePath)

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
