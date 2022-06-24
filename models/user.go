package models

import "time"

type User struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	Username  string     `gorm:"index" json:"username" validate:"required"`
	Password  string     `json:"-" validate:"required,min=5" `
	Language  string     `json:"language"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}
