package helpers

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"

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

func ParseToken(authorization string, tokenSecret string) (jwt.MapClaims, error) {
	tokenFields := strings.Fields(authorization)

	if len(tokenFields) != 2 {
		return nil, errors.New("invalid token")
	}

	tokenString := tokenFields[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(GetEnv(tokenSecret)), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
