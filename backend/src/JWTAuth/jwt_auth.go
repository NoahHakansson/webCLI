// Package jwtauth provides jwtauth
package jwtauth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type userClaims struct {
	ID string `json:"id"`
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

// GenerateJWT function
func GenerateJWT(id string) (string, error) {
	claims := &userClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			Issuer:    "webCLI",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(signKey))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return signedToken, nil
}

// ValidateJWT get user ID from JWT
func ValidateJWT(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if claims, ok := token.Claims.(*userClaims); ok && token.Valid {
		// valid token
		return claims.ID, nil
	}

	return "", err
}
