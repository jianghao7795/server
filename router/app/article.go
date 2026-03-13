package app

import (
	v1 "server/api/v1/app"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type ArticleRouter struct{}

// InitArticleRouter 初始化 article 路由信息
func (s *ArticleRouter) InitArticleRouter(Router fiber.Router) {
	articleApi := new(v1.ArticleApi)
	articleRouter := Router.Group("article")

	articleRouter.Post("createArticle", middleware.OperationRecord, articleApi.CreateArticle).Name("createArticle") // 新建article
	articleRouter.Delete("deleteArticle/:id", middleware.OperationRecord, articleApi.DeleteArticle)                 // 删除article
	articleRouter.Delete("deleteArticleByIds", middleware.OperationRecord, articleApi.DeleteArticleByIds)           // 批量删除article
	articleRouter.Put("updateArticle/:id", middleware.OperationRecord, articleApi.UpdateArticle)                    // 更新article
	articleRouter.Put("PutArticleByIds", middleware.OperationRecord, articleApi.PutArticleByIds)                    // 批量更新 是否首页显示article

	articleRouter.Get("findArticle/:id", articleApi.FindArticle)         // 根据ID获取article
	articleRouter.Get("getArticleList", articleApi.GetArticleList)       // 获取article列表
	articleRouter.Get("getArticleReading", articleApi.GetArticleReading) // 获取阅读量
}
