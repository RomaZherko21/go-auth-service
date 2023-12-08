package handlers

import (
	"database/sql"
	"exampleApi/helpers"
	"exampleApi/helpers/log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	authorization := c.GetHeader("authorization")

	tokenFields := strings.Fields(authorization)

	if len(tokenFields) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no access token"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, "no access token")
		c.Abort()
		return
	}

	tokenString := tokenFields[1]

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(helpers.GetEnv("ACCESS_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	c.Next()
}
