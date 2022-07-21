package src

import (
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type userClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

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
	signedToken, err := token.SignedString(signKey)

	return signedToken, err
}
