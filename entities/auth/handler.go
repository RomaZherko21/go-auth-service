package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"exampleApi/consts"
	"exampleApi/helpers"
	"exampleApi/helpers/log"
)

func SignIn(c *gin.Context) {
	const USER_ERROR_MESSAGE = "wrong email or password"
	var body User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": USER_ERROR_MESSAGE})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("%v: %v", consts.INVALID_BODY, err.Error()))
		return
	}

	userMeta, err := UserServiceInstance.GetUserPassword(c, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": USER_ERROR_MESSAGE})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("cant't get user from db: %v", err.Error()))
		return
	}

	if isPasswordCorrect := helpers.CheckPassword(body.Password, userMeta.Password); !isPasswordCorrect {
		c.JSON(http.StatusBadRequest, gin.H{"error": USER_ERROR_MESSAGE})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("wrong email or password on uid: %v", userMeta.ID))
		return
	}

	accessDetails, err := helpers.CreateAccessToken(c.Request.Header, userMeta.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("can't create access token: %v", err.Error()))
		return
	}

	refreshDetails, err := helpers.CreateRefreshToken(c.Request.Header, userMeta.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("can't create refresh token: %v", err.Error()))
		return
	}

	err = helpers.SetAccessTokenCookie(c, accessDetails.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("can't set access token to cookie: %v", err.Error()))
		return
	}

	err = helpers.SetRefreshTokenCookie(c, refreshDetails.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Warn, http.StatusInternalServerError, fmt.Sprintf("can't set refresh token to cookie: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user authenticated"})
	log.HttpLog(c, log.Warn, http.StatusOK, fmt.Sprintf("user authenticated: %v", userMeta.ID))
}

func SignUp(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": consts.INVALID_BODY})
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("%v: %v", consts.INVALID_BODY, err.Error()))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("can't hash password: %v", err.Error()))
		return
	}

	user.Password = password

	if err := UserServiceInstance.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("can't create user: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user was created successfully"})
	log.HttpLog(c, log.Info, http.StatusOK, fmt.Sprintf("user was created successfully: %v", user.Email))
}

func SignOut(c *gin.Context) {
	refreshToken, err := c.Cookie(consts.REFRESH_TOKEN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no refresh token"})
		log.HttpLog(c, log.Error, http.StatusBadRequest, fmt.Sprintf("can't extract refresh token: %v", err.Error()))
		return
	}

	redisDb := c.MustGet("redis_db").(*redis.Client)

	err = helpers.SetRefreshTokenToRedis(redisDb, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("сan't set token to redis: %v", err.Error()))
		return
	}

	c.SetCookie(consts.ACCESS_TOKEN, "", -1, "/", "", false, true)
	c.SetCookie(consts.REFRESH_TOKEN, "", -1, "/", "", false, true)

	c.JSON(http.StatusUnauthorized, gin.H{"message": "user sign out"})
	log.HttpLog(c, log.Info, http.StatusUnauthorized, "user sign out")
}

func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(consts.REFRESH_TOKEN)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no refresh token"})
		log.HttpLog(c, log.Error, http.StatusBadRequest, fmt.Sprintf("can't extract refresh token: %v", err.Error()))
		return
	}

	claims, err := helpers.ParseToken(refreshToken, helpers.GetEnv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is invalid"})
		log.HttpLog(c, log.Error, http.StatusBadRequest, fmt.Sprintf("refresh token is invalid: %v", err.Error()))
		return
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, "сan't extract user_id claim")
		return
	}

	redisDb := c.MustGet("redis_db").(*redis.Client)

	err = helpers.SetRefreshTokenToRedis(redisDb, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("can't set token to redis: %v", err.Error()))
		return
	}

	accessDetails, err := helpers.CreateAccessToken(c.Request.Header, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("can't create access token: %v", err.Error()))
		return
	}

	refreshDetails, err := helpers.CreateRefreshToken(c.Request.Header, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("can't create refresh token: %v", err.Error()))
		return
	}

	err = helpers.SetAccessTokenCookie(c, accessDetails.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("can't set access to cookie: %v", err.Error()))
		return
	}

	err = helpers.SetRefreshTokenCookie(c, refreshDetails.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, fmt.Sprintf("can't set refresh to cookie: %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "refresh successfully"})
	log.HttpLog(c, log.Warn, http.StatusOK, "refresh successfully")
}

func SignOutFromAllDevices(c *gin.Context) {
	redisDb := c.MustGet("redis_db").(*redis.Client)

	userId, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": consts.SOMETHING_WENT_WRONG})
		log.HttpLog(c, log.Error, http.StatusInternalServerError, "can't extract user_id claim")
		return
	}

	ctx := context.Background()
	iter := redisDb.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		currentValue, err := redisDb.Get(ctx, iter.Val()).Result()
		if err != nil {
			log.HttpLog(c, log.Warn, http.StatusUnprocessableEntity, fmt.Sprintf("can't remove refresh token of user: %v. Err: %v", userId, err.Error()))
		}

		if currentValue == userId {
			if err := redisDb.Del(context.Background(), iter.Val()).Err(); err != nil {
				log.HttpLog(c, log.Warn, http.StatusUnprocessableEntity, fmt.Sprintf("can't remove refresh token: %v, of user: %v. Err: %v", currentValue, userId, err.Error()))
			}
		}
	}
	if err := iter.Err(); err != nil {
		log.HttpLog(c, log.Warn, http.StatusBadRequest, fmt.Sprintf("can't remove refresh token of user: %v. Err: %v", userId, err.Error()))
	}

	c.SetCookie(consts.ACCESS_TOKEN, "", -1, "/", "", false, true)
	c.SetCookie(consts.REFRESH_TOKEN, "", -1, "/", "", false, true)

	c.JSON(http.StatusUnauthorized, gin.H{"message": "user sign out from all devices"})
	log.HttpLog(c, log.Info, http.StatusUnauthorized, "user sign out from all devices")
}
