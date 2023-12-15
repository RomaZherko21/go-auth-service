package auth

import (
	"exampleApi/helpers"
	"exampleApi/helpers/log"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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

	err = helpers.SetTokensToRedis(redisDb, userMeta.ID, tokenDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set tokens to redis"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set tokens to redis: %v", err.Error()))
		return
	}

	c.Set("user_id", userMeta.ID)

	err = helpers.SetTokensToCookie(c, tokenDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set tokens to cookie"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set tokens to cookie: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user authenticated"})
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
	tokens, err := helpers.ExtractTokens(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := helpers.ParseToken(tokens.AccessToken, helpers.GetEnv("ACCESS_TOKEN_SECRET"))

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

func Refresh(c *gin.Context) {
	type RefreshToken struct {
		RefreshToken string `json:"refresh_token"`
	}
	var body RefreshToken

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, err.Error())
		return
	}

	token, err := jwt.Parse(body.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(helpers.GetEnv("REFRESH_TOKEN_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, err.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh expired"})
		log.HttpLog(c, log.Warn, http.StatusUnauthorized, "refresh expired")
		return
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "refresh expired"})
		log.HttpLog(c, log.Warn, http.StatusUnprocessableEntity, "refresh expired")
		return
	}

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Error occurred")
		return
	}

	redisDb := c.MustGet("redis_db").(*redis.Client)
	_, err = redisDb.Del(refreshUuid).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant delete refresh token"})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, "cant delete refresh token")
		return
	}

	tokenDetails, createErr := helpers.CreateTokens(int(userId))
	if createErr != nil {
		c.JSON(http.StatusForbidden, createErr.Error())
		return
	}

	err = helpers.SetTokensToRedis(redisDb, int(userId), tokenDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cant set tokens to redis"})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("Cant set tokens to redis: %v", err.Error()))
		return
	}

	tokens := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}

	c.JSON(http.StatusCreated, tokens)
}
