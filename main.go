package main

import (
	"log"
	"net/http"
	"net"
	"fmt"
	"exampleApi/entities/user"
	"exampleApi/db"
)

type User struct {
    firstName string 
    lastName string 
    age int 
}

type UserI interface {
    getFullName() string
}


func (user User) getFullName() string {
    return user.firstName + " " + user.lastName
}

func (user User) getAge() int {
    return user.age
}


func some(u UserI, us *User) string{
	us.firstName = "KEKE"
	us.getFullName()

	return u.getFullName()
}

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

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server has been started on port 8000")

	if err := http.Serve(l, nil); err != nil {
		fmt.Println(err)
	}
}