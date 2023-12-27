package helpers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

func GetEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Panicf("Error loading .env file: %v", err)
	}

	fmt.Println("asdasdasd")

	return os.Getenv(key)
}

func CheckRequiredEnvs() {
	requiredEnvVars := []string{"SERVER_PORT", "DB_PORT", "DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "LOG_LEVEL", "REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "ACCESS_TOKEN_SECRET", "ACCESS_TOKEN_EXP_MIN", "REFRESH_TOKEN_SECRET", "REFRESH_TOKEN_EXP_MIN", "ACCESS_TOKEN_SECRET_PUBLICK", "ACCESS_TOKEN_SECRET_PRIVATE"}

	for _, envVar := range requiredEnvVars {
		if value, exists := os.LookupEnv(envVar); !exists || value == "" {
			log.Panic(fmt.Sprintf("Error: Environment variable %v is not set.", envVar))
		}
	}
}
