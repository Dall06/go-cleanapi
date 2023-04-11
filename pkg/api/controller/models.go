package controller

type User struct {
	ID       string `json:"uid"`
	Email    string `json:"mail"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Users []User

type PostRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required"`
}

type PutRequest struct {
	Email    string `json:"email" validate:"email"`
	Phone    string `json:"phone" validate:"omitempty"`
	Password string `json:"password" validate:"required"`
}

type DeleteRequest struct {
	Password string `json:"password" validate:"required"`
}
