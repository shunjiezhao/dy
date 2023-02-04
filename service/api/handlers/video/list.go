package video

import (
	"first/pkg/errno"
	"first/service/api/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func (s *Service) List() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 检查参数
		var (
			err   error
			param ListRequest
		)

		err = c.ShouldBind(&param)
		if err != nil {
			goto ParamErr

		}

		log.Println("获取到 参数", param)
		//	2. 获取数据 绑定

		err = c.ShouldBind(&param)
		if err != nil {
			goto ParamErr

		}

		handlers.BuildResponse(errno.Success)
		return
	ParamErr:
		handlers.BuildResponse(errno.ParamErr)
		c.Abort()
	}
}
