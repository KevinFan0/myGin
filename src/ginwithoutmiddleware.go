package main

import (
	"github.com/gin-gonic/gin"
	"time"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Using middleware
func main1()  {
	// Creates a router without any middleware by default
	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	// Per route middleware, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger, benchEndpoint)
	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired)
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)
		// // nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}
	r.Run(":8080")
}

func benchEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "benchEndpoint")
}

func loginEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "loginEndpoint")
}

func submitEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "submitEndpoint")
}


func readEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "readEndpoint")
}


func analyticsEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "analyticsEndpoint")
}

func AuthRequired(c *gin.Context) {
	c.String(http.StatusOK, "AuthRequired")
}

func MyBenchLogger(c *gin.Context) {
	c.String(http.StatusOK, "MyBenchLogger")
}

// How to write log file
func main4()  {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()
	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Use the following code if you need to write the logs to file and console at the same time.
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context)  {
		c.String(200, "pong")
	})
	router.Run(":8080")
}

// Custom Log Format
func main5() {
	router := gin.New()
	
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
						param.ClientIP,
						param.TimeStamp.Format(time.RFC1123),
						param.Method,
						param.Path,
						param.Request.Proto,
						param.StatusCode,
						param.Latency,
						param.Request.UserAgent(),
						param.ErrorMessage,
					)
	}))
	router.Use(gin.Recovery())
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.Run(":8080")

}

// Binding from JSON
type Login struct {
	User		string	`form:"user" json:"user" xml:"user" binding:"required"`
	Password	string	`form: "password" json:"password" xml:"password" binding:"required"`
}

func main6()  {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		if json.User != "manu" || json.Password != "123" {
			c.JSON(http.StatusOK,gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
		
	})
	// Example for binding XML (
	//	<?xml version="1.0" encoding="UTF-8"?>
	//	<root>
	//		<user>user</user>
	//		<password>123</password>
	//	</root>)
	router.POST("/loginXML", func(c *gin.Context) {
		var xml Login
		if err := c.ShouldBindXML(&xml); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if xml.User != "manu" || xml.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})

	// Example for binding a HTML form (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if form.User != "manu" || form.Password != "123" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
	})
	router.Run(":8081")
}