package auth

import (
	"exampleApi/helpers"
	"exampleApi/helpers/log"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

func SignIn(c *gin.Context) {
	redisDb := c.MustGet("redis_db").(*redis.Client)

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

	tokenDetails, err := helpers.CreateTokens(userMeta.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}

	at := time.Unix(tokenDetails.AtExpires, 0)
	rt := time.Unix(tokenDetails.RtExpires, 0)
	now := time.Now()

	err = redisDb.Set(tokenDetails.AccessUuid, strconv.Itoa(int(userMeta.ID)), at.Sub(now)).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}
	err = redisDb.Set(tokenDetails.RefreshUuid, strconv.Itoa(int(userMeta.ID)), rt.Sub(now)).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, err.Error())
		return
	}

	c.Set("user_id", userMeta.ID)

	c.JSON(http.StatusOK, gin.H{"access_token": tokenDetails})
	log.HttpLog(c, log.Warn, http.StatusBadRequest, "user authenticated")
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
	claims, err := helpers.ParseToken(c.GetHeader("authorization"), helpers.GetEnv("ACCESS_TOKEN_SECRET"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	accessUuid, ok := claims["access_uuid"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant get access_uuid claim"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, "cant get access_uuid claim")
		return
	}

	redisDb := c.MustGet("redis_db").(*redis.Client)

	_, err = redisDb.Del(accessUuid).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant delete access token"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, "cant delete access token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User sign out"})
	log.HttpLog(c, log.Info, http.StatusOK, "User sign out")
}
