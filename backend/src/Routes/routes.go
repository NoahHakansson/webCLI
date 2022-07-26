package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NoahHakansson/webCLI/backend/src/JWTAuth"
)

var test_user string = "test"
var test_pass string = "test"

type UserCreds struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func SetupRoutes() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// routes
	r.GET("/api/hello", hello)
	r.POST("/api/login", login)
	// auth protected routes
	auth := r.Group("/", authMiddleware)
	auth.GET("/api/auth-hello", hello)

	r.Run(":5000")
}

func hello(c *gin.Context) {
	msg := gin.H{"message": "Hi there!"}
	c.IndentedJSON(http.StatusOK, msg)
}

func authMiddleware(c *gin.Context) {
	// check cookie for valid JWT to see if user is already logged in
	cookie, err := c.Cookie("web_cli")

	if err != nil {
		fmt.Println("user NOT logged in")
		fmt.Println(err)
		cookie = "NotSet"
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("Cookie:", cookie)

	id, err := jwtauth.ValidateJWT(cookie)

	if err != nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("User logged in")
	fmt.Println("ID:", id)
	// c.JSON(200, gin.H{"message": "User already logged in"})
}

func checkUserLogin(user string, pass string) (id string, err error) {
	// TODO: integrate with future database, instead of hardcoded data
	id = "test_ID"
	if user == test_user && pass == test_pass {
		return id, nil
	}
	return "", errors.New("checkUserLogin: User credentials are not valid")
}

func login(c *gin.Context) {
	var user UserCreds

	// bind form data and return error if it fails
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "Missing user credentials"})
		return
	}

	// check login information
	fmt.Printf("User: %#v\n", user)
	userId, err := checkUserLogin(user.Username, user.Password)

	if err != nil {
		fmt.Println(err)
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// generate a JWT for the user with user ID
	token, err := jwtauth.GenerateJWT(userId)

	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// create a cookie that's valid for 2 hours
	c.SetCookie("web_cli", token, 60*60*2, "/", "localhost", true, true)

	fmt.Println("Token:", token)

	c.JSON(200, gin.H{"message": "Success"})
}
