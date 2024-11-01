package util

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type UserClaims struct {
	jwt.StandardClaims

	Id string `json:"id"`
}

var JWT_SECRET string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	JWT_SECRET = os.Getenv("JWT_SECRET")
}

func GetUserIdFromJWT(jwtToken string) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(jwtToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return "", err
	}

	if !parsedToken.Valid {
		// comment this "if" if you want to test with expired token
		return "", errors.New("invalid token")
	}

	userClaims := parsedToken.Claims.(*UserClaims)

	return userClaims.Id, nil
}

func GetUserIdFromRequestObject(r *http.Request) (string, error) {
	jwtCookie, err := r.Cookie("jwt")
	if err != nil {
		return "", err
	}

	jwtToken := jwtCookie.Value

	return GetUserIdFromJWT(jwtToken)
}
