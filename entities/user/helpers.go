package user

import (
	"exampleApi/helpers"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func isPasswordCorrect(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func createToken(userEmail string) (string, error) {
	var err error

	atClaims := jwt.MapClaims{}

	atClaims["authorized"] = true
	atClaims["user_email"] = userEmail

	accessTokenExp, err := strconv.Atoi(helpers.GetEnv("ACCESS_TOKEN_EXP"))
	if err != nil {
		return "", err
	}

	atClaims["exp"] = time.Now().Add(time.Minute * time.Duration(accessTokenExp)).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(helpers.GetEnv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}
