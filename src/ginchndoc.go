package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main3()  {
	router := gin.Default()
	http.ListenAndServe(":8080", router)
}