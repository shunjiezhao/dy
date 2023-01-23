package main

import (
	"first/service/api/handlers/user"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	// server.Default() creates a Hertz with recovery middleware.
	// If you need a pure hertz, you can use server.New()

	r := server.Default()
	dy := r.Group("/douyin")

	// 用户相关
	userGroup := dy.Group("user")
	{
		userGroup.POST("register", user.Register)
		userGroup.GET("login", user.Login)
	}

	r.Spin()
}
