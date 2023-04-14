package controller

// User is a struct model for users interaction in controller layer
type User struct {
	ID       string `json:"uid"`
	Email    string `json:"mail"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Users is a struct model for a slice of users interaction in controller layer
type Users []User

// PostRequest is a struct model for post requests in controller layer
type PostRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required"`
}

// PutRequest is a struct model for put requests in controller layer
type PutRequest struct {
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required"`
}

// DeleteRequest is a struct model for delete requests in controller layer
type DeleteRequest struct {
	Password string `json:"password" validate:"required"`
}
