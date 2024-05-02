package admin

import (
	"life/dto"
	"life/models"
	"life/service"
	"life/utils"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	AdminBaseController
}

// 获取首页文章
func (c *ArticleController) HomeArticle(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	//获取文章列表
	query := dto.ArticleQuery{}
	query.Limit = limit
	query.Status = 1
	service_article := new(service.ArticleService)
	list := service_article.All(query)

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["list"] = list
	c.SuccessJson(ctx, "success", resp)
}

// 获取列表
func (c *ArticleController) Index(ctx *gin.Context) {
	catid, _ := strconv.Atoi(ctx.Query("catid"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	page_size, _ := strconv.Atoi(ctx.Query("page_size"))
	if page < 1 {
		page = 1
	}
	if page_size < 1 {
		page_size = 10
	}

	//获取文䓬列表
	query := dto.ArticleQuery{}
	query.Catid = catid
	query.Page = page
	query.PageSize = page_size
	query.Status = 1
	service_article := new(service.ArticleService)
	list, total := service_article.PageList(query)

	//组装数据
	resp := make(map[string]interface{}) //创建1个空集合
	resp["total"] = total
	resp["list"] = list
	c.SuccessJson(ctx, "success", resp)
}

// 获取详情
func (c *ArticleController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Query("id"))

	//参数验证
	entity := models.Article{Id: id}
	valid := validation.Validation{}
	valid.Required(entity.Id, "id")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	info, err := service_article.GetById(id)
	if err != nil {
		c.ErrorJson(ctx, -2, err.Error(), nil)
		return
	}
	c.SuccessJson(ctx, "success", info)
}

// 保存
func (c *ArticleController) Save(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	catid, _ := strconv.Atoi(ctx.PostForm("catid"))
	title := ctx.PostForm("title")
	info := ctx.PostForm("info")
	content := ctx.PostForm("content")
	author := ctx.PostForm("author")
	status, _ := strconv.Atoi(ctx.PostForm("status"))

	//参数验证
	entity := models.Article{
		Id:      id,
		Catid:   catid,
		Title:   title,
		Info:    info,
		Content: content,
		Author:  author,
		Status:  int8(status),
	}
	valid := validation.Validation{}
	if entity.Id > 0 {
		valid.Required(entity.Id, "id")
	}
	valid.Required(entity.Title, "title")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.Save(entity)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	//修改时
	if entity.Id > 0 {
		//干掉緩存
		cache_key := "article:id:" + strconv.Itoa(id)
		utils.Redis.Del(cache_key)
	}

	c.SuccessJson(ctx, "success", nil)
}

// 删除
func (c *ArticleController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	//参数验证
	entity := models.Article{Id: id}
	valid := validation.Validation{}
	valid.Required(entity.Id, "id")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	err := ctx.ShouldBind(entity)
	if err != nil {
		c.ErrorJson(ctx, -1, err.Error(), nil)
		return
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.DeleteById(id)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	//干掉緩存
	cache_key := "article:id:" + strconv.Itoa(id)
	utils.Redis.Del(cache_key)

	c.SuccessJson(ctx, "success", nil)
}

// 启用
func (c *ArticleController) Enable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	//参数验证
	entity := models.Article{Id: id}
	valid := validation.Validation{}
	valid.Required(entity.Id, "id")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.EnableById(id)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	//干掉緩存
	cache_key := "article:id:" + strconv.Itoa(id)
	utils.Redis.Del(cache_key)

	c.SuccessJson(ctx, "success", nil)
}

// 禁用
func (c *ArticleController) Disable(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.PostForm("id"))

	//参数验证
	entity := models.Article{Id: id}
	valid := validation.Validation{}
	valid.Required(entity.Id, "id")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.ErrorJson(ctx, -1, err.Key+err.Error(), nil)
			return
		}
	}

	service_article := new(service.ArticleService)
	stat, err := service_article.DisableById(id)
	if stat < 0 {
		c.ErrorJson(ctx, stat, err.Error(), nil)
		return
	}
	//干掉緩存
	cache_key := "article:id:" + strconv.Itoa(id)
	utils.Redis.Del(cache_key)

	c.SuccessJson(ctx, "success", nil)
}