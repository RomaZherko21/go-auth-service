package helpers

import (
	"context"
	"errors"

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

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUuid   string `json:"access_uuid"`
	RefreshUuid  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

func CreateTokens(userId int) (*TokenDetails, error) {
	var err error

	td := &TokenDetails{}

	//Creating Access Token
	atExp, err := strconv.Atoi(GetEnv("ACCESS_TOKEN_EXP_MIN"))
	if err != nil {
		return nil, err
	}

	td.AtExpires = time.Now().Add(time.Minute * time.Duration(atExp)).Unix()
	td.AccessUuid = uuid.NewV4().String()

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userId
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(GetEnv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtExp, err := strconv.Atoi(GetEnv("REFRESH_TOKEN_EXP_HOUR"))
	if err != nil {
		return nil, err
	}

	td.RtExpires = time.Now().Add(time.Hour * time.Duration(rtExp)).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	rtClaims := jwt.MapClaims{}
	rtClaims["authorized"] = true
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["exp"] = td.AtExpires

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

func ExtractTokens(c *gin.Context) (*Tokens, error) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		return nil, err
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		return nil, err
	}

	return &Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
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

func SetTokensToRedis(redis *redis.Client, userid int, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err := redis.Set(context.Background(), td.AccessUuid, strconv.Itoa(userid), at.Sub(now)).Err()
	if err != nil {
		return err
	}
	err = redis.Set(context.Background(), td.RefreshUuid, strconv.Itoa(userid), rt.Sub(now)).Err()
	if err != nil {
		return err
	}
	return nil
}

func SetTokensToCookie(c *gin.Context, tokenDetails *TokenDetails) error {
	atExp, err := strconv.Atoi(GetEnv("ACCESS_TOKEN_EXP_MIN"))
	if err != nil {
		return err
	}

	rtExp, err := strconv.Atoi(GetEnv("REFRESH_TOKEN_EXP_HOUR"))
	if err != nil {
		return err
	}

	c.SetCookie("access_token", tokenDetails.AccessToken, int((time.Duration(atExp) * time.Minute).Seconds()), "/", "", false, true)
	c.SetCookie("refresh_token", tokenDetails.RefreshToken, int((time.Duration(rtExp) * time.Hour).Seconds()), "/", "", false, true)

	return nil
}
