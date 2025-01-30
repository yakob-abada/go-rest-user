package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	FirstName string     `json:"first_name" gorm:"not null"`
	LastName  string     `json:"last_name" gorm:"not null"`
	Nickname  string     `json:"nickname" gorm:"not null"`
	Email     string     `json:"email" gorm:"unique;not null"` // UNIQUE CONSTRAINT
	Password  string     `json:"password" gorm:"not null"`
	Country   string     `json:"country" gorm:"not null"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Country   string `json:"country"`
}
