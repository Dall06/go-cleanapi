package controller

type User struct {
	ID       string `json:"uid"`
	Email    string `json:"mail"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Users []User
