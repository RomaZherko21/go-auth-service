package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"database/sql"
	"net"
	"net/http"

	"exampleApi/config"
	"exampleApi/db"
	"exampleApi/entities/user"
	"exampleApi/helpers"
)

func init() {
	config.InitLogger()

}

func handlers(db *sql.DB) {
	const (
		GET    = "GET"
		POST   = "POST"
		PUT    = "PUT"
		DELETE = "DELETE"
	)

	const (
		USERS = "/users"
		POSTS = "/posts"
	)

	http.HandleFunc(USERS, func(w http.ResponseWriter, req *http.Request) {

		switch req.URL.Path {
		case USERS:
			switch req.Method {
			// case GET:
			// 	user.GetUser(w, req, db)
			case POST:
				user.CreateUser(w, req, db)
			case PUT:
				log.Println("three")
			case DELETE:
				log.Println("three")
			}
		case POSTS:
			log.Println("two")
		}
	})
}

func main() {
	dataBase := db.ConnectDb()

	defer dataBase.Close()

	handlers(dataBase)

	var SERVER_PORT = helpers.GetEnv("SERVER_PORT")

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", SERVER_PORT))
	if err != nil {
		log.Panic(err)
	}

	log.Infof("Server has been started on port %v", SERVER_PORT)

	if err := http.Serve(l, nil); err != nil {
		log.Panic(err)
	}
}
