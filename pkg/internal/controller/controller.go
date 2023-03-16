package controller

import (
	"dall06/go-cleanapi/pkg/internal/usecases"
	"dall06/go-cleanapi/utils"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	response utils.ResponseRepository
	usecases *usecases.UseCases
}

func NewController(uc *usecases.UseCases, r utils.ResponseRepository) *Controller {
	return &Controller{
		usecases: uc,
		response: r,
	}
}

func(c *Controller) GetByID(context *fiber.Ctx) error {
	// TODO: Implement
	return nil
}