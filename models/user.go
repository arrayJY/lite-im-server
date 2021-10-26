package models

import "gorm.io/gorm"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type CreateUserRequestBody struct {
	User
}

type UserModel struct {
	gorm.Model
	User
}
