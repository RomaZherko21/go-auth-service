package user

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type UserService struct {
}

func (u *UserService) CreateUser(c *gin.Context, user *User) {
	db := c.MustGet("db").(*sql.DB)

	sqlStatement := `INSERT INTO users (email, password, nickname, phone_number, created_at)
	VALUES ($1, $2, $3, $4, $5);`

	_, err := db.Exec(sqlStatement, user.Email, user.Password, user.Nickname, user.PhoneNumber, user.CreatedAt)
	// _, err := db.Exec(sqlStatement, "user.Email", "user.Password", "user.Nickname", "user.PhoneNumber", "user.CreatedAt")
	if err != nil {
		panic(err)
	}
}

var UserServiceInstance = UserService{}
