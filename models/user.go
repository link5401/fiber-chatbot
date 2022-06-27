package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id" gorm:"primaryKey" swaggerignore:"true"`
	Username  string    `gorm:"index" json:"username" validate:"required"`
	Password  string    `validate:"required,min=5"`
	Language  string    `json:"language" swaggerignore:"true"`
	CreatedAt time.Time `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at" swaggerignore:"true"`
	DeletedAt time.Time `json:"deleted_at" swaggerignore:"true"`
}

func HashPassword(password string) string {
	cost := 9 // Min: 4, Max: 31
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes)
}

func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
