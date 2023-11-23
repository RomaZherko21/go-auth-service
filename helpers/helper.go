package helpers

import (
	"github.com/joho/godotenv"
	"encoding/json"
	"net/http"
	"log"
	"os"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")
  
	if err != nil {
	  log.Fatalf("Error loading .env file")
	}
  
	return os.Getenv(key)
  }

func HttpSend(data interface{}, w  http.ResponseWriter) {
	jData, err := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		errData, _ := json.Marshal(err)
		w.Write(errData)
	}

	w.Write(jData)
	return
}