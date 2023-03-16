package usecases_test

import (
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/pkg/internal/controller"
	"dall06/go-cleanapi/pkg/internal/repository"
	"dall06/go-cleanapi/pkg/internal/usecases"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

const (
	spCreate  = "CALL `go_cleanapi`.`sp_create_user`(?, ?, ?, ?);"
	spRead    = "CALL `go_cleanapi`.`sp_read_user`(?);"
	spReadAll = "CALL `go_cleanapi`.`sp_read_users`();"
	spUpdate  = "CALL `go_cleanapi`.`sp_update_user`(?, ?, ?, ?);"
	spDelete  = "CALL `go_cleanapi`.`sp_delete_user`(?, ?);"
)

func TestRegister(test *testing.T) {
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
		Password: "password",
	}

	// mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(spCreate)).
		WithArgs(user.ID, user.Email, user.Phone, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	uc := usecases.NewUseCases(r)
	if uc == nil {
		test.Fatalf("an error was not expected when creating usecases")
	}

	test.Run("it should create", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err != nil {
			t.Fatalf("error was not expected while registering an user: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, empty ID", func(t *testing.T) {
		u := &controller.User{
			ID:       "",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, empty User", func(t *testing.T) {
		u := &controller.User{}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, User is nil", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = uc.RegisterUser(nil)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, empty Email", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, empty password", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "",
		}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should create, even if it has empty phone data", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}

func TestIndexByID(test *testing.T) {
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
		Email:    "email",
		Phone:    "phone",
		Password: "password",
	}

	model := &controller.User{
		ID:       "1",
		Email:    "",
		Phone:    "",
		Password: "password",
	}

	columns := []string{"id", "email", "phone"}
	mock.ExpectQuery(regexp.QuoteMeta(spRead)).
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(user.ID, user.Email, user.Phone))

	uc := usecases.NewUseCases(r)

	test.Run("it should index an user", func(t *testing.T) {
		m := &controller.User{}

		i, err := uc.IndexByID(model)
		if err != nil {
			t.Errorf("error was not expected while indexing data: %s", err)
		}

		err = mapstructure.Decode(i, &m)
		if err != nil {
			t.Errorf("error was not expected while decoding indexed data: %s", err)
		}

		fmt.Println(m)
		t.Log("successfully indexed")
	})

	test.Run("it should not index, model is nil", func(t *testing.T) {
		m := &controller.User{}

		i, err := uc.IndexByID(nil)
		if err == nil {
			t.Errorf("error was expected")
		}

		err = mapstructure.Decode(i, &m)
		if err != nil {
			t.Errorf("error was not expected while decoding indexed data: %s", err)
		}

		t.Log("successfully failed to index")
	})

	test.Run("it should not index, model is empty", func(t *testing.T) {
		model := &controller.User{}
		m := &controller.User{}

		i, err := uc.IndexByID(model)
		if err == nil {
			t.Errorf("error was expected")
		}

		err = mapstructure.Decode(i, &m)
		if err != nil {
			t.Errorf("error was not expected while decoding indexed data: %s", err)
		}
		
		t.Log("successfully failed to index")
	})

	test.Run("it should not index, model ID is empty string", func(t *testing.T) {
		model := &controller.User{
			ID: "",
		}
		m := &controller.User{}

		i, err := uc.IndexByID(model)
		if err == nil {
			t.Errorf("error was expected")
		}

		err = mapstructure.Decode(i, &m)
		if err != nil {
			t.Errorf("error was not expected while decoding indexed data: %s", err)
		}
		
		t.Log("successfully failed to index")
	})

	test.Run("it should not index, model ID is not in db", func(t *testing.T) {
		model := &controller.User{
			ID: "3",
		}
		m := &controller.User{}

		i, err := uc.IndexByID(model)
		if err == nil {
			t.Errorf("error was expected")
		}

		err = mapstructure.Decode(i, &m)
		if err != nil {
			t.Errorf("error was not expected while decoding indexed data: %s", err)
		}
		
		t.Log("successfully failed to index")
	})
}

func TestIndexAll(test *testing.T) {
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
	uc := usecases.NewUseCases(r)
	if uc == nil {
		test.Fatalf("an error was not expected when creating usecases")
	}

	expected_one := &internal.User{
		ID: "1",
		Email: "test1@test.com",
		Phone: "1234567890",
	}

	expected_two := &internal.User{
		ID: "2",
		Email: "test2@test.com",
		Phone: "0987654321",
	}
	
	test.Run("it should index an user slice", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(spReadAll)).
			WillReturnRows(sqlmock.NewRows(columns).				
				AddRow("1", "test1@test.com", "1234567890").
				AddRow("2", "test2@test.com", "0987654321"))

		m := controller.Users{}

		i, err := uc.IndexAll()
		if err != nil {
			t.Errorf("error was not expected while indexing data: %s", err)
		}

		err = mapstructure.Decode(i, &m)
		if err != nil {
			t.Errorf("error was not expected while decoding indexed data: %s", err)
		}

		fmt.Println(m)
		t.Log("successfully indexed")
	})

	test.Run("it should not read, db is empty", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(spReadAll)).
		WillReturnRows(
			sqlmock.NewRows(columns))
		
		users, err := uc.IndexAll()

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(users) != 0 {
			t.Fatalf("expected to be an empty list")
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

func TestModify(test *testing.T) {
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
		Password: "password",
	}

	// mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(spUpdate)).
		WithArgs(user.ID, user.Email, user.Phone, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	uc := usecases.NewUseCases(r)
	if uc == nil {
		test.Fatalf("an error was not expected when creating usecases")
	}

	test.Run("it should modify", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.ModifyUser(u)
		if err != nil {
			t.Fatalf("error was not expected while registering an user: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not modify, empty ID", func(t *testing.T) {
		u := &controller.User{
			ID:       "",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.ModifyUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not modify, empty User", func(t *testing.T) {
		u := &controller.User{}

		// Call the Create method with the user instance.
		err = uc.ModifyUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not modify, User is nil", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = uc.ModifyUser(nil)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should modify even if empty Email", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should modify even if empty phone", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.RegisterUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not modify, user email and phone missing", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.ModifyUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not modify, empty password", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "",
		}

		// Call the Create method with the user instance.
		err = uc.ModifyUser(u)
		if err == nil {
			t.Fatalf("error was expected")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}

func TestDestroy(test *testing.T) {
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
		Password: "password",
	}

	// mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(spDelete)).
		WithArgs(user.ID, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))


	uc := usecases.NewUseCases(r)
	if uc == nil {
		test.Fatalf("an error was not expected when creating usecases")
	}

	test.Run("it should destroy", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.DestroyUser(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})


	test.Run("it should destroy, besides email is missing", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.DestroyUser(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should destroy, besides phone is missing", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.DestroyUser(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not destroy, ID is missing", func(t *testing.T) {
		u := &controller.User{
			ID:       "",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = uc.DestroyUser(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not destroy, User is nil", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = r.Delete(nil)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not destroy, User is empty", func(t *testing.T) {
		u := &controller.User{}

		// Call the Create method with the user instance.
		err = uc.DestroyUser(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not destroy, Password is missing", func(t *testing.T) {
		u := &controller.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "",
		}

		// Call the Create method with the user instance.
		err = uc.DestroyUser(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}