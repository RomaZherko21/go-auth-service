package handlers

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
)

func InitMiddlewares(r *gin.Engine, db *sql.DB) {
	r.Use(func(c *gin.Context) {
		setDbMiddleware(c, db)
	})

	r.Use(func(c *gin.Context) {
		setStartTime(c)
	})
}

func setDbMiddleware(c *gin.Context, db *sql.DB) {
	c.Set("db", db)
	c.Next()
}

func setStartTime(c *gin.Context) {
	startTime := time.Now()
	c.Set("startTime", startTime)
	c.Next()
}
