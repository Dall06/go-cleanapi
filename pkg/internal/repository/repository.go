// Package repository provides the methods that intercat with data source
package repository

import (
	"dall06/go-cleanapi/pkg/internal"
	"database/sql"
	"fmt"
)

const (
	spCreate  = "CALL `go_cleanapi`.`sp_create_user`(?, ?, ?, ?);"
	spRead    = "CALL `go_cleanapi`.`sp_read_user`(?);"
	spReadAll = "CALL `go_cleanapi`.`sp_read_users`();"
	spUpdate  = "CALL `go_cleanapi`.`sp_update_user`(?, ?, ?, ?);"
	spDelete  = "CALL `go_cleanapi`.`sp_delete_user`(?, ?);"
	spLogin   = "CALL `go_cleanapi`.`sp_login_user`(?, ?, ?);"
)

// Repository is an interface that extends the repository
type Repository interface {
	Create(user *internal.User) error
	Read(user *internal.User) (*internal.User, error)
	ReadAll() (internal.Users, error)
	Update(user *internal.User) error
	Delete(user *internal.User) error
	Login(user *internal.User) (*internal.User, error)
}

var _ Repository = (*repository)(nil)

type repository struct {
	dbConn *sql.DB
}

// NewRepository is a constructor for a repository
func NewRepository(db *sql.DB) Repository {
	return &repository{
		dbConn: db,
	}
}

func (r *repository) Login(user *internal.User) (*internal.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is required")
	}
	if user.Email == "" && user.Phone == "" {
		return nil, fmt.Errorf("data is required")
	}
	if user.Email != "" && user.Phone != "" {
		return nil, fmt.Errorf("only one parameter is required")
	}
	if user.Password == "" {
		return nil, fmt.Errorf("password is required")
	}
	// Add more validation checks as needed.
	row := r.dbConn.QueryRow(spLogin,
		user.Email,
		user.Phone,
		user.Password)
	if row == nil {
		empty := &internal.User{}
		return empty, nil
	}

	u := &internal.User{}

	err := row.Scan(&u.ID)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *repository) Create(user *internal.User) error {
	if user == nil {
		return fmt.Errorf("user is empty")
	}
	if user.ID == "" {
		return fmt.Errorf("ID is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	// Add more validation checks as needed.

	_, err := r.dbConn.Exec(spCreate,
		user.ID,
		user.Email,
		user.Phone,
		user.Password)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return nil
}

func (r repository) Read(user *internal.User) (*internal.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is required")
	}
	if user.ID == "" {
		return nil, fmt.Errorf("ID is required")
	}

	u := &internal.User{}

	row := r.dbConn.QueryRow(spRead, user.ID)
	if row == nil {
		empty := &internal.User{}
		return empty, nil
	}

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Phone)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *repository) ReadAll() (internal.Users, error) {
	rows, err := r.dbConn.Query(spReadAll)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			err = cerr
		}
	}()

	users := make(internal.Users, 0) // allocate slice

	for rows.Next() {
		user := &internal.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Phone,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r repository) Update(user *internal.User) error {
	if user == nil {
		return fmt.Errorf("user is required")
	}
	if user.ID == "" {
		return fmt.Errorf("ID is resquired")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	if user.Email == "" && user.Phone == "" {
		return fmt.Errorf("user data is required")
	}

	res, err := r.dbConn.Exec(spUpdate,
		user.ID,
		user.Email,
		user.Phone,
		user.Password)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to obtain rows affected: %v", err)
	}

	if affected == 0 {
		return fmt.Errorf("user not updated")
	}

	return nil
}

func (r repository) Delete(user *internal.User) error {
	if user == nil {
		return fmt.Errorf("user is required")
	}
	if user.ID == "" {
		return fmt.Errorf("ID is required")
	}
	if user.Password == "" {
		return fmt.Errorf("password is required")
	}

	res, err := r.dbConn.Exec(spDelete, user.ID, user.Password)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to obtain rows affected: %v", err)
	}

	if affected == 0 {
		return fmt.Errorf("user not deleted")
	}

	return nil
}
