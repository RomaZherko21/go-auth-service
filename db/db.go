package db

import (
	_ "github.com/lib/pq"

  "database/sql"
  "fmt"
  "os"	
  "os/exec"
  
  "exampleApi/helpers"
)

var DB_PORT = helpers.GetEnv("DB_PORT")
var DB_HOST = helpers.GetEnv("DB_HOST")
var DB_NAME = helpers.GetEnv("DB_NAME")
var DB_USER = helpers.GetEnv("DB_USER")
var DB_PASSWORD = helpers.GetEnv("DB_PASSWORD")

func runMigrations() {
	psqlInfo := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

    cmd := exec.Command("migrate", "-path", "db/migrations", "-database", psqlInfo, "up")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err := cmd.Run()
    if err != nil {
        panic(err)
    }
}

func ConnectDb(){
	runMigrations()

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