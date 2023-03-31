package controller

import (
	"dall06/go-cleanapi/pkg/internal/usecases"
	"dall06/go-cleanapi/utils"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
)

const (
	statusOK                  = fiber.StatusOK
	statusCreated             = fiber.StatusCreated
	statusBadRequest          = fiber.StatusBadRequest
	statusNotFound            = fiber.StatusNotFound
	statusInternalServerError = fiber.StatusInternalServerError

	requestError  = "request error"
	internalError = "internal error"
	notFound      = "not Found error"
	missingId     = "missing id parameter"
	userIsNil     = "user is null"
	usersAreNil   = "users are null"
	registered    = "account registered successfully"
	modified      = "account modified successfully"
	deleted       = "account deleted successfully"

	processed = "request processed"
)

type Controller interface {
	Post(context *fiber.Ctx) error
	Get(context *fiber.Ctx) error
	GetAll(context *fiber.Ctx) error
	Put(context *fiber.Ctx) error
	Delete(context *fiber.Ctx) error
}

type controller struct {
	usecases usecases.UseCases
	validate validator.Validate
	logger   utils.Logger
}

var _ Controller = (*controller)(nil)

func NewController(
	uc usecases.UseCases,
	v validator.Validate,
	l utils.Logger,
) Controller {
	return &controller{
		usecases: uc,
		validate: v,
		logger:   l,
	}
}

func (c *controller) Post(ctx *fiber.Ctx) error {
	req := &PostRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	if err := c.validate.Struct(req); err != nil {
		c.logger.Error("%s: %s", requestError, err)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, err))
	}

	userInput := &User{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}

	err := c.usecases.RegisterUser(userInput)
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	c.logger.Info("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": registered})
}

func (c *controller) Get(ctx *fiber.Ctx) error {
	// Get the id parameter from the request context
	id := ctx.Params("id")
	if id == "" {
		// Return an error response if the id parameter is missing
		c.logger.Error("%s: %s", requestError, missingId)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, missingId))
	}

	// Call the use case to retrieve the user by id
	userInput := &User{ID: id}
	userData, err := c.usecases.IndexByID(userInput)
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if userData == nil {
		c.logger.Error("%s: %s", statusNotFound, userIsNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", internalError, userIsNil))
	}

	// Convert the user data to the output format
	userOutput := &User{}
	err = mapstructure.Decode(userData, &userOutput)
	if err != nil {
		// Return an error response if the user data cannot be converted
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if err == sql.ErrNoRows {
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), notFound, userIsNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, err))
	}
	if userOutput == nil {
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), notFound, userIsNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, err))
	}

	// Return a success response with the user data
	c.logger.Info("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": userOutput})
}

func (c *controller) GetAll(ctx *fiber.Ctx) error {
	users, err := c.usecases.IndexAll()
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if users == nil {
		c.logger.Error("%s: %s", notFound, usersAreNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, usersAreNil))
	}

	// Convert the user data to the output format
	userOutput := &Users{}
	err = mapstructure.Decode(users, &userOutput)
	if err != nil {
		// Return an error response if the user data cannot be converted
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if userOutput == nil {
		c.logger.Error("%s: %s", notFound, usersAreNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, usersAreNil))
	}

	// Return a success response with the user data
	c.logger.Info("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": userOutput})
}

func (c *controller) Put(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		// Return an error response if the id parameter is missing
		c.logger.Error("%s: %s", statusBadRequest, missingId)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, missingId))
	}

	req := &PutRequest{}
	if err := ctx.BodyParser(req); err != nil {
		c.logger.Error("%s: %s", statusInternalServerError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	if err := c.validate.Struct(req); err != nil {
		c.logger.Error("%s: %s", statusBadRequest, err)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, err))
	}

	userInput := &User{
		ID:       id,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}
	err := c.usecases.ModifyUser(userInput)
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	c.logger.Info("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"msg": modified})
}

func (c *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		// Return an error response if the id parameter is missing
		c.logger.Error("%s: %s", requestError, missingId)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, missingId))
	}

	req := &DeleteRequest{}
	if err := ctx.BodyParser(req); err != nil {
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if err := c.validate.Struct(req); err != nil {
		c.logger.Error("%s: %s", requestError, missingId)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, err))
	}

	userInput := &User{
		ID:       id,
		Password: req.Password,
	}

	err := c.usecases.DestroyUser(userInput)
	if err != nil {
		c.logger.Error("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	c.logger.Info("%s -> [%s, %s] %s: (%s)", ctx.Method(), time.Now(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"msg": deleted})
}
