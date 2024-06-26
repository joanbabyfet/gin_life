// 父控制器
package controllers

import (
	"life/models"
	"life/service"

	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type FeedbackController struct {
	BaseController
}

// 保存
func (c *FeedbackController) Save(ctx *gin.Context) {
	name := ctx.PostForm("name")
	mobile := ctx.PostForm("mobile")
	email := ctx.PostForm("email")
	content := ctx.PostForm("content")

	//参数验证
	entity := models.Feedback{
		Name:    name,
		Mobile:  mobile,
		Email:   email,
		Content: content,
	}
	rules := govalidator.MapData{}
	rules["name"] = []string{"required"}
	messages := govalidator.MapData{}
	messages["name"] = []string{"required:name 不能为空"}
	opts := govalidator.Options{
		Data:            &entity,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: false,
	}
	valid := govalidator.New(opts)
	e := valid.ValidateStruct()
	if len(e) > 0 {
		for _, err := range e {
			c.ErrorJson(ctx, -1, err[0], nil)
			return
		}
	}

	service_feedback := new(service.FeedbackService)
	stat, err := service_feedback.Save(entity)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", nil)
}
