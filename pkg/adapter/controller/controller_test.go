// Package controller_test is a test for controller
package controller_test

import (
	"bytes"
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/adapter/controller"
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/pkg/internal/repository"
	"dall06/go-cleanapi/pkg/internal/usecases"
	"dall06/go-cleanapi/utils"
	"database/sql"
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
	dbUserModel := &internal.User{
		Email:    "test@test.com",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}

	dbUserOne := &internal.User{
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	dbUserTwo := &internal.User{
		Email:    "",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}

	rowsSetOne := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetTwo := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetThree := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetFour := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetFive := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetSix := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetSeven := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetEight := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)
	rowsSetNine := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		"im an ID",
	)

	formValuesEmail := url.Values{}
	formValuesEmail.Set("user", "test@test.com")
	formValuesEmail.Set("password", "12345pAsSWORd*")
	formValuesEmail2 := url.Values{}
	formValuesEmail2.Set("user", "test2@test.com")
	formValuesEmail2.Set("password", "12345pAsSWORd*")
	formValuesBadEmail := url.Values{}
	formValuesBadEmail.Set("user", "testtest.com")
	formValuesBadEmail.Set("password", "12345pAsSWORd*")

	formValuesPhone := url.Values{}
	formValuesPhone.Set("user", "+991234567890")
	formValuesPhone.Set("password", "12345pAsSWORd*")
	formValuesPhone2 := url.Values{}
	formValuesPhone2.Set("user", "+521234567890")
	formValuesPhone2.Set("password", "12345pAsSWORd*")
	formValuesBadPhone := url.Values{}
	formValuesBadPhone.Set("user", "+991234ee7890")
	formValuesBadPhone.Set("password", "12345pAsSWORd*")

	formValuesEmpty := url.Values{}
	formValuesEmpty.Set("user", "")
	formValuesEmpty.Set("password", "")
	formValuesEmptyPass := url.Values{}
	formValuesEmptyPass.Set("user", "test@test.com")
	formValuesEmptyPass.Set("password", "")

	formValuesNil := url.Values{}

	successfulCases := []struct {
		testID         string
		name           string
		rows           *sqlmock.Rows
		dbUser         *internal.User
		reqForm        url.Values
		expectedStatus int
		expectedBody   *controller.User
	}{
		{
			testID:         "test1",
			name:           "it should login (mocked), whit email",
			rows:           rowsSetOne,
			dbUser:         dbUserOne,
			reqForm:        formValuesEmail,
			expectedStatus: fiber.StatusAccepted,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
		{
			testID:         "test2",
			name:           "it should login (mocked), whit phone",
			rows:           rowsSetTwo,
			dbUser:         dbUserTwo,
			reqForm:        formValuesPhone,
			expectedStatus: fiber.StatusAccepted,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
	}

	failedCases := []struct {
		testID         string
		name           string
		rows           *sqlmock.Rows
		dbUser         *internal.User
		reqForm        url.Values
		expectedStatus int
		expectedBody   *controller.User
	}{
		{
			testID:         "test3",
			name:           "it should not login (mocked), empty values",
			rows:           rowsSetThree,
			dbUser:         dbUserModel,
			reqForm:        formValuesEmpty,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
		{
			testID:         "test4",
			name:           "it should not login (mocked), nil values",
			rows:           rowsSetFour,
			dbUser:         dbUserModel,
			reqForm:        formValuesNil,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
		{
			testID:         "test5",
			name:           "it should not login (mocked), bad email",
			rows:           rowsSetFive,
			dbUser:         dbUserModel,
			reqForm:        formValuesBadEmail,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
		{
			testID:         "test6",
			name:           "it should not login (mocked), bad phone",
			rows:           rowsSetSix,
			dbUser:         dbUserModel,
			reqForm:        formValuesBadPhone,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
		{
			testID:         "test7",
			name:           "it should not login (mocked), empty pass",
			rows:           rowsSetSeven,
			dbUser:         dbUserModel,
			reqForm:        formValuesEmptyPass,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
		{
			testID:         "test8",
			name:           "it should not login (mocked), phone not found",
			rows:           rowsSetEight,
			dbUser:         dbUserModel,
			reqForm:        formValuesPhone2,
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
		{
			testID:         "test9",
			name:           "it should not login (mocked), email not found",
			rows:           rowsSetNine,
			dbUser:         dbUserModel,
			reqForm:        formValuesEmail2,
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: &controller.User{
				ID: "im an ID",
			},
		},
	}

	sCfg := fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	app := fiber.New(sCfg)

	cfg := config.NewConfig("8080", "1.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	v := validator.New()
	l := utils.NewLogger(*vars)
	err = l.Initialize()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	uuid := utils.NewUUIDMock()
	jwt := utils.NewJWTMock()
	val := utils.NewValidations()

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

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

			app.Post("/auth/"+tc.testID, ctrl.Auth)

			req := httptest.NewRequest("POST", "/auth/"+tc.testID, strings.NewReader(tc.reqForm.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(req, 1)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

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

			app.Post("/auth/"+tc.testID, ctrl.Auth)

			req := httptest.NewRequest("POST", "/auth/"+tc.testID, strings.NewReader(tc.reqForm.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := app.Test(req, 1)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestPost(test *testing.T) {
	dbUser1 := &internal.User{
		Email:    "test@test.com",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}
	dbUser2 := &internal.User{
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	successfulCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		body           string
		expectedStatus int
	}{
		{
			testID:         "test1",
			name:           "it should post user (mocked)",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusCreated,
			body:           `{"email":"test@test.com","phone":"+991234567890","password":"12345pAsSWORd*"}`,
		},
		{
			testID:         "test2",
			name:           "it should post user (mocked), whit empty phone",
			dbUser:         dbUser2,
			expectedStatus: fiber.StatusCreated,
			body:           `{"email":"test@test.com","phone":"","password":"12345pAsSWORd*"}`,
		},
	}

	failedCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		body           string
		expectedStatus int
	}{
		{
			testID:         "test3",
			name:           "it should post user (mocked), empty email",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusBadRequest,
			body:           `{"email":"","phone":"+991234567890","password":"12345pAsSWORd*"}`,
		},
		{
			testID:         "test4",
			name:           "it should post user (mocked), empty password",
			dbUser:         dbUser2,
			expectedStatus: fiber.StatusBadRequest,
			body:           `{"email":"test@test.com","phone":"","password":""}`,
		},
		{
			testID:         "test5",
			name:           "it should not post user (mocked), empty string",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusInternalServerError,
			body:           `{""}`,
		},
		{
			testID:         "test6",
			name:           "it should post user (mocked), internal error",
			dbUser:         dbUser2,
			expectedStatus: fiber.StatusInternalServerError,
			body:           `{"email":"test2@test.com","phone":"","password":"12345pAsSWORd*"}`,
		},
	}

	sCfg := fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	app := fiber.New(sCfg)

	cfg := config.NewConfig("8080", "1.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	v := validator.New()
	l := utils.NewLogger(*vars)
	err = l.Initialize()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	uuid := utils.NewUUIDMock()
	jwt := utils.NewJWTMock()
	val := utils.NewValidations()

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spCreate)).WithArgs(
				sqlmock.AnyArg(),
				&tc.dbUser.Email,
				&tc.dbUser.Phone,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 0))
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Post("/post/"+tc.testID, ctrl.Post)

			req := httptest.NewRequest("POST", "/post/"+tc.testID, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					test.Fatalf("error closing database connection: %v", err)
				}
			}()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spCreate)).WithArgs(
				sqlmock.AnyArg(),
				&tc.dbUser.Email,
				&tc.dbUser.Phone,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 0))
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Post("/post/"+tc.testID, ctrl.Post)

			req := httptest.NewRequest("POST", "/post/"+tc.testID, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					test.Fatalf("error closing database connection: %v", err)
				}
			}()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestGet(test *testing.T) {
	dbUser1 := &internal.User{
		ID:       "im_an_id",
		Email:    "test@test.com",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}

	rowsSet1 := sqlmock.NewRows([]string{
		"id_user",
		"user_email",
		"user_phone",
	}).AddRow(
		dbUser1.ID,
		dbUser1.Email,
		dbUser1.Phone,
	)

	rowsSet2 := sqlmock.NewRows([]string{
		"id_user",
		"user_email",
		"user_phone",
	})

	successfulCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		id             string
		expectedStatus int
		rows           *sqlmock.Rows
	}{
		{
			testID:         "test1",
			name:           "it should read user (mocked)",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusOK,
			rows:           rowsSet1,
			id:             "im_an_id",
		},
	}

	failedCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		id             string
		expectedStatus int
		rows           *sqlmock.Rows
	}{
		{
			testID:         "test2",
			name:           "it should not read user (mocked), empty id",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusOK,
			rows:           rowsSet1,
			id:             "im_an_id",
		},
		{
			testID:         "test3",
			name:           "it should not read user (mocked), empty db",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusNotFound,
			rows:           rowsSet2,
			id:             "",
		},
		{
			testID:         "test4",
			name:           "it should not read user (mocked), id not found",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusInternalServerError,
			rows:           rowsSet1,
			id:             "im_an_id_but_worng",
		},
	}

	sCfg := fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	app := fiber.New(sCfg)

	cfg := config.NewConfig("8080", "1.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	v := validator.New()
	l := utils.NewLogger(*vars)
	err = l.Initialize()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	uuid := utils.NewUUIDMock()
	jwt := utils.NewJWTMock()
	val := utils.NewValidations()

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spRead)).WithArgs(&dbUser1.ID).WillReturnRows(tc.rows)
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Get("/users/"+tc.testID+"/:id", ctrl.Get)

			// Make a request to the route with the test user ID
			req := httptest.NewRequest(fiber.MethodGet, "/users/"+tc.testID+"/"+tc.id, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spRead)).WithArgs(&dbUser1.ID).WillReturnRows(tc.rows)
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Get("/users/"+tc.testID+"/:id", ctrl.Get)

			// Make a request to the route with the test user ID
			req := httptest.NewRequest(fiber.MethodGet, "/users/"+tc.testID+"/"+tc.id, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone)
		})
	}
}

