package handlers

import (
	"database/sql"
	"exampleApi/consts"
	"exampleApi/helpers"
	"exampleApi/helpers/log"
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

	r.Use(func(c *gin.Context) {
		setStartTime(c)
	})
}

func setStartTime(c *gin.Context) {
	startTime := time.Now()
	c.Set("startTime", startTime)
	c.Next()
}

func authMiddleware(c *gin.Context) {
	accessToken, err := c.Cookie(consts.ACCESS_TOKEN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	claims, err := helpers.ParseToken(accessToken, helpers.GetEnv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, "no user_id in token")
		c.Abort()
		return
	}

	c.Set("user_id", userId)

	c.Next()
}
