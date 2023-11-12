package main

import (
	"github.com/gin-gonic/gin"
	"inf_test/handler"
)

func main() {
	r := gin.Default()

	r.POST("/signup", handler.Signup)
	r.GET("users/:user_id", handler.GetUser)
	r.PATCH("users/:user_id", handler.UpdateUser)
	r.DELETE("/close", handler.DeleteUser)

	r.Run()
}
