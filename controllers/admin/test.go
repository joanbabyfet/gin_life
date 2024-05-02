package admin

import (
	"github.com/gin-gonic/gin"
)

type TestController struct {
	AdminBaseController
}

// 测试用
func (c *TestController) Test(ctx *gin.Context) {
	c.SuccessJson(ctx, "success", nil)
}
