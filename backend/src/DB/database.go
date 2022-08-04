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

// constants
const PASS_MAX_LENGTH int = 50
const USER_MAX_LENGTH int = 20

// variables
var db *gorm.DB
var err error

// Functions
func SetupDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=web_cli port=5432 sslmode=disable TimeZone=Europe/Stockholm"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	// migrate database models
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Command{})

	// create some commands
	command, err := createCommand(
		"test",
		"some description what this command does",
		"the actual output when running the command",
		"https://github.com/NoahHakansson")

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Command: %#v\n", command)

	command, err = createCommand(
		"hello",
		"some description what this command does",
		"the actual output when running the command",
		"")

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Command: %#v\n", command)

	// create admin account
	err = CreateUser(&User{Username: "admin", Password: "pass"})

	if err != nil {
		fmt.Println(err.Error())
	}
}

func checkInputLength(username string, password string) (err error) {
	if len(username) > USER_MAX_LENGTH || len(password) > PASS_MAX_LENGTH {
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

// Leave link as an empty string if no link is needed for the command.
func createCommand(keyword string, description string, text string, link string) (cmd Command, err error) {
	// Disallow reserved command keyword "help"
	if keyword == "help" {
		return Command{}, errors.New(`Restricted keyword "help" is not allowed`)
	}

	// create command
	command := Command{
		Keyword:     keyword,
		Description: description,
		Text:        text,
		Link:        link,
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
	}

	result := db.Create(&command)

	if result.Error != nil {
		return Command{}, result.Error
	}

	return command, nil
}

func CreateUser(user *User) (err error) {
	// check username and password length
	err = checkInputLength(user.Username, user.Password)

	if err != nil {
		return err
	}

	// set creation date
	user.CreatedAt = time.Now()

	// hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPass)

	if err != nil {
		return err
	}

	// create user in database
	result := db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
