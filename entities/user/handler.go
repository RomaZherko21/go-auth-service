package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// var users = make(map[int]User)

// func GetUser(w http.ResponseWriter, req *http.Request, db *sql.DB) {
// 	userId := req.URL.Query().Get("id")

// 	if userId == "" {
// 		log.Printf("Don't have query parameter: userId")
// 		helpers.HttpSend(types.HttpMessage{Message: "Don't have query parameter: userId"}, w)
// 		return
// 	}

// 	intUserId, _ := strconv.Atoi(userId)

// 	user, ok := users[intUserId]
// 	if !ok {
// 		message := fmt.Sprintf("User with id %v does not exits!", userId)
// 		log.Printf(message)
// 		helpers.HttpSend(types.HttpMessage{Message: message}, w)
// 		return
// 	}

// 	helpers.HttpSend(user, w)
// }

func CreateUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Received email:", user.Email)
	fmt.Println("Received name:", user.Nickname)

	UserServiceInstance.CreateUser(c, &user)

	c.JSON(http.StatusOK, gin.H{"message": "JSON processed successfully"})

	// log.Printf("User was created")
}
