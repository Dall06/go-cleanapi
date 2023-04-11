package controller_test

import (
	"bufio"
	"bytes"
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/api/controller"
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/pkg/internal/repository"
	"dall06/go-cleanapi/pkg/internal/usecases"
	"dall06/go-cleanapi/utils"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

const (
	spCreate  = "CALL `go_cleanapi`.`sp_create_user`(?, ?, ?, ?);"
	spRead    = "CALL `go_cleanapi`.`sp_read_user`(?);"
	spReadAll = "CALL `go_cleanapi`.`sp_read_users`();"
	spUpdate  = "CALL `go_cleanapi`.`sp_update_user`(?, ?, ?, ?);"
	spDelete  = "CALL `go_cleanapi`.`sp_delete_user`(?, ?);"
)

var user = &internal.User{
	ID:       "1",
	Email:    "johndoe@example.com",
	Phone:    "+1234567890",
	Password: "password123",
}

var expectedModel = &controller.User{
	ID:       "1",
	Email:    "johndoe@example.com",
	Phone:    "+1234567890",
	Password: "",
}

var postReq = controller.PostRequest{
	Email:    "johndoe@example.com",
	Phone:    "+1234567890",
	Password: "password123",
}

var putReq = controller.PutRequest{
	Email:    "johndoe@example.com",
	Phone:    "+1234567890",
	Password: "password123",
}

var delReq = controller.DeleteRequest{
	Password: "password123",
}

func TestPost(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewRepository(db)
	if r == nil {
		test.Fatalf("an error was not expected when creating repository")
	}

	user := &internal.User{
		ID:       "1",
		Email:    "johndoe@example.com",
		Phone:    "+1234567890",
		Password: "password123",
	}

	conf := config.NewConfig("8080")
	err = conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to create config: %v", err)
	}

	l := utils.NewLogger()
	if l == nil {
		test.Fatalf("failed to load logger")
	}

	err = l.Initialize()
	if err != nil {
		test.Fatalf("failed to initialize logger: %v", err)
	}

	v := validator.New()

	uidRepo := utils.NewMockUUIDRepository()
	uc := usecases.NewUseCases(r, uidRepo)
	c := cache.New(5*time.Minute, 10*time.Minute)
	ctrl := controller.NewController(uc, *v, l, *c)

	app := fiber.New(fiber.Config{
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
	})

	// mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(spCreate)).
		WithArgs(sqlmock.AnyArg(), user.Email, sqlmock.AnyArg(), user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	test.Run("it should create", func(t *testing.T) {
		// Encode the person struct as a JSON string
		body, err := json.Marshal(postReq)
		if err != nil {
			panic(err)
		}

		app.Post("/test_post_controller/", ctrl.Post)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_post_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 201 {
			t.Errorf("Expected status code %d, but got %d", 201, res.StatusCode)
		}
	})

	test.Run("it should not create, empty req", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PostRequest{}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Post("/test_post_controller/", ctrl.Post)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_post_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not create, nil req", func(t *testing.T) {
		// Encode the person struct as a JSON string
		body, err := json.Marshal(nil)
		if err != nil {
			panic(err)
		}

		app.Post("/test_post_controller/", ctrl.Post)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_post_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not create, empty email", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PostRequest{
			Password: "password123",
			Phone:    "+1234567890",
		}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Post("/test_post_controller/", ctrl.Post)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_post_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not create, empty password", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PostRequest{
			Email:    "johndoe@example.com",
			Password: "",
			Phone:    "+1234567890",
		}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Post("/test_post_controller/", ctrl.Post)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_post_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not create, bad email", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PostRequest{
			Email:    "johndoeexample.com",
			Password: "password123",
			Phone:    "+1234567890",
		}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Post("/test_post_controller/", ctrl.Post)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_post_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not create, bad pass", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PostRequest{
			Email:    "johndoe@example.com",
			Password: "",
			Phone:    "",
		}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Post("/test_post_controller/", ctrl.Post)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_post_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})
}

func TestGet(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewRepository(db)
	if r == nil {
		test.Fatalf("an error was not expected when creating repository")
	}

	columns := []string{"id", "email", "phone"}
	mock.ExpectQuery(regexp.QuoteMeta(spRead)).
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(user.ID, user.Email, user.Phone))

	conf := config.NewConfig("8080")
	err = conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to create config: %v", err)
	}

	l := utils.NewLogger()
	if l == nil {
		test.Fatalf("failed to load logger")
	}

	err = l.Initialize()
	if err != nil {
		test.Fatalf("failed to initialize logger: %v", err)
	}

	v := validator.New()

	uidRepo := utils.NewMockUUIDRepository()
	uc := usecases.NewUseCases(r, uidRepo)
	c := cache.New(5*time.Minute, 10*time.Minute)
	ctrl := controller.NewController(uc, *v, l, *c)

	app := fiber.New(fiber.Config{
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
	})

	test.Run("it should index an user model", func(t *testing.T) {

		app.Get("/test_index_controller/:id", ctrl.Get)

		req := httptest.NewRequest("GET", "/test_index_controller/"+user.ID, nil)
		res, err := app.Test(req)
		if err != nil {
			t.Errorf("error was not expected when indexing data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 200 {
			t.Errorf("Expected status code %d, but got %d", 200, res.StatusCode)
		}

		var buf bytes.Buffer
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			buf.Write(scanner.Bytes())
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("Failed to read response body: %s", err)
		}

		// Convert the response body to a map
		var m map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Errorf("Failed to deserialize response body to map: %s", err)
		}

		fmt.Println(m)
		content, ok := m["data"]
		if !ok {
			t.Errorf("map key not found: %s", err)
		}
		fmt.Println(content)

		var result map[string]interface{}
		err = mapstructure.Decode(content, &result)
		if err != nil {
			t.Errorf("caonnot convert interface to map: %s", err)
		}

		model := &controller.User{
			ID:    result["uid"].(string),
			Email: result["mail"].(string),
			Phone: result["phone"].(string),
		}

		assert.Equal(t, expectedModel, model)
	})

	test.Run("it should not index an user model, empty ID", func(t *testing.T) {

		app.Get("/test_index_controller/:id", ctrl.Get)
		id := ""

		req := httptest.NewRequest("GET", "/test_index_controller/"+id, nil)
		res, err := app.Test(req)
		if err != nil {
			t.Errorf("error was not expected when indexing data: %s", err)
		}

		fmt.Println(res)
		fmt.Println(res.StatusCode)

		// Check that the response status code is 200 OK
		if res.StatusCode != 404 {
			t.Errorf("status code not expected: %d", res.StatusCode)
		}
		t.Log("failed successfully")
	})

	test.Run("it should not index an user model, empty ID", func(t *testing.T) {

		app.Get("/test_index_controller/:id", ctrl.Get)
		id := ""

		req := httptest.NewRequest("GET", "/test_index_controller/"+id, nil)
		res, err := app.Test(req)
		if err != nil {
			t.Errorf("error was not expected when indexing data: %s", err)
		}

		fmt.Println(res)
		fmt.Println(res.StatusCode)

		// Check that the response status code is 200 OK
		if res.StatusCode != 404 {
			t.Errorf("status code not expected: %d", res.StatusCode)
		}
		t.Log("failed successfully")
	})

	test.Run("it should not index an user model, empty ID", func(t *testing.T) {

		app.Get("/test_index_controller/:id", ctrl.Get)
		id := ""

		req := httptest.NewRequest("GET", "/test_index_controller/"+id, nil)
		res, err := app.Test(req)
		if err != nil {
			t.Errorf("error was not expected when indexing data: %s", err)
		}

		fmt.Println(res)
		fmt.Println(res.StatusCode)

		// Check that the response status code is 200 OK
		if res.StatusCode != 404 {
			t.Errorf("status code not expected: %d", res.StatusCode)
		}
		t.Log("failed successfully")
	})
}

func TestGetAll(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewRepository(db)
	if r == nil {
		test.Fatalf("an error was not expected when creating repository")
	}

	columns := []string{"id", "email", "phone"}

	conf := config.NewConfig("8080")
	err = conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to create config: %v", err)
	}

	l := utils.NewLogger()
	if l == nil {
		test.Fatalf("failed to load logger")
	}

	err = l.Initialize()
	if err != nil {
		test.Fatalf("failed to initialize logger: %v", err)
	}

	v := validator.New()

	uidRepo := utils.NewMockUUIDRepository()
	uc := usecases.NewUseCases(r, uidRepo)
	c := cache.New(5*time.Microsecond, 10*time.Microsecond)
	ctrl := controller.NewController(uc, *v, l, *c)
	if uc == nil {
		test.Fatalf("an error was not expected when creating usecases")
	}

	expected_one := controller.User{
		ID:    "1",
		Email: "test1@test.com",
		Phone: "1234567890",
	}

	expected_two := controller.User{
		ID:    "2",
		Email: "test2@test.com",
		Phone: "0987654321",
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
	})

	test.Run("it should index an user slice", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(spReadAll)).
			WillReturnRows(sqlmock.NewRows(columns).
				AddRow("1", "test1@test.com", "1234567890").
				AddRow("2", "test2@test.com", "0987654321"))

		app.Get("/test_index_all_controller/", ctrl.GetAll)

		req := httptest.NewRequest("GET", "/test_index_all_controller/", nil)
		res, err := app.Test(req)
		if err != nil {
			t.Errorf("error was not expected when indexing data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 200 {
			t.Errorf("Expected status code %d, but got %d", 200, res.StatusCode)
		}

		var buf bytes.Buffer
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			buf.Write(scanner.Bytes())
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("Failed to read response body: %s", err)
		}

		var m map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Errorf("Failed to deserialize response body to map: %s", err)
		}

		fmt.Println(m)
		content, ok := m["data"]
		if !ok {
			t.Errorf("map key not found: %s", err)
		}
		fmt.Println(content)

		var users controller.Users

		for _, d := range content.([]interface{}) {
			u := controller.User{}
			m := d.(map[string]interface{})
			u.ID = m["uid"].(string)
			u.Email = m["mail"].(string)
			u.Phone = m["phone"].(string)
			users = append(users, u)
		}

		fmt.Println(users)
		if len(users) == 0 {
			t.Errorf("no users inside")
		}

		assert.Equal(t, expected_one, users[0])
		assert.Equal(t, expected_two, users[1])

	})

	test.Run("it should not read, db is empty", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(spReadAll)).
			WillReturnRows(
				sqlmock.NewRows(columns))

		app.Get("/test_index_all_controller/", ctrl.GetAll)

		req := httptest.NewRequest("GET", "/test_index_all_controller/", nil)
		res, err := app.Test(req)
		if err != nil {
			t.Errorf("error was not expected when indexing data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 200 {
			t.Errorf("Expected status code %d, but got %d", 200, res.StatusCode)
		}

		var buf bytes.Buffer
		scanner := bufio.NewScanner(res.Body)
		for scanner.Scan() {
			buf.Write(scanner.Bytes())
		}
		if err := scanner.Err(); err != nil {
			t.Errorf("Failed to read response body: %s", err)
		}

		var m map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &m)
		if err != nil {
			t.Errorf("Failed to deserialize response body to map: %s", err)
		}

		fmt.Println(m)
		content, ok := m["data"]
		if !ok {
			t.Errorf("map key not found: %s", err)
		}
		fmt.Println(content)

		var users controller.Users

		for _, d := range content.([]interface{}) {
			u := controller.User{}
			m := d.(map[string]interface{})
			u.ID = m["uid"].(string)
			u.Email = m["mail"].(string)
			u.Phone = m["phone"].(string)
			users = append(users, u)
		}
		if len(users) != 0 {
			t.Errorf("expected to be empty")
		}

		for _, u := range users {
			if u.ID == "1" {
				assert.NotEqual(t, expected_one, u)
			}

			if u.ID == "2" {
				assert.NotEqual(t, expected_two, u)
			}
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled expectations: %s", err)
		}
	})
}

func TestPut(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewRepository(db)
	if r == nil {
		test.Fatalf("an error was not expected when creating repository")
	}

	user := &internal.User{
		ID:       "1",
		Email:    "johndoe@example.com",
		Phone:    "+1234567890",
		Password: "password123",
	}

	conf := config.NewConfig("8080")
	err = conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to create config: %v", err)
	}

	l := utils.NewLogger()
	if l == nil {
		test.Fatalf("failed to load logger")
	}

	err = l.Initialize()
	if err != nil {
		test.Fatalf("failed to initialize logger: %v", err)
	}

	v := validator.New()

	uidRepo := utils.NewMockUUIDRepository()
	uc := usecases.NewUseCases(r, uidRepo)
	c := cache.New(5*time.Minute, 10*time.Minute)
	ctrl := controller.NewController(uc, *v, l, *c)

	app := fiber.New(fiber.Config{
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
	})

	// mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(spUpdate)).
		WithArgs(sqlmock.AnyArg(), user.Email, sqlmock.AnyArg(), user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	test.Run("it should put", func(t *testing.T) {
		// Encode the person struct as a JSON string
		body, err := json.Marshal(putReq)
		if err != nil {
			panic(err)
		}

		app.Put("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("PUT", "/test_put_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 200 {
			t.Errorf("Expected status code %d, but got %d", 200, res.StatusCode)
		}
	})

	test.Run("it should not put, empty req", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PutRequest{}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Put("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("PUT", "/test_put_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not put, nil req", func(t *testing.T) {
		// Encode the person struct as a JSON string
		body, err := json.Marshal(nil)
		if err != nil {
			panic(err)
		}

		app.Put("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("PUT", "/test_put_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not put, empty email", func(t *testing.T) {
		pr := &controller.PutRequest{
			Phone:    "+1234567890",
			Password: "password123",
		}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Put("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("PUT", "/test_put_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not put, empty password", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PostRequest{
			Email:    "johndoeexample.com",
			Password: "",
			Phone:    "+1234567890",
		}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Put("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("PUT", "/test_put_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not put, bad email", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PutRequest{
			Email:    "johndoeexample.com",
			Password: "password123",
			Phone:    "+1234567890",
		}

		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Put("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("PUT", "/test_put_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not put, no id param", func(t *testing.T) {
		// Encode the person struct as a JSON string
		body, err := json.Marshal(putReq)
		if err != nil {
			panic(err)
		}

		app.Post("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_put_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 404 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})
}

func TestDelete(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewRepository(db)
	if r == nil {
		test.Fatalf("an error was not expected when creating repository")
	}

	user := &internal.User{
		ID:       "1",
		Email:    "johndoe@example.com",
		Phone:    "+1234567890",
		Password: "password123",
	}

	conf := config.NewConfig("8080")
	err = conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to create config: %v", err)
	}

	l := utils.NewLogger()
	if l == nil {
		test.Fatalf("failed to load logger")
	}

	err = l.Initialize()
	if err != nil {
		test.Fatalf("failed to initialize logger: %v", err)
	}

	v := validator.New()

	uidRepo := utils.NewMockUUIDRepository()
	uc := usecases.NewUseCases(r, uidRepo)
	c := cache.New(5*time.Minute, 10*time.Minute)
	ctrl := controller.NewController(uc, *v, l, *c)

	app := fiber.New(fiber.Config{
		ReadTimeout:  50 * time.Second,
		WriteTimeout: 50 * time.Second,
	})

	// mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(spDelete)).
		WithArgs(sqlmock.AnyArg(), user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	test.Run("it should delete", func(t *testing.T) {
		// Encode the person struct as a JSON string
		body, err := json.Marshal(delReq)
		if err != nil {
			panic(err)
		}

		app.Delete("/test_delete_controller/:id", ctrl.Delete)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("DELETE", "/test_delete_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 200 {
			t.Errorf("Expected status code %d, but got %d", 200, res.StatusCode)
		}
	})

	test.Run("it should not delete, empty req", func(t *testing.T) {
		// Encode the person struct as a JSON string
		dr := &controller.DeleteRequest{}
		body, err := json.Marshal(dr)
		if err != nil {
			panic(err)
		}

		app.Delete("/test_delete_controller/:id", ctrl.Delete)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("DELETE", "/test_delete_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not delete, nil req", func(t *testing.T) {
		// Encode the person struct as a JSON string
		body, err := json.Marshal(nil)
		if err != nil {
			panic(err)
		}

		app.Post("/test_delete_controller/:id", ctrl.Delete)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_delete_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not delete, empty password", func(t *testing.T) {
		// Encode the person struct as a JSON string
		dr := &controller.DeleteRequest{
			Password: "",
		}
		body, err := json.Marshal(dr)
		if err != nil {
			panic(err)
		}

		app.Delete("/test_delete_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("DELETE", "/test_delete_controller/1", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 400 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})

	test.Run("it should not create, no id param", func(t *testing.T) {
		// Encode the person struct as a JSON string
		pr := &controller.PostRequest{
			Email:    "johndoeexample.com",
			Password: "password123",
			Phone:    "+1234567890",
		}
		body, err := json.Marshal(pr)
		if err != nil {
			panic(err)
		}

		app.Post("/test_put_controller/:id", ctrl.Put)
		reqBody := bytes.NewBuffer(body)
		req := httptest.NewRequest("POST", "/test_put_controller/", reqBody)
		req.Header.Set("Content-Type", "application/json") // add content type header
		res, err := app.Test(req)

		if err != nil {
			t.Errorf("error was not expected when post data: %s", err)
		}

		// Check that the response status code is 200 OK
		if res.StatusCode != 404 {
			t.Errorf("Expected status code %d, but got %d", 400, res.StatusCode)
		}
	})
}