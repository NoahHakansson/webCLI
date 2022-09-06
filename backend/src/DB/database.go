// Package db provides db
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
const passMaxLength int = 50
const userMaxLength int = 20

// variables
var db *gorm.DB
var err error

// Functions

// SetupDatabase function
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
	command := &Command{
		Keyword:     "test",
		Description: "some description what this command does",
		Text:        "the actual output when running the command",
		Link:        "https://github.com/NoahHakansson",
	}
	err := CreateCommand(command)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Command: %#v\n", command)

	command = &Command{
		Keyword:     "hello",
		Description: "some description what this command does",
		Text:        "the actual output when running the command",
		Link:        "",
	}
	err = CreateCommand(command)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Command: %#v\n", command)

	// create admin account
	err = CreateUser(&User{Username: "admin", Password: "pass"})
	if err != nil {
		fmt.Println(err.Error())
	}

	commands, err := ListCommands()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("\nCommands: %#v\n", commands)
}

func checkInputLength(username string, password string) (err error) {
	if len(username) > userMaxLength || len(password) > passMaxLength {
		return errors.New("Error: username or password exceeds max length")
	}
	return nil
}

// AuthUser function
func AuthUser(username string, password string) (userID string, err error) {
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
	if err != nil {
		return "", err // return if password doesn't match
	}

	// return user ID
	userID = strconv.Itoa(int(user.ID))
	return userID, nil
}

// CreateCommand function
// Leave link as an empty string if no link is needed for the command.
func CreateCommand(command *Command) error {
	// Disallow reserved command keyword "help"
	if command.Keyword == "help" {
		return errors.New(`Restricted keyword "help" is not allowed`)
	}

	// set creation date
	command.CreatedAt = time.Now()

	// create command in database
	result := db.Create(&command)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ListCommands function
// returns array with all commands from database
func ListCommands() (commands []Command, err error) {
	result := db.Find(&commands)
	if result.Error != nil {
		return nil, result.Error
	}

	return commands, nil
}

// CreateUser function
func CreateUser(user *User) error {
	// check username and password length
	err = checkInputLength(user.Username, user.Password)
	if err != nil {
		return err
	}

	// set creation date
	user.CreatedAt = time.Now()

	// hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPass)

	// create user in database
	result := db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
