package user

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

func Getuser(ctx context.Context, c *app.RequestContext) {
	var param UserParam
	err := c.BindAndValidate(&param)
	if err != nil {
		c.String(http.StatusUnauthorized, "检查参数")
		return
	}
	c.String(http.StatusOK, "成功")
}
