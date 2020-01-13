package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func healthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}


func MultipartUrlencodedForm(c *gin.Context)  {
	message := c.PostForm("message")
	nick := c.DefaultPostForm("nick", "anonymous")

	c.JSON(200, gin.H{
		"status": "posted",
		"message": message,
		"nick": nick,
	})
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	/*
	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.HEAD("/someHead", head)
	router.PATCH("/somePatch", patching)
	router.OPTIONS("/someOptions", options)
	*/

	
	router.GET("/hs", healthCheck)
	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", func (c *gin.Context)  {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context)  {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})
	// For each matched request Context will hold the route definition
	// router.POST("/user/:name/*action", func(c *gin.Context) {
	// 	c.FullPath() == "/user/:name/*action"
	// })
	

	// Querystring parameters
	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")				// shortcut for c.Request.URL.Query().Get("lastname")
		c.String(http.StatusOK, "hello %s %s", firstname, lastname)
	})


	router.POST("/form_post", MultipartUrlencodedForm)

	// Grouping routes

	// Simple group:v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}
	// Simple group:v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}


	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run() // listen and serve on 0.0.0.0:8080
	// router.Run(":3000") for a hard coded port
}
