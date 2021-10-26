package main

import (
	"github.com/arrayJY/lite-im-server/services"
	"github.com/gin-gonic/gin"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

func main() {
	r := gin.Default()

	r.POST("/token", services.CreateToken)
	r.PATCH("/token", services.RefreshToken)

	r.Run()
}
