package main

import (
	"first/pkg/constants"
	"first/service/api/router"
	"first/service/api/rpc"
	"github.com/gin-gonic/gin"
)

func Init() {
	rpc.InitRPC()
}
func main() {
	// server.Default() creates a Hertz with recovery middleware.
	// If you need a pure hertz, you can use server.New()
	Init()
	engine := gin.Default()
	router.InitRouter(engine)

	engine.Run(constants.ApiServerAddress)
}
