package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/user", func(context *gin.Context) {
		var user User
		if context.BindJSON(&user) == nil {
			fmt.Println(user.Password)
			fmt.Println(user.UserName)
			context.JSON(200,user);
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type User struct {
	UserName string
	Password string
}
