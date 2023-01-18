package main

import (
	"first/service/api/handlers/user"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// server.Default() creates a Hertz with recovery middleware.
	// If you need a pure hertz, you can use server.New()
	h := server.Default()

	h.GET("/hello", user.Getuser)

	h.Spin()
}
