package server

import (
	"dall06/go-cleanapi/pkg/infrastructure/database"
	"dall06/go-cleanapi/pkg/infrastructure/middleware"
	"dall06/go-cleanapi/pkg/internal/controller"
	"dall06/go-cleanapi/pkg/internal/repository"
	"dall06/go-cleanapi/pkg/internal/usecases"
	"dall06/go-cleanapi/pkg/server/routes"
	"dall06/go-cleanapi/utils"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	mycache *cache.Cache
	logger utils.LoggerRepository
	responses utils.ResponseRepository
}

func NewServer(c cache.Cache, l utils.LoggerRepository, rsp utils.ResponseRepository) *Server {
	return &Server{
		mycache: &c,
		logger: l,
		responses: rsp,
	}
}

func (s Server) Start() {

	cfg := fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "remit",
		AppName:       "remit-api-user-v1.0.0",
	}

	app := fiber.New(cfg)

	// init middleware
	mw := middleware.NewMiddleware()
	//app.Use(middleware.CORS)
	app.Use(mw.Compress())
	app.Use(mw.Helmet())
	app.Use(mw.EncryptCookie())
	app.Use(mw.ETag())
	app.Use(mw.Recover())
	app.Use(mw.JwtWare())
	app.Use(mw.KeyAuth())

	// init database
	dbConn := database.NewDBConn()
	conn, err := dbConn.Open()
	if err != nil {
		s.logger.Error("Failed to open database connection", err)
	}

	// generate internal controllers
	// user
	repo := repository.NewRepository(&conn)
	usecases := usecases.NewUseCases(repo)
	ctrl := controller.NewController(usecases, s.responses)

	// generate routing
	rts := routes.NewRoutesV1(*app, *s.mycache, *ctrl)
	rts.Set()
	
	// run gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
					fmt.Println("Gracefully shutting down...")
					_ = app.Shutdown()
	}()

	// ...

	if err := app.Listen(":8080"); err != nil {
		s.logger.Error("Failed to listen on port 8080", err)
		log.Panic(err)
	}

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
}