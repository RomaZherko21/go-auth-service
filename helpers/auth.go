package helpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"strconv"
	"time"

	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func CheckPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

type AccessTokenDetails struct {
	AccessToken string `json:"access_token"`
	AtExpires   int64  `json:"at_expires"`
}

func CreateAccessToken(userId int) (*AccessTokenDetails, error) {
	var err error

	td := &AccessTokenDetails{}

	//Creating Access Token
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

type RefreshTokenDetails struct {
	RefreshToken string `json:"refresh_token"`
	RefreshUuid  string `json:"refresh_uuid"`
	RtExpires    int64  `json:"rt_expires"`
}

func CreateRefreshToken(headers http.Header, userId int) (*RefreshTokenDetails, error) {
	var err error

	td := &RefreshTokenDetails{}

	userAgent := headers.Get("User-Agent")

	//Creating Refresh Token
	rtExp, err := strconv.Atoi(GetEnv("REFRESH_TOKEN_EXP_HOUR"))
	if err != nil {
		return nil, err
	}

	td.RtExpires = time.Now().Add(time.Hour * time.Duration(rtExp)).Unix()
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

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func ParseToken(tokenString string, tokenSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("cant parse invalid token")
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

	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
	if err != nil {
		return errors.New("can't extract user_id claim")
	}

	exp, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["exp"]), 10, 64)
	if err != nil {
		return errors.New("can't extract refresh exp claim")
	}

	fmt.Println("EHHE", exp)

	rt := time.Unix(int64(exp), 0)
	now := time.Now()

	err = redis.Set(context.Background(), refreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
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

	c.SetCookie("access_token", token, int((time.Duration(atExp) * time.Minute).Seconds()), "/", "", false, true)

	return nil
}

func SetRefreshTokenCookie(c *gin.Context, token string) error {
	rtExp, err := strconv.Atoi(GetEnv("REFRESH_TOKEN_EXP_HOUR"))
	if err != nil {
		return err
	}

	c.SetCookie("refresh_token", token, int((time.Duration(rtExp) * time.Hour).Seconds()), "/", "", false, true)

	return nil
}
