package user

import "time"

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Nickname    string    `json:"nickname"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}
