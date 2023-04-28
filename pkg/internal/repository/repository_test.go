// Package repository_test contains test for repository
package repository_test

import (
	"dall06/go-cleanapi/pkg/internal"
	"dall06/go-cleanapi/pkg/internal/repository"
	"database/sql"
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
	spLogin   = "CALL `go_cleanapi`.`sp_login_user`(?, ?, ?);"
)

func TestLogin(test *testing.T) {
	dbUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	dbUserTwo := &internal.User{
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}
	dbUserThree := &internal.User{
		Email:    "",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	rowsSetOne := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		&dbUserOne.ID,
	)

	rowsSetOneTwo := sqlmock.NewRows([]string{
		"id_user",
	}).AddRow(
		&dbUserOne.ID,
	)

	rowsSetTwo := sqlmock.NewRows([]string{
		"id_user",
	})

	inputUserOne := &internal.User{
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	inputUserTwo := &internal.User{
		Email:    "",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserThree := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserFour := &internal.User{
		Email:    "test2@test.com",
		Password: "12345pAsSWORd*",
	}

	inputUserFive := &internal.User{
		Phone:    "+8812324524",
		Password: "12345pAsSWORd*",
	}

	expectedOne := &internal.User{
		ID: "im an id",
	}

	successfulCases := []struct {
		name     string
		input    *internal.User
		rows     *sqlmock.Rows
		dbUser   *internal.User
		expected *internal.User
	}{
		{
			name:     "it should login (mocked), whit email",
			input:    inputUserOne,
			rows:     rowsSetOne,
			dbUser:   dbUserTwo,
			expected: expectedOne,
		},
		{
			name:     "it should login (mocked), with phone",
			input:    inputUserTwo,
			rows:     rowsSetOneTwo,
			dbUser:   dbUserThree,
			expected: expectedOne,
		},
	}

	failedCases := []struct {
		name     string
		input    *internal.User
		rows     *sqlmock.Rows
		dbUser   *internal.User
		expected *internal.User
	}{
		{
			name:     "it should not login (mocked), empty user",
			input:    nil,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not login (mocked), empty db",
			input:    inputUserTwo,
			rows:     rowsSetTwo,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not login (mocked), no id found",
			input:    inputUserThree,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not login (mocked), no password found",
			input:    inputUserThree,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not login (mocked), no email or phone",
			input:    inputUserThree,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not login (mocked), both params found",
			input:    inputUserThree,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not login (mocked), email not found",
			input:    inputUserFour,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not login (mocked), phone not found",
			input:    inputUserFive,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spLogin)).WithArgs(
				tc.dbUser.Email,
				tc.dbUser.Phone,
				tc.dbUser.Password,
			).WillReturnRows(tc.rows)

			assert.NoError(t, err)

			r := repository.NewRepository(db)
			res, err := r.Login(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, res)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spLogin)).WithArgs(
				tc.dbUser.Email,
				tc.dbUser.Phone,
				tc.dbUser.Password,
			).WillReturnRows(tc.rows)
			assert.NoError(t, err)

			r := repository.NewRepository(db)
			res, err := r.Login(tc.input)
			assert.Error(t, err)
			assert.NotEqual(t, tc.expected, res)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestCreate(test *testing.T) {
	dbUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	dbUserTwo := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	inputUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserTwo := &internal.User{
		ID:       "",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserThree := &internal.User{
		ID:       "im an id",
		Email:    "",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserFour := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	inputUserFive := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "",
	}

	inputUserSix := &internal.User{
		ID:       "im an id but wrong",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	successfulCases := []struct {
		name   string
		input  *internal.User
		dbUser *internal.User
	}{
		{
			name:   "it should create an user (mocked)",
			input:  inputUserOne,
			dbUser: dbUserOne,
		},
		{
			name:   "it should create an user (mocked) besides empty phone",
			input:  inputUserFour,
			dbUser: dbUserTwo,
		},
	}

	failedCases := []struct {
		name   string
		input  *internal.User
		dbUser *internal.User
	}{
		{
			name:   "it should not create an user (mocked), empty id",
			input:  inputUserTwo,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not create an user (mocked), empty email",
			input:  inputUserThree,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not create an user (mocked), empty password",
			input:  inputUserFive,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not create an user (mocked), nil user",
			input:  nil,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not create an user (mocked), internal error",
			input:  inputUserSix,
			dbUser: dbUserOne,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spCreate)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Email,
				&tc.dbUser.Phone,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 0))

			r := repository.NewRepository(db)
			err = r.Create(tc.input)
			assert.NoError(t, err)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spCreate)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Email,
				&tc.dbUser.Phone,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 0))

			r := repository.NewRepository(db)
			err = r.Create(tc.input)
			assert.Error(t, err)

			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestRead(test *testing.T) {
	dbUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	rowsSetOne := sqlmock.NewRows([]string{
		"id_user",
		"user_email",
		"user_phone",
	}).AddRow(
		&dbUserOne.ID,
		&dbUserOne.Email,
		&dbUserOne.Phone,
	)

	rowsSetTwo := sqlmock.NewRows([]string{
		"id_user",
		"user_email",
		"user_phone",
	})

	inputUserOne := &internal.User{
		ID: "im an id",
	}

	inputUserTwo := &internal.User{
		ID: "",
	}

	inputUserThree := &internal.User{
		ID: "im an id two",
	}

	expectedOne := &internal.User{
		ID:    "im an id",
		Email: "test@test.com",
		Phone: "+7812324524",
	}

	successfulCases := []struct {
		name     string
		input    *internal.User
		rows     *sqlmock.Rows
		dbUser   *internal.User
		expected *internal.User
	}{
		{
			name:     "it should read an user (mocked)",
			input:    inputUserOne,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
	}

	failedCases := []struct {
		name     string
		input    *internal.User
		rows     *sqlmock.Rows
		dbUser   *internal.User
		expected *internal.User
	}{
		{
			name:     "it should not read an user (mocked), empty id",
			input:    inputUserTwo,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not read an user (mocked), empty user",
			input:    nil,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not read an user (mocked), empty db",
			input:    inputUserTwo,
			rows:     rowsSetTwo,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
		{
			name:     "it should not read an user (mocked), no id found",
			input:    inputUserThree,
			rows:     rowsSetOne,
			dbUser:   dbUserOne,
			expected: expectedOne,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spRead)).WithArgs(
				&tc.dbUser.ID,
			).WillReturnRows(tc.rows)

			r := repository.NewRepository(db)
			res, err := r.Read(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, res)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spRead)).WithArgs(
				&tc.dbUser.ID,
			).WillReturnRows(tc.rows)

			r := repository.NewRepository(db)
			res, err := r.Read(tc.input)
			assert.Error(t, err)
			assert.NotEqual(t, tc.expected, res, "")
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestReadAll(test *testing.T) {
	dbUserOne := &internal.User{
		ID:    "im an id",
		Email: "test@test.com",
		Phone: "+7812324524",
	}
	dbUserTwo := &internal.User{
		ID:    "im an id 2",
		Email: "test2@test.com",
		Phone: "+8812324524",
	}
	dbUserThree := &internal.User{
		ID:    "im an id 3",
		Email: "test2@test.com",
		Phone: "+8812324524",
	}

	rowsSetOne := sqlmock.NewRows([]string{
		"id_user",
		"user_email",
		"user_phone",
	}).AddRow(
		&dbUserOne.ID,
		&dbUserOne.Email,
		&dbUserOne.Phone,
	).AddRow(
		&dbUserTwo.ID,
		&dbUserTwo.Email,
		&dbUserTwo.Phone,
	)

	rowsSetTwo := sqlmock.NewRows([]string{"id_user", "user_email", "user_phone"})

	expectedOnes := make(internal.Users, 0)
	expectedOnes = append(expectedOnes, dbUserOne)
	expectedOnes = append(expectedOnes, dbUserTwo)

	expectedOnesTwo := make(internal.Users, 0)

	expectedOnesThree := make(internal.Users, 0)
	expectedOnesThree = append(expectedOnesThree, dbUserOne)
	expectedOnesThree = append(expectedOnesThree, dbUserThree)

	successfulCases := []struct {
		name     string
		rows     *sqlmock.Rows
		expected internal.Users
	}{
		{
			name:     "it should read many users (mocked)",
			rows:     rowsSetOne,
			expected: expectedOnes,
		},
		{
			name:     "it should read many users (mocked), but empty db",
			rows:     rowsSetTwo,
			expected: expectedOnesTwo,
		},
	}

	failedCases := []struct {
		name     string
		rows     *sqlmock.Rows
		dbUser   internal.Users
		expected internal.Users
	}{
		{
			name:     "it should not read many users (mocked), wrong values returned",
			rows:     rowsSetOne,
			expected: expectedOnesThree,
		},
	}

	for _, tc := range successfulCases {
		tc := tc
		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectQuery(regexp.QuoteMeta(spReadAll)).WillReturnRows(tc.rows)

			r := repository.NewRepository(db)
			res, err := r.ReadAll()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, res)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc
		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)
			m.ExpectQuery(regexp.QuoteMeta(spReadAll)).WillReturnRows(tc.rows)

			r := repository.NewRepository(db)
			res, err := r.ReadAll()
			assert.NoError(t, err)
			assert.NotEqual(t, tc.expected, res)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestUpdate(test *testing.T) {
	dbUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	dbUserTwo := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	dbUserThree := &internal.User{
		ID:       "im an id",
		Email:    "",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserTwo := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "",
		Password: "12345pAsSWORd*",
	}

	inputUserThree := &internal.User{
		ID:       "im an id",
		Email:    "",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserFour := &internal.User{
		ID:       "",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserFive := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "",
	}

	inputUserSix := &internal.User{
		ID:       "im an id but wrong",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	successfulCases := []struct {
		name   string
		input  *internal.User
		dbUser *internal.User
	}{
		{
			name:   "it should update an user (mocked)",
			input:  inputUserOne,
			dbUser: dbUserOne,
		},
		{
			name:   "it should update an user (mocked) besides empty phone",
			input:  inputUserTwo,
			dbUser: dbUserTwo,
		},
		{
			name:   "it should update an user (mocked) besides empty email",
			input:  inputUserThree,
			dbUser: dbUserThree,
		},
	}

	failedCases := []struct {
		name   string
		input  *internal.User
		dbUser *internal.User
	}{
		{
			name:   "it should not update an user (mocked), empty id",
			input:  inputUserFour,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not update an user (mocked), empty password",
			input:  inputUserFive,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not update an user (mocked), nil user",
			input:  nil,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not update an user (mocked), id not found",
			input:  inputUserSix,
			dbUser: dbUserOne,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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

			r := repository.NewRepository(db)
			err = r.Update(tc.input)
			assert.NoError(t, err)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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
			).WillReturnResult(sqlmock.NewResult(0, 0))

			r := repository.NewRepository(db)
			err = r.Update(tc.input)
			assert.Error(t, err)
			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}

func TestDelete(test *testing.T) {
	dbUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserOne := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserFour := &internal.User{
		ID:       "",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	inputUserFive := &internal.User{
		ID:       "im an id",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "",
	}

	inputUserSix := &internal.User{
		ID:       "im an id but wrong",
		Email:    "test@test.com",
		Phone:    "+7812324524",
		Password: "12345pAsSWORd*",
	}

	successfulCases := []struct {
		name   string
		input  *internal.User
		dbUser *internal.User
	}{
		{
			name:   "it should delete an user (mocked)",
			input:  inputUserOne,
			dbUser: dbUserOne,
		},
	}

	failedCases := []struct {
		name   string
		input  *internal.User
		dbUser *internal.User
	}{
		{
			name:   "it should not delete an user (mocked), empty id",
			input:  inputUserFour,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not delete an user (mocked), empty password",
			input:  inputUserFive,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not delete an user (mocked), nil user",
			input:  nil,
			dbUser: dbUserOne,
		},
		{
			name:   "it should not delete an user (mocked), id not found",
			input:  inputUserSix,
			dbUser: dbUserOne,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spDelete)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 1))

			r := repository.NewRepository(db)
			err = r.Delete(tc.input)
			assert.NoError(t, err)

			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, m, err := sqlmock.New()
			if err != nil {
				test.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			assert.NoError(t, err)

			m.ExpectExec(regexp.QuoteMeta(spDelete)).WithArgs(
				&tc.dbUser.ID,
				&tc.dbUser.Password,
			).WillReturnResult(sqlmock.NewResult(0, 1))

			r := repository.NewRepository(db)
			err = r.Delete(tc.input)
			assert.Error(t, err)

			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

			m.ExpectClose().WillReturnError(sql.ErrConnDone) // expect a call to Close() but return an error to indicate that it was not expected

		})
	}
}
