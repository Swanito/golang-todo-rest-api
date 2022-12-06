package model

import (
	"time"
)

type NewUser struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"email,required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
	CreatedAt       time.Time
}

type User struct {
	Id            int
	Username      string
	Email         string
	Password      string
	CreatedAt     time.Time
	LoginAttempts int
	Activated     bool
	ActivatedAt   bool
}

type Login struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Email       string    `json:"email" validate:"required"`
	Username    string    `json:"username" validate:"required"`
	CreatedAt   time.Time `json:"createdAt" validate:"required"`
	AccessToken string    `json:"accessToken" validate:"required"`
}

type RegisterResponse struct {
	UserId int `json:"userId" validate:"required"`
}

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	Register(username string, password string, email string, createdAt time.Time) (int, error)
	GetLoginAttempts(email string) (int, error)
	AddInvalidLoginAttempt(email string) error
	AddLoginSession(id int, sessionDate time.Time) error
	ActivateUser(userEmail string) error
}
