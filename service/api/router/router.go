package router

import (
	"first/service/api/handlers/user"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	user.InitRouter(engine, user.New()) // 避免 [import cycle]
}
