package main

import (
	"fmt"
	"log"

	"github.com/NoahHakansson/webCLI/backend/src/JWTAuth"
)

func main() {
  token, err := jwtauth.GenerateJWT("test")

  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Token:", token)
}
