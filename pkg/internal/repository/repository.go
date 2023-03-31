package repository

import (
	"dall06/go-cleanapi/pkg/internal"
	"database/sql"
	"fmt"
)

type Repository interface {
	Create(user *internal.User) error
	Read(user *internal.User) (*internal.User, error)
	ReadAll() (internal.Users, error)
	Update(user *internal.User) error
	Delete(user *internal.User) error
}

var _ Repository = (*repository)(nil)

type repository struct {
	dbConn *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		dbConn: db,
	}
}

func (r *repository) Create(user *internal.User) error {
	if user == nil {
		return fmt.Errorf("user is empty")
	}
	if user.ID == "" {
		return fmt.Errorf("user ID is empty")
	}
	if user.Email == "" {
		return fmt.Errorf("user email is empty")
	}
	if user.Password == "" {
		return fmt.Errorf("user password is empty")
	}
	// Add more validation checks as needed.

	res, err := r.dbConn.Exec(spCreate,
		user.ID,
		user.Email,
		user.Phone,
		user.Password)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to obtain rows affected: %v", err)
	}

	if lastId == 0 {
		return fmt.Errorf("user not created")
	}

	return nil
}

func (r repository) Read(user *internal.User) (*internal.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is empty")
	}
	if user.ID == "" {
		return nil, fmt.Errorf("user ID is empty")
	}

	u := &internal.User{}

	err := r.dbConn.QueryRow(spRead, user.ID).Scan(
		&u.ID,
		&u.Email,
		&u.Phone)
	if err == sql.ErrNoRows {
		// double nil because i consider that no rows found is more as a goo request but with an empty result
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
	defer rows.Close()

	users := make(internal.Users, 0, 100) // allocate slice with initial capacity of 100

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

	return users, nil
}

func (r repository) Update(user *internal.User) error {
	if user == nil {
		return fmt.Errorf("user is empty")
	}
	if user.ID == "" {
		return fmt.Errorf("user ID is empty")
	}
	if user.Password == "" {
		return fmt.Errorf("user password is empty")
	}
	if user.Email == "" && user.Phone == "" {
		return fmt.Errorf("user data is empty")
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
		return fmt.Errorf("user not created")
	}

	return nil
}

func (r repository) Delete(user *internal.User) error {
	if user == nil {
		return fmt.Errorf("user is empty")
	}
	if user.ID == "" {
		return fmt.Errorf("user ID is empty")
	}
	if user.Password == "" {
		return fmt.Errorf("user email is password")
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
		return fmt.Errorf("user not created")
	}

	return nil
}
