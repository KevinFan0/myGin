package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main()  {
	// Creates a router without any middleware by default
	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	// Per route middleware, you can add as many as you desire.
	r.GET("benchmark", MyBenchLogger(), benchEndpoint)
	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/", AuthRequired())
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.

}