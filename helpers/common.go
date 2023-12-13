package helpers

import (
	"os"

	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	return os.Getenv(key)
}
