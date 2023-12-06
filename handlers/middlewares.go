package handlers

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

func SetDbMiddleware(c *gin.Context, db *sql.DB) {
	c.Set("db", db)
	c.Next()
}

func SetStartTime(c *gin.Context) {
	startTime := time.Now()
	c.Set("startTime", startTime)
	c.Next()
}
