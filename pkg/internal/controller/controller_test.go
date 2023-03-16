package controller_test

import (
	"testing"
)

func TestIndex(test *testing.T) {
	/*db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rsp := utils.NewResponsesUtils()

	r := repository.NewRepository(*db)
	uc := usecases.NewUseCases(r)
	c := controller.NewController(uc, rsp)

	user := internal.User{
		ID:       "1",
		Email:    "email",
		Phone:    "phone",
		Password: "password",
	}

	userRows := sqlmock.NewRows([]string{"id", "email", "phone"}).
		AddRow(user.ID, user.Email, user.Phone)

	mock.ExpectQuery(regexp.QuoteMeta(spRead)).
		WithArgs(user.ID).
		WillReturnRows(userRows)

	app := fiber.New()

	test.Run("it should index an user model", func(t *testing.T) {

		app.Get("/test_index_controller/:id", c.Index)

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

		statusCode, ok := m["status"]
		if !ok {
			t.Errorf("map key not found: %s", err)
		}
		statusCodeStr := statusCode.(string)
		
		fmt.Println(m["content"])
		content, ok := m["content"]
		if !ok {
			t.Errorf("map key not found: %s", err)
		}
		fmt.Println(content)

		assert.Equalf(t, "ok", string(statusCodeStr), "status code match")
		assert.Equalf(t, 200, res.StatusCode, "status code match")
	})*/
}
