package main

import (
	"log"
	"net/http"
	"net"
	"fmt"
	"exampleApi/entities/user"
	"exampleApi/db"
	"exampleApi/helpers"
)

func handlers(){
	const (
		GET = "GET"
		POST = "POST"
		PUT = "PUT"
		DELETE = "DELETE"
	)

	const (
		USERS = "/users"
		POSTS = "/posts"
	)

	http.HandleFunc(USERS, func(w http.ResponseWriter, req *http.Request){

		switch req.URL.Path {
		case USERS:
			switch req.Method {
			case GET:
				user.GetUser(w, req)
			case POST:
				user.PostUser(w, req)
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
	db.ConnectDb()

	handlers()

	var SERVER_PORT = helpers.GetEnv("SERVER_PORT")

	l, err := net.Listen("tcp", ":" + SERVER_PORT)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Server has been started on port %v", SERVER_PORT)

	if err := http.Serve(l, nil); err != nil {
		fmt.Println(err)
	}
}