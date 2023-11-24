package db

import (
	_ "github.com/lib/pq"

  "database/sql"
  "fmt"
  
  "exampleApi/helpers"
)

var DB_PORT = helpers.GetEnv("DB_PORT")
var DB_HOST = helpers.GetEnv("DB_HOST")
var DB_NAME = helpers.GetEnv("DB_NAME")
var DB_USER = helpers.GetEnv("DB_USER")
var DB_PASSWORD = helpers.GetEnv("DB_PASSWORD")


func ConnectDb(){
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB was successfully connected!")
}