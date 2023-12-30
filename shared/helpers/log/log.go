package log

import (
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type HttpLogInfo struct{}

type LogLevel int

const (
	Panic LogLevel = iota
	Error
	Warn
	Info
)

func HttpLog(c *gin.Context, level LogLevel, httpStatus int, message string) {
	startTime := c.MustGet("startTime")

	fields := log.Fields{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"statusCode":  httpStatus,
		"client_ip":   c.ClientIP(),
		"duration_ms": time.Since(startTime.(time.Time)) / 1000000,
		"user_agent":  c.Request.UserAgent(),
	}

	switch level {
	case Panic:
		log.WithFields(fields).Panic(message)
	case Error:
		log.WithFields(fields).Error(message)
	case Warn:
		log.WithFields(fields).Warn(message)
	case Info:
		log.WithFields(fields).Info(message)
	default:
		log.WithFields(fields).Info(message)
	}
}
