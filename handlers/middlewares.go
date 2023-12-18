package handlers

import (
	"context"
	"database/sql"
	"exampleApi/helpers"
	"exampleApi/helpers/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitMiddlewares(r *gin.Engine, db *sql.DB, redisDb *redis.Client) {
	r.Use(func(c *gin.Context) {
		setDbMiddleware(c, db)
	})

	r.Use(func(c *gin.Context) {
		setRedisDbMiddleware(c, redisDb)
	})

	r.Use(func(c *gin.Context) {
		setStartTime(c)
	})
}

func setDbMiddleware(c *gin.Context, db *sql.DB) {
	c.Set("db", db)
	c.Next()
}

func setRedisDbMiddleware(c *gin.Context, redisDb *redis.Client) {
	c.Set("redis_db", redisDb)
	c.Next()
}

func setStartTime(c *gin.Context) {
	startTime := time.Now()
	c.Set("startTime", startTime)
	c.Next()
}

func authMiddleware(c *gin.Context) {
	tokens, err := helpers.ExtractTokens(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		c.Abort()
		return
	}

	claims, err := helpers.ParseToken(tokens.AccessToken, helpers.GetEnv("ACCESS_TOKEN_SECRET"))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, "no access_uuid in token")
		c.Abort()
		return
	}

	redisDb := c.MustGet("redis_db").(*redis.Client)

	val, err := redisDb.Get(context.Background(), accessUuid).Result()
	if len(val) == 0 || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	userId, ok := claims["user_id"]

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, "no user_id in token")
		c.Abort()
		return
	}
	c.Set("user_id", userId)

	c.Next()
}
