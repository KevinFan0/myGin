package main

import "github.com/gin-gonic/gin"

func main() {
	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	router := gin.Default()
	//创建不带中间件的路由：
	//r := gin.New()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})


	router.GET("/user/:name", func(c *gin.Context) {
		name :=c.Param("name")
		c.String(http.StatusOK, "Hello % s", name)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
