package config

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"exampleApi/helpers"
)

func InitLogger() {
	logLevel := helpers.GetEnv("LOG_LEVEL")

	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	intNumber, err := strconv.Atoi(logLevel)
	if err != nil {
		log.Panic("Error converting string to int:", err)
		return
	}

	log.SetLevel(log.Level(uint32(intNumber)))
}
