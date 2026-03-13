package frontend

import (
	fileUpload "server/api/v1/app"
	v1 "server/api/v1/frontend"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type FrontendRouter struct{}

func (s *FrontendRouter) InitFrontendRouter(Router fiber.Router) {
	frontend := Router.Group("")
	frontendTagApi := new(v1.TagApi)
	{
		frontend.Get("getTagList", frontendTagApi.GetTagList)
		frontend.Get("getTagArticleList/:id", frontendTagApi.GetTag)
	}
	frontendArticleApi := new(v1.ArticleApi)
	{
		frontend.Get("getArticleList", frontendArticleApi.GetArticleList)
		frontend.Get("getArticle/:id", frontendArticleApi.GetArticleDetail)
		frontend.Get("getSearchArticle/:name/:value", frontendArticleApi.GetSearchArticle)
	}
	frontendCommentApi := new(v1.CommentApi)
	{
		frontend.Get("getArticleComment/:articleId", frontendCommentApi.GetCommentByArticleId)
		frontend.Post("createdComment", middleware.OperationRecord, frontendCommentApi.CreatedComment)
	}
	frontendUserApi := new(v1.User)
	{
		frontend.Get("getImages", middleware.JWTAuth, frontendUserApi.GetImages)
		frontend.Post("login", frontendUserApi.Login)
		frontend.Get("getCurrentUser", middleware.JWTAuth, frontendUserApi.GetCurrent)
		frontend.Put("updateBackgroundImage", middleware.JWTAuth, middleware.OperationRecord, frontendUserApi.UpdateUserBackgroudImage)
		frontend.Put("resetPassword", middleware.JWTAuth, middleware.OperationRecord, frontendUserApi.UpdatePassword)
		frontend.Post("register", middleware.JWTAuth, frontendUserApi.RegisterUser)
		frontend.Put("updateUser", middleware.JWTAuth, middleware.OperationRecord, frontendUserApi.UpdateUser)
	}
	frontendUploadApi := new(fileUpload.FileUploadAndDownloadApi)
	{
		frontend.Post("upload", middleware.OperationRecord, frontendUploadApi.UploadFile)
	}
}
