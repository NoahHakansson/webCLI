package db

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
}

var db *gorm.DB
var err error

func SetupDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=web_cli port=5432 sslmode=disable TimeZone=Europe/Stockholm"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})
	err := CreateUser("admin", "pass")

	if err != nil {
		fmt.Println(err.Error())
	}
}

func checkInputLength(username string, password string) (err error) {
	if len(username) > 20 || len(password) > 50 {
		return errors.New("Error: username or password exceeds max length")
	}
	return nil
}

func AuthUser(username string, password string) (userId string, err error) {
	// check username and password length
	err = checkInputLength(username, password)

	if err != nil {
		return "", err
	}
	// get user from database
	var user User
	result := db.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return "", result.Error
	}

	// compare password and hashedPassword
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// return if password doesn't match
	if err != nil {
		return "", err
	}

	// return user ID
	userId = strconv.Itoa(int(user.ID))
	return userId, nil
}

func CreateUser(username string, password string) (err error) {
	// check username and password length
	err = checkInputLength(username, password)

	if err != nil {
		return err
	}
	// hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// create user
	user := User{
		Username: username,
		Password: string(hashedPass),
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
	}

	result := db.Create(&user)
	// result := db.Select("Username", "Password").Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
