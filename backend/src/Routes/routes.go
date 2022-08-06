package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NoahHakansson/webCLI/backend/src/DB"
	"github.com/NoahHakansson/webCLI/backend/src/JWTAuth"
)

func SetupRoutes() {
	db.SetupDatabase()
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// routes
	r.GET("/api/hello", hello)
	r.POST("/api/login", login)
	r.GET("/api/list-commands", listCommands)
	// auth protected routes
	auth := r.Group("/", authMiddleware)
	auth.GET("/api/auth-hello", hello)
	auth.POST("/api/create-user", createUser)
	auth.POST("/api/create-command", createCommand)

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
		fmt.Println(err.Error())
		cookie = "NotSet"
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("Cookie:", cookie)

	id, err := jwtauth.ValidateJWT(cookie)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
		return
	}

	fmt.Println("User logged in")
	fmt.Println("ID:", id)
}

func createCommand(c *gin.Context) {
	var command db.Command

	if err := c.ShouldBind(&command); err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Command: %#v\n", command)
	if err := db.CreateCommand(&command); err != nil {
		fmt.Println(err.Error())
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, command)
}

func listCommands(c *gin.Context) {
	commands, err := db.ListCommands()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, commands)
}

func createUser(c *gin.Context) {
	var user db.User

	// bind form data and return error if it fails
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{"error": "Missing user credentials"})
		return
	}

	// Create user
	fmt.Printf("User: %#v\n", user)
	if err := db.CreateUser(&user); err != nil {
		fmt.Println(err.Error())
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"username": user.Username, "creation-date": user.CreatedAt})
}

func login(c *gin.Context) {
	var user db.User

	// bind form data and return error if it fails
	if err := c.ShouldBind(&user); err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{"error": "Missing user credentials"})
		return
	}

	// Try to authenticate user
	fmt.Printf("User: %#v\n", user)
	userId, err := db.AuthUser(user.Username, user.Password)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(401, gin.H{"error": "Failed to authenticate user, username or password is wrong."})
		return
	}

	// generate a JWT for the user with user ID
	token, err := jwtauth.GenerateJWT(userId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// create a cookie that's valid for 2 hours to match the JWT 2 hour expiration time
	c.SetCookie("web_cli", token, 60*60*2, "/", "localhost", true, true)

	fmt.Println("Token:", token)

	c.JSON(200, gin.H{"message": "Success"})
}
