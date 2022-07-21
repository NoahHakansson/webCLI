package jwtauth

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type userClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func getEnv() string {
	var secret string
	val, ok := os.LookupEnv("SECRET_SIGN_KEY")
	if !ok {
		// env not set
		secret = "supersecretkey"
	} else {
		// env is set
		secret = val
	}

	return secret
}

var signKey = getEnv()
// var signKey = []byte("supersecretkey")


func GenerateJWT(id string) (string, error) {
	claims := &userClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "webCLI",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(signKey))

	return signedToken, err
}

func ValidateJWT(token string) (string, error) {
	return "", nil
}
