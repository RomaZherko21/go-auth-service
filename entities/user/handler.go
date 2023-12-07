package user

import (
	"exampleApi/helpers"
	"exampleApi/helpers/log"
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
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	validationResult := helpers.Validate(&user)

	if !validationResult.OK {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationResult.Errors})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, "validation error")
		return
	}

	if err := UserServiceInstance.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	log.HttpLog(c, log.Info, http.StatusOK, "User created successfully")
}
