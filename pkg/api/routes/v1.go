// +build !coverage
// this package should not be tested
// it is just an implementation and integration fo methos that are goning to be tested on integration

package routes

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/api/controller"

	"strings"

	"github.com/gofiber/fiber/v2"
)

type RoutesRepository interface {
	Set()
}

var _ RoutesRepository = (*Routes)(nil)

type Routes struct {
	app        *fiber.App
	controller controller.Controller
}

func NewRoutesV1(app fiber.App, ctrl controller.Controller) RoutesRepository {
	return Routes{
		app:      &app,
		controller: ctrl,
	}
}

func (routes Routes) Set() {
	var sb strings.Builder
	sb.WriteString(config.ApiBasePath)
	sb.WriteString("/user")

	userPath := sb.String()
	userGroup := routes.app.Group(userPath)

	userGroup.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("hello to user route path, human!")
	})
	userGroup.Post("/", routes.controller.Post)
	userGroup.Get("/:id", routes.controller.Get)
	userGroup.Get("/all", routes.controller.GetAll)
	userGroup.Put("/:id", routes.controller.Put)
	userGroup.Delete("/:id", routes.controller.Delete)
}
