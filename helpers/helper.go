package helpers

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	return os.Getenv(key)
}

func HttpSend(data interface{}, w http.ResponseWriter) {
	jData, err := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		errData, _ := json.Marshal(err)
		w.Write(errData)
	}

	w.Write(jData)
}
