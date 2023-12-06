package user

import "time"

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,min=8"`
	Nickname    string    `json:"nickname" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required,min=10,max=15"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
}
