package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type customClaims struct {
	Username string
	jwt.StandardClaims
}

func CreateToken(username string, issuer string, key string) (string, error) {
	var claim = jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 0, 1).Unix(),
		Issuer:    issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, customClaims{
		username,
		claim,
	})

	return token.SignedString([]byte(key))
}

func GetUsernameFromToken(token string, key string) (string, error) {
	tokenParsed, err := jwt.ParseWithClaims(token, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	return tokenParsed.Claims.(*customClaims).Username, nil
}
