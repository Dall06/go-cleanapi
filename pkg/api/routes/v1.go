package routes

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/api/controller"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

type Routes struct {
	app        *fiber.App
	theCache   *cache.Cache
	controller controller.Controller
}

func NewRoutesV1(app fiber.App, c cache.Cache, ctrl controller.Controller) Routes {
	return Routes{
		app:      &app,
		theCache: &c,
		controller: ctrl,
	}
}

func (routes *Routes) Set() {
	var sb strings.Builder
	sb.WriteString(config.ProyectPath)
	sb.WriteString("/user")

	group := routes.app.Group(sb.String())

	group.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("hello to user route path, human!")
	})
	group.Group("/user/:id", routes.controller.Get)
}
