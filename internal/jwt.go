package tracker

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = os.Getenv("JWT_KEY")

func JwtGenerator(userDetails LoginData) (string, error) {
	// Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userDetails.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
