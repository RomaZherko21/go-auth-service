package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func SetDbMiddleware(c *gin.Context, db *sql.DB) {
	c.Set("db", db)
	c.Next()
}
