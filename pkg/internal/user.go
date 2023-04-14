// Package internal contains bussines rules
package internal

// User struct is the model of the actor called user
type User struct {
	ID       string
	Email    string
	Phone    string
	Password string
}

// Users is an array type of User
type Users []*User
