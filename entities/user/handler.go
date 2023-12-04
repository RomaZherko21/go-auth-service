package user

import (
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
	// var user User

	// err := json.NewDecoder(req.Body).Decode(&user)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// UserServiceInstance.CreateUser(db, &user)

	// log.Printf("User was created")
}
