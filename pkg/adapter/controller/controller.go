//go:build !exclude_Permision
// +build !exclude_Permision

// Package controller package is a package that provides handlers of an http to intercat with a data source
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
	"github.com/patrickmn/go-cache"
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
	missingID     = "missing id parameter"
	userIsNil     = "user is null"
	usersAreNil   = "users are null"
	registered    = "account registered successfully"
	modified      = "account modified successfully"
	deleted       = "account deleted successfully"

	processed = "request processed"
)

// Controller is an interface for controller
type Controller interface {
	Auth(context *fiber.Ctx) error
	Post(context *fiber.Ctx) error
	Get(context *fiber.Ctx) error
	GetAll(context *fiber.Ctx) error
	Put(context *fiber.Ctx) error
	Delete(context *fiber.Ctx) error
}

type controller struct {
	usecases    usecases.UseCases
	validate    validator.Validate
	logger      utils.Logger
	jwt         utils.JWT
	validations utils.Validations
	cache       *cache.Cache
}

var _ Controller = (*controller)(nil)

// NewController is a Constructor for controller
func NewController(
	uc usecases.UseCases,
	v validator.Validate,
	l utils.Logger,
	j utils.JWT,
	val utils.Validations,
	c cache.Cache,
) Controller {
	return &controller{
		usecases:    uc,
		validate:    v,
		logger:      l,
		jwt:         j,
		validations: val,
		cache:       &c,
	}
}

// @Summary Auth as user
// @Description auth a as user with phone or mail
// @Accept json
// @Produce json
// @Param user body PostRequest true "PostRequest object"
// @Success 200 {string} Accepted
// @Security ApiKeyAuth
// @Router /users/auth [post]
func (c *controller) Auth(ctx *fiber.Ctx) error {
	req := &AuthRequest{
		UserName: ctx.FormValue("user"),
		Password: ctx.FormValue("password"),
	}
	userInput := &User{}

	if err := c.validate.Struct(req); err != nil {
		c.logger.Error("%s: %s", requestError, err)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, err))
	}

	userName := req.UserName
	switch {
	case c.validations.IsEmail(userName):
		userInput.Email = userName
	case c.validations.IsPhone(userName):
		userInput.Phone = userName
	default:
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, "invalid user format"))
	}
	userInput.Password = req.Password

	res, err := c.usecases.AuthUser(userInput)
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if res.ID == "" {
		// Return an error response if the use case returns an error
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), requestError, missingID)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, missingID))
	}

	accessToken, err := c.jwt.CreateUserJWT(res.ID)
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, userIsNil)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, userIsNil))
	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = accessToken
	cookie.Expires = time.Now().Add(15 * time.Hour)

	c.logger.Info("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), processed, ctx.BaseURL())
	ctx.Cookie(cookie)
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{"msg": registered})
}

// @Summary Create a user
// @Description Create a new user
// @Accept json
// @Produce json
// @Param user body PostRequest true "PostRequest object"
// @Success 201 {string} Created
// @Security ApiKeyAuth
// @Router /users [post]
func (c *controller) Post(ctx *fiber.Ctx) error {
	req := &PostRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
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
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	c.logger.Info("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"msg": registered})
}

// @Summary Get a user by ID
// @Description Retrieve a single user by ID
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Security ApiKeyAuth
// @Security JwtTokenAuth
// @Router /users/{id} [get]
func (c *controller) Get(ctx *fiber.Ctx) error {
	// Get the id parameter from the request context
	id := ctx.Params("id")
	if id == "" {
		// Return an error response if the id parameter is missing
		c.logger.Error("%s: %s", requestError, missingID)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, missingID))
	}

	// Call the use case to retrieve the user by id
	userInput := &User{ID: id}
	empty := &User{}
	userData, err := c.usecases.IndexUserByID(userInput)
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
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
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if err == sql.ErrNoRows {
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), notFound, userIsNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, err))
	}
	if userOutput == nil {
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), notFound, userIsNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, err))
	}
	if userOutput == empty {
		c.logger.Info("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), notFound, userIsNil)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": empty, "msg": notFound})
	}

	// Return a success response with the user data
	c.logger.Info("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": userOutput})
}

// @Summary Get all users
// @Description Retrieve all users
// @Produce json
// @Success 200 {array} User
// @Security ApiKeyAuth
// @Security JwtTokenAuth
// @Router /users [get]
func (c *controller) GetAll(ctx *fiber.Ctx) error {
	// check if exists in cache, if yes returns value, if not, continues
	cachedUsers, found := c.cache.Get("users")
	if found {
		usersOutput := cachedUsers.(*Users)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": usersOutput})
	}

	users, err := c.usecases.IndexUsers()
	if err != nil {
		// Return an error response if the use case returns an error
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if users == nil {
		c.logger.Error("%s: %s", notFound, usersAreNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, usersAreNil))
	}

	// Convert the user data to the output format
	usersOutput := &Users{}
	err = mapstructure.Decode(users, &usersOutput)
	if err != nil {
		// Return an error response if the user data cannot be converted
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if usersOutput == nil {
		c.logger.Error("%s: %s", notFound, usersAreNil)
		return fiber.NewError(statusNotFound, fmt.Sprintf("%s: %s", notFound, usersAreNil))
	}
	if len(*usersOutput) == 0 {
		c.logger.Info("%s: %s", notFound, usersAreNil)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": usersOutput})
	}

	// Set new cache
	c.cache.Set("users", usersOutput, cache.DefaultExpiration)

	// Return a success response with the user data
	c.logger.Info("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": usersOutput})
}

// @Summary Update a user
// @Description Update a user with a given ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body PutRequest true "PutRequest object"
// @Success 200 {string} Updated
// @Security ApiKeyAuth
// @Security JwtTokenAuth
// @Router /users/{id} [put]
func (c *controller) Put(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		// Return an error response if the id parameter is missing
		c.logger.Error("%s: %s", statusBadRequest, missingID)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, missingID))
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
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	c.logger.Info("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"msg": modified})
}

// @Summary Delete a user
// @Description Delete a user with a given ID
// @Param id path int true "User ID"
// @Param user body DeleteRequest true "DeleteRequest object"
// @Success 204
// @Security ApiKeyAuth
// @Security JwtTokenAuth
// @Router /users/{id} [delete]
func (c *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		// Return an error response if the id parameter is missing
		c.logger.Error("%s: %s", requestError, missingID)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, missingID))
	}

	req := &DeleteRequest{}
	if err := ctx.BodyParser(req); err != nil {
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}
	if err := c.validate.Struct(req); err != nil {
		c.logger.Error("%s: %s", requestError, missingID)
		return fiber.NewError(statusBadRequest, fmt.Sprintf("%s: %s", requestError, err))
	}

	userInput := &User{
		ID:       id,
		Password: req.Password,
	}

	err := c.usecases.DestroyUser(userInput)
	if err != nil {
		c.logger.Error("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), internalError, err)
		return fiber.NewError(statusInternalServerError, fmt.Sprintf("%s: %s", internalError, err))
	}

	c.logger.Info("%s path[%s] -> %s: %s", ctx.Method(), ctx.Path(), processed, ctx.BaseURL())
	return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{"msg": deleted})
}
