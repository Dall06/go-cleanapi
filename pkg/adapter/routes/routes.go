//go:build !coverage
// +build !coverage

// Package routes contains the routes per version
// this package should not be tested
// it is just an implementation and integration fo methos that are goning to be tested on integration
package routes

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/adapter/controller"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler

	// docs are generated by Swag CLI, you have to import them.
	_ "dall06/go-cleanapi/docs"
)

// Routes is an interface that extends routes
type Routes interface {
	Set()
}

var _ Routes = (*routes)(nil)

type routes struct {
	app        *fiber.App
	config     config.Vars
	controller controller.Controller
}

// NewRoutes is a constructor for routes generator
func NewRoutes(app *fiber.App, vars config.Vars, ctrl controller.Controller) Routes {
	return &routes{
		app:        app,
		config:     vars,
		controller: ctrl,
	}
}

func (routes *routes) Set() {
	basePath := routes.config.APIBasePath

	swaggerPath := fmt.Sprintf("%s/swagger/*", basePath)
	routes.app.Get(swaggerPath, swagger.HandlerDefault)

	usersPath := fmt.Sprintf("%s/users", basePath)
	usersGroup := routes.app.Group(usersPath)
	usersGroup.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("welcome to go-cleanapi user path ...")
	})
	usersGroup.Post("/auth", routes.controller.Auth)
	usersGroup.Post("/signup", routes.controller.Post)
	usersGroup.Get("/:id", routes.controller.Get)
	usersGroup.Get("/all", routes.controller.GetAll)
	usersGroup.Put("/modify/:id", routes.controller.Put)
	usersGroup.Delete("/delete/:id", routes.controller.Delete)
}
