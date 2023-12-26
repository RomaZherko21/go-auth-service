package auth

import (
	"context"
	"exampleApi/helpers"
	"exampleApi/helpers/log"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SignIn(c *gin.Context) {
	var body User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	userMeta, err := UserServiceInstance.GetUserPassword(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	if isPasswordCorrect := helpers.CheckPassword(body.Password, userMeta.Password); !isPasswordCorrect {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong email or password"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("wrong email or password on uid: %v", userMeta.ID))
		return
	}

	accessDetails, err := helpers.CreateAccessToken(userMeta.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create access token"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}

	refreshDetails, err := helpers.CreateRefreshToken(c.Request.Header, userMeta.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create refresh token"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}

	err = helpers.SetAccessTokenCookie(c, accessDetails.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set access token to cookie"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set access token to cookie: %v", err.Error()))
		return
	}

	err = helpers.SetRefreshTokenCookie(c, refreshDetails.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set refresh token to cookie"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set refresh token to cookie: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user authenticated"})
	log.HttpLog(c, log.Warn, http.StatusOK, "user authenticated")
}

func SignUp(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	validationResult := helpers.Validate(&user)

	if !validationResult.OK {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationResult.Errors})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, "validation error")
		return
	}

	password, err := helpers.HashPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, err.Error())
		return
	}

	user.Password = password

	if err := UserServiceInstance.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	log.HttpLog(c, log.Info, http.StatusOK, "User created successfully")
}

func SignOut(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't extract refresh token"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	// add refresh token to black list
	redisDb := c.MustGet("redis_db").(*redis.Client)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	err = helpers.SetRefreshTokenToRedis(redisDb, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set tokens to redis"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set tokens to redis: %v", err.Error()))
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusUnauthorized, gin.H{"message": "User sign out"})
	log.HttpLog(c, log.Info, http.StatusUnauthorized, "User sign out")
}

func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant extract refresh token"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := helpers.ParseToken(refreshToken, helpers.GetEnv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant parse refresh token"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, "Error occurred")
		return
	}

	redisDb := c.MustGet("redis_db").(*redis.Client)

	err = helpers.SetRefreshTokenToRedis(redisDb, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set tokens to redis"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set tokens to redis: %v", err.Error()))
		return
	}

	accessDetails, err := helpers.CreateAccessToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create access token"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}

	refreshDetails, err := helpers.CreateRefreshToken(c.Request.Header, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create refresh token"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}

	err = helpers.SetAccessTokenCookie(c, accessDetails.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set tokens to cookie"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set access to cookie: %v", err.Error()))
		return
	}

	err = helpers.SetRefreshTokenCookie(c, refreshDetails.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set tokens to cookie"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set refresh to cookie: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "refresh successfully"})
	log.HttpLog(c, log.Warn, http.StatusOK, "refresh successfully")
}

func SignOutFromAllDevices(c *gin.Context) {
	redisDb := c.MustGet("redis_db").(*redis.Client)

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't extract user id"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, "can't extract user id")
		return
	}

	ctx := context.Background()
	iter := redisDb.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		currentValue, err := redisDb.Get(ctx, iter.Val()).Result()
		if err != nil {
			log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("Cant remove refresh token of user: %v. Err: %v", userId, err.Error()))
		}

		if currentValue == userId {
			if err := redisDb.Del(context.Background(), iter.Val()).Err(); err != nil {
				log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("Cant remove refresh token: %v, of user: %v. Err: %v", currentValue, userId, err.Error()))
			}
		}
	}
	if err := iter.Err(); err != nil {
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("Cant remove refresh token of user: %v. Err: %v", userId, err.Error()))
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusUnauthorized, gin.H{"message": "User sign out from all devices"})
	log.HttpLog(c, log.Info, http.StatusUnauthorized, "User sign out from all devices")
}