func TestGetAll(test *testing.T) {
	dbUser1 := &internal.User{
		ID:       "im_an_id",
		Email:    "test@test.com",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}
	dbUser2 := &internal.User{
		ID:    "im an id 2",
		Email: "test2@test.com",
		Phone: "+891234567891",
	}

	rowsSet1 := sqlmock.NewRows([]string{
		"id_user",
		"user_email",
		"user_phone",
	}).AddRow(
		dbUser1.ID,
		dbUser1.Email,
		dbUser1.Phone,
	).AddRow(
		dbUser2.ID,
		dbUser2.Email,
		dbUser2.Phone,
	)

	rowsSet2 := sqlmock.NewRows([]string{
		"id_user",
		"user_email",
		"user_phone",
	})

	successfulCases := []struct {
		testID         string
		name           string
		expectedStatus int
		rows           *sqlmock.Rows
	}{
		{
			testID:         "test1",
			name:           "it should read users (mocked)",
			expectedStatus: fiber.StatusOK,
			rows:           rowsSet1,
		},
		{
			testID:         "test2",
			name:           "it should read users (mocked), but empty db",
			expectedStatus: fiber.StatusOK,
			rows:           rowsSet2,
		},
	}

	sCfg := fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	app := fiber.New(sCfg)

	cfg := config.NewConfig("8080", "1.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	v := validator.New()
	l := utils.NewLogger(*vars)
	err = l.Initialize()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	uuid := utils.NewUUIDMock()
	jwt := utils.NewJWTMock()
	val := utils.NewValidations()

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spReadAll)).WillReturnRows(tc.rows)
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Get("/users/"+tc.testID, ctrl.GetAll)

			// Make a request to the route with the test user ID
			req := httptest.NewRequest(fiber.MethodGet, "/users/"+tc.testID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestPut(test *testing.T) {
	dbUser1 := &internal.User{
		ID:       "im_an_id",
		Email:    "test@test.com",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}
	dbUser2 := &internal.User{
		ID:       "im_an_id",
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}
	dbUser3 := &internal.User{
		ID:       "im_an_id",
		Email:    "",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}

	successfulCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		body           string
		id             string
		expectedStatus int
	}{
		{
			testID:         "test1",
			name:           "it should put user (mocked)",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusOK,
			id:             "im_an_id",
			body:           `{"email":"test@test.com","phone":"+991234567890","password":"12345pAsSWORd*"}`,
		},
		{
			testID:         "test2",
			name:           "it should put user (mocked), whit empty phone",
			dbUser:         dbUser2,
			expectedStatus: fiber.StatusOK,
			id:             "im_an_id",
			body:           `{"email":"test@test.com","phone":"","password":"12345pAsSWORd*"}`,
		},
		{
			testID:         "test3",
			name:           "it should put user (mocked), whit empty email",
			dbUser:         dbUser3,
			expectedStatus: fiber.StatusOK,
			id:             "im_an_id",
			body:           `{"email":"","phone":"+991234567890","password":"12345pAsSWORd*"}`,
		},
	}

	failedCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		body           string
		id             string
		expectedStatus int
	}{
		{
			testID:         "test4",
			name:           "it should post user (mocked), empty id",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusNotFound,
			id:             "",
			body:           `{"email":"","phone":"+991234567890","password":"12345pAsSWORd*"}`,
		},
		{
			testID:         "test5",
			name:           "it should post user (mocked), empty password",
			dbUser:         dbUser2,
			id:             "im_an_id",
			expectedStatus: fiber.StatusBadRequest,
			body:           `{"email":"test@test.com","phone":"","password":""}`,
		},
		{
			testID:         "test6",
			name:           "it should not post user (mocked), empty string",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusInternalServerError,
			id:             "im_an_id",
			body:           `{""}`,
		},
		{
			testID:         "test7",
			name:           "it should post user (mocked), internal error",
			dbUser:         dbUser2,
			id:             "im_an_id_2",
			expectedStatus: fiber.StatusInternalServerError,
			body:           `{"email":"test@test.com","phone":"","password":"12345pAsSWORd*"}`,
		},
	}

	sCfg := fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	app := fiber.New(sCfg)

	cfg := config.NewConfig("8080", "1.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	v := validator.New()
	l := utils.NewLogger(*vars)
	err = l.Initialize()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	uuid := utils.NewUUIDMock()
	jwt := utils.NewJWTMock()
	val := utils.NewValidations()

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spUpdate)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Email,
				&tc.dbUser.Phone,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 1))
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Put("/put/"+tc.testID+"/:id", ctrl.Put)

			req := httptest.NewRequest("PUT", "/put/"+tc.testID+"/"+tc.id, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					test.Fatalf("error closing database connection: %v", err)
				}
			}()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spUpdate)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Email,
				&tc.dbUser.Phone,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 1))
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Put("/put/"+tc.testID+"/:id", ctrl.Put)

			req := httptest.NewRequest("PUT", "/put/"+tc.testID+"/"+tc.id, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					test.Fatalf("error closing database connection: %v", err)
				}
			}()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestDelet(test *testing.T) {
	dbUser1 := &internal.User{
		ID:       "im_an_id",
		Email:    "test@test.com",
		Phone:    "+991234567890",
		Password: "12345pAsSWORd*",
	}
	successfulCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		body           string
		id             string
		expectedStatus int
	}{
		{
			testID:         "test1",
			name:           "it should delete user (mocked)",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusNoContent,
			id:             "im_an_id",
			body:           `{"password":"12345pAsSWORd*"}`,
		},
	}

	failedCases := []struct {
		testID         string
		name           string
		dbUser         *internal.User
		body           string
		id             string
		expectedStatus int
	}{
		{
			testID:         "test4",
			name:           "it should post user (mocked), empty id",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusNotFound,
			id:             "",
			body:           `{"password":"12345pAsSWORd*"}`,
		},
		{
			testID:         "test5",
			name:           "it should post user (mocked), empty password",
			dbUser:         dbUser1,
			id:             "im_an_id",
			expectedStatus: fiber.StatusBadRequest,
			body:           `{"password":""}`,
		},
		{
			testID:         "test6",
			name:           "it should not post user (mocked), empty string",
			dbUser:         dbUser1,
			expectedStatus: fiber.StatusInternalServerError,
			id:             "im_an_id",
			body:           `{""}`,
		},
		{
			testID:         "test7",
			name:           "it should post user (mocked), internal error",
			dbUser:         dbUser1,
			id:             "im_an_id_2",
			expectedStatus: fiber.StatusInternalServerError,
			body:           `{"password":"12345pAsSWORd*"}`,
		},
	}

	sCfg := fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	app := fiber.New(sCfg)

	cfg := config.NewConfig("8080", "1.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	v := validator.New()
	l := utils.NewLogger(*vars)
	err = l.Initialize()
	if err != nil {
		test.Fatalf("expected no error but got %v", err)
	}
	uuid := utils.NewUUIDMock()
	jwt := utils.NewJWTMock()
	val := utils.NewValidations()

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spDelete)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 1))
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Delete("/delete/"+tc.testID+"/:id", ctrl.Delete)

			req := httptest.NewRequest("DELETE", "/delete/"+tc.testID+"/"+tc.id, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					test.Fatalf("error closing database connection: %v", err)
				}
			}()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spDelete)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 1))
			assert.Empty(t, err, "expected no error, but got:", err)

			myCache := cache.New(5*time.Minute, 10*time.Minute)

			r := repository.NewRepository(db)
			uc := usecases.NewUseCases(r, uuid)
			ctrl := controller.NewController(uc, *v, l, jwt, val, *myCache)

			app.Delete("/delete/"+tc.testID+"/:id", ctrl.Delete)

			req := httptest.NewRequest("DELETE", "/delete/"+tc.testID+"/"+tc.id, bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			defer func() {
				if err := resp.Body.Close(); err != nil {
					test.Fatalf("error closing database connection: %v", err)
				}
			}()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}
