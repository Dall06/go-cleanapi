package repository_test

import (
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/pkg/internal/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

const (
	spCreate  = "CALL `go_cleanapi`.`sp_create_user`(?, ?, ?, ?);"
	spRead    = "CALL `go_cleanapi`.`sp_read_user`(?);"
	spReadAll = "CALL `go_cleanapi`.`sp_read_users`();"
	spUpdate  = "CALL `go_cleanapi`.`sp_update_user`(?, ?, ?, ?);"
	spDelete  = "CALL `go_cleanapi`.`sp_delete_user`(?, ?);"
)

func TestCreate(test *testing.T) {
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

	test.Run("it should create", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = r.Create(user)
		if err != nil {
			t.Fatalf("failed to create user: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, ID is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Create(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, user is nil", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = r.Create(nil)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, User is empty", func(t *testing.T) {
		u := &internal.User{}

		// Call the Create method with the user instance.
		err = r.Create(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, Password is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "",
		}

		// Call the Create method with the user instance.
		err = r.Create(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not create, Email is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Create(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should create besides that Phone is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Create(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}

func TestRead(test *testing.T) {
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

	columns := []string{"id", "email", "phone"}
	mock.ExpectQuery(regexp.QuoteMeta(spRead)).
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(user.ID, user.Email, user.Phone))

	test.Run("it should read an user", func(t *testing.T) {
		// now we execute our method
		res, err := r.Read(user)
		if err != nil {
			t.Errorf("error when trying to read an user: %s", err)
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("error when trying to read an user: %s", err)
		}

		expected := &internal.User{
			ID:    "1",
			Email: "email",
			Phone: "phone",
		}

		assert.Equal(t, expected, res)
	})

	test.Run("it should not read, ID is empty", func(t *testing.T) {
		u := internal.User{
			ID: "",
		}
		// now we execute our method
		_, err = r.Read(&u)
		if err == nil {
			t.Errorf("error no generated")
		}
		t.Log("user not read")
	})

	test.Run("it should not read, user is empty", func(t *testing.T) {
		u := internal.User{}
		// now we execute our method
		_, err = r.Read(&u)
		if err == nil {
			t.Errorf("error no generated")
		}
		t.Log("user not read")
	})

	test.Run("it should not read, user is nil", func(t *testing.T) {
		// now we execute our method
		_, err = r.Read(nil)
		if err == nil {
			t.Errorf("error no generated")
		}
		t.Log("user not read")
	})

	test.Run("it should not read, ID does not exist", func(t *testing.T) {
		u := &internal.User{
			ID: "3",
		}
		// now we execute our method
		_, err = r.Read(u)
		if err == nil {
			t.Errorf("error no generated")
		}
		t.Log("user not read")
	})
}

func TestReadAll(test *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := repository.NewRepository(db)
	if r == nil {
		test.Fatalf("an error was not expected when creating repository")
	}

	// set up the mock database to expect a query and return the expected rows
	columns := []string{"id", "email", "phone"}

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
	

	test.Run("it should read all", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(spReadAll)).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow("1", "test1@test.com", "1234567890").
				AddRow("2", "test2@test.com", "0987654321"))
		// now we execute our method
		users, err := r.ReadAll()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(users) == 0 {
			t.Fatalf("expected at least one user, got %d", len(users))
		}

		for _, u := range users {
			if u.ID == "1" {
				assert.Equal(t, expected_one, u)
			}

			if u.ID == "2" {
				assert.Equal(t, expected_two, u)
			}
		}

		// check that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled expectations: %s", err)
		}
	})

	test.Run("it should not read, db is empty", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(spReadAll)).
		WillReturnRows(
			sqlmock.NewRows(columns))
		// now we execute our method
		users, err := r.ReadAll()

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

		// check that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unfulfilled expectations: %s", err)
		}
	})
}

func TestUpdate(test *testing.T) {
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

	test.Run("it should update", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = r.Update(user)
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should update, besides email is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Update(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should update, besides phone is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Update(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not update, ID is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Update(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not update, User is nil", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = r.Update(nil)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not update, User is empty", func(t *testing.T) {
		u := &internal.User{}

		// Call the Create method with the user instance.
		err = r.Update(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not update, Password is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "",
		}

		// Call the Create method with the user instance.
		err = r.Update(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not update, user email and phone is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Update(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
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
		Password: "password",
	}

	// mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(spDelete)).
		WithArgs(user.ID, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	test.Run("it should delete", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = r.Delete(user)
		if err != nil {
			t.Fatalf("failed to update user: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should delete, besides email is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Delete(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should delete, besides phone is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Delete(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not delete, ID is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Delete(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not delete, User is nil", func(t *testing.T) {
		// Call the Create method with the user instance.
		err = r.Delete(nil)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not delete, User is empty", func(t *testing.T) {
		u := &internal.User{}

		// Call the Create method with the user instance.
		err = r.Delete(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should not delete, Password is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "johndoe@example.com",
			Phone:    "+1234567890",
			Password: "",
		}

		// Call the Create method with the user instance.
		err = r.Delete(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})

	test.Run("it should delete besides user email and phone is missing", func(t *testing.T) {
		u := &internal.User{
			ID:       "1",
			Email:    "",
			Phone:    "",
			Password: "password",
		}

		// Call the Create method with the user instance.
		err = r.Update(u)
		if err == nil {
			t.Fatalf("expecting an error")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
