package user

import "time"

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,password"`
	Nickname    string    `json:"nickname" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required,min=10,max=15,phone"`
	CreatedAt   time.Time `json:"created_at"`
}
