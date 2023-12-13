package helpers

import (
	"errors"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

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
