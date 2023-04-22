package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id" db:"id"`
	Nickname     string    `json:"nickname" db:"nickname"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	ConfirmToken string    `json:"confirm_token" db:"confirm_token"`
	Confirmed    bool      `json:"confirmed" db:"confirmed"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
