package handlers

import (
	"database/sql"
	"exampleApi/consts"
	"exampleApi/helpers"
	"exampleApi/helpers/log"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitMiddlewares(r *gin.Engine, db *sql.DB, redisDb *redis.Client) {
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.Use(func(c *gin.Context) {
		c.Set("redis_db", redisDb)
		c.Next()
	})

	// set start time on every request to check response time
	r.Use(func(c *gin.Context) {
		startTime := time.Now()
		c.Set("startTime", startTime)
		c.Next()
	})
}

func authMiddleware(c *gin.Context) {
	const TOKEN_EXPIRED_ERROR = "token is expired"

	accessToken, err := c.Cookie(consts.ACCESS_TOKEN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no access token"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("no access token: %v", err.Error()))
		c.Abort()
		return
	}

	claims, err := helpers.ParseToken(accessToken, helpers.GetEnv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": TOKEN_EXPIRED_ERROR})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, fmt.Sprintf("can't parse token: %v", err.Error()))
		c.Abort()
		return
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": TOKEN_EXPIRED_ERROR})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, "can't extract user_id claim")
		c.Abort()
		return
	}

	c.Set("user_id", userId)

	c.Next()
}
