package router

import (
	"first/service/api/handlers/user"
	"first/service/api/handlers/video"
	"github.com/gin-gonic/gin"
)

func InitRouter(engine *gin.Engine) {
	user.InitRouter(engine, user.New(nil, nil)) // 避免 [import cycle]
	video.InitRouter(engine)
}
