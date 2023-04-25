// Package controller_test is a test for controller
package controller_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/adapter/controller"
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/pkg/internal/repository"
	"dall06/go-cleanapi/pkg/internal/usecases"
	"dall06/go-cleanapi/utils"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

const (
	spCreate  = "CALL `go_cleanapi`.`sp_create_user`(?, ?, ?, ?);"
	spRead    = "CALL `go_cleanapi`.`sp_read_user`(?);"
	spReadAll = "CALL `go_cleanapi`.`sp_read_users`();"
	spUpdate  = "CALL `go_cleanapi`.`sp_update_user`(?, ?, ?, ?);"
	spDelete  = "CALL `go_cleanapi`.`sp_delete_user`(?, ?);"
	spLogin   = "CALL `go_cleanapi`.`sp_login_user`(?, ?, ?);"
)

func TestAuth(test *testing.T) {
	dbUserTwo := &internal.User{
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	rowsSetOne := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)

	formValuesEmail := url.Values{}
	formValuesEmail.Set("user", "test@test.com")
	formValuesEmail.Set("password", "12345pAsSWORd*")

	formValuesPhone := url.Values{}
	formValuesPhone.Set("user", "+991234567890")
	formValuesPhone.Set("password", "12345pAsSWORd*")

	successfulCases := []struct {
		name           string
		rows           *sqlmock.Rows
		dbUser         *internal.User
		reqForm        url.Values
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "it should login (mocked), whit email",
			rows:           rowsSetOne,
			dbUser:         dbUserTwo,
			reqForm:        formValuesEmail,
			expectedStatus: fiber.StatusCreated,
			expectedBody:   `{"msg":"account registered successfully"}`,
		},
	}

	app := fiber.New()

	cfg := config.NewConfig("8080", "1.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	v := validator.New()
	l := utils.NewLogger(*vars)
	uuid := utils.NewUUIDMock()
	jwt := utils.NewJWTMock()
	val := utils.NewValidations()

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()
			assert.Empty(test, err, "expected no error, but got:", err)

			m.ExpectQuery(regexp.QuoteMeta(spLogin)).WithArgs(
				&tc.dbUser.Email,
				&tc.dbUser.Phone,
				&tc.dbUser.Password,
			).WillReturnRows(tc.rows)
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Post("/auth", ctrl.Auth)

			req := httptest.NewRequest("POST", "/auth", strings.NewReader(tc.reqForm.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(req, 1)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
