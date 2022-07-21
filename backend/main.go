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

  id, err := jwtauth.ValidateJWT("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Token:", token)
  fmt.Println("ID:", id)
}
