package helpers

import (
	"context"
	"errors"
	"exampleApi/consts"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenDetails struct {
	AccessToken string `json:"access_token"`
	AtExpires   int64  `json:"at_expires"`
}

type RefreshTokenDetails struct {
	RefreshToken string `json:"refresh_token"`
	RefreshUuid  string `json:"refresh_uuid"`
	RtExpires    int64  `json:"rt_expires"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func CheckPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func CreateAccessToken(userId string) (*AccessTokenDetails, error) {
	td := &AccessTokenDetails{}

	atExp, err := strconv.Atoi(GetEnv("ACCESS_TOKEN_EXP_MIN"))
	if err != nil {
		return nil, err
	}

	td.AtExpires = time.Now().Add(time.Minute * time.Duration(atExp)).Unix()

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userId
	atClaims["role"] = "admin"
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(GetEnv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func CreateRefreshToken(headers http.Header, userId string) (*RefreshTokenDetails, error) {
	td := &RefreshTokenDetails{}

	userAgent := headers.Get("User-Agent")

	rtExp, err := strconv.Atoi(GetEnv("REFRESH_TOKEN_EXP_MIN"))
	if err != nil {
		return nil, err
	}

	td.RtExpires = time.Now().Add(time.Duration(rtExp) * time.Minute).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["user_agent"] = userAgent
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(GetEnv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ParseToken(tokenString string, tokenSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("can't parse invalid token")
	}

	return claims, nil
}

func SetRefreshTokenToRedis(redis *redis.Client, refreshToken string) error {
	claims, err := ParseToken(refreshToken, GetEnv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return errors.New("can't parse refresh token")
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)
	if !ok {
		return errors.New("can't extract refresh_uuid claim")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return errors.New("can't extract user_id claim")
	}

	exp, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
	if err != nil {
		return errors.New("can't extract refresh exp claim")
	}

	rt := time.Unix(int64(exp), 0)
	now := time.Now()

	err = redis.Set(context.Background(), refreshUuid, userId, rt.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetAccessTokenCookie(c *gin.Context, token string) error {
	atExp, err := strconv.Atoi(GetEnv("ACCESS_TOKEN_EXP_MIN"))
	if err != nil {
		return err
	}

	c.SetCookie(consts.ACCESS_TOKEN, token, int((time.Duration(atExp) * time.Minute).Seconds()), "/", "", false, true)

	return nil
}

func SetRefreshTokenCookie(c *gin.Context, token string) error {
	rtExp, err := strconv.Atoi(GetEnv("REFRESH_TOKEN_EXP_MIN"))
	if err != nil {
		return err
	}

	c.SetCookie(consts.REFRESH_TOKEN, token, int((time.Duration(rtExp) * time.Minute).Seconds()), "/", "", false, true)

	return nil
}
