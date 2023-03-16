package internal

type User struct {
	ID string
	Email string
	Phone string
	Password string
}

type Users []*User