package handlers

import (
	"database/sql"
	"net/http"
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

func authMiddleware(c *gin.Context) {
	_, ok := c.GetQuery("access_token")

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
