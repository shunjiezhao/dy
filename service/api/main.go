package main

import (
	"context"
	"first/pkg/constants"
	"first/service/api/router"
	"first/service/api/rpc"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	server := &http.Server{
		Addr:           constants.ApiServerAddress,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go server.ListenAndServe()
	//优雅退出
	sig := <-ch
	fmt.Println("got a signal", sig)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	_ = server.Shutdown(ctx)
	fmt.Println("------exited--------")
}
