package user

import (
	"database/sql"
	"log"
)

type UserService struct {

}

func (u *UserService) CreateUser(db *sql.DB, user *User) {
	log.Printf("User was created", user)

	sqlStatement := `INSERT INTO users (email, password, nickname, phone_number, created_at)
	VALUES ($1, $2, $3, $4, $5);`

	_, err := db.Exec(sqlStatement, user.Email, user.Password, user.Nickname, user.PhoneNumber, user.CreatedAt)
	// _, err := db.Exec(sqlStatement, "user.Email", "user.Password", "user.Nickname", "user.PhoneNumber", "user.CreatedAt")
	if err != nil {
		panic(err)
	}
}

var UserServiceInstance = UserService{}