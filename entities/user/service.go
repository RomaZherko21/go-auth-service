package user

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserService struct {
}

func (u *UserService) GetUserPassword(c *gin.Context, user *User) (User, error) {
	result := User{}

	db := c.MustGet("db").(*sql.DB)

	sqlStatement := `SELECT password, id
	FROM users
	WHERE users.email=$1;`

	err := db.QueryRow(sqlStatement, user.Email).Scan(&result.Password, &result.ID)

	fmt.Println("EHHE", result)

	return result, err
}

func (u *UserService) CreateUser(c *gin.Context, user *User) error {
	db := c.MustGet("db").(*sql.DB)

	sqlStatement := `INSERT INTO users (email, password, nickname, phone_number)
	VALUES ($1, $2, $3, $4);`

	_, err := db.Exec(sqlStatement, user.Email, user.Password, user.Nickname, user.PhoneNumber)

	return err
}

var UserServiceInstance = UserService{}
