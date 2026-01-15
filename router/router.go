package router

import (
	"server-fiber/router/app"
	"server-fiber/router/example"
	"server-fiber/router/frontend"
	"server-fiber/router/mobile"
	"server-fiber/router/system"
)

// AppRouter App 路由组
type AppRouter struct {
	app.ArticleRouter
	app.CommentRouter
	app.BaseMessageRouter
	app.UserRouter
	app.TaskRouter
	app.TagRouter
	app.LikeRouter
}

// ExampleRouter Example 路由组
type ExampleRouter struct {
	example.CustomerRouter
	example.ExcelRouter
	example.FileUploadAndDownloadRouter
}

// FrontendRouter Frontend 路由组
type FrontendRouter struct {
	frontend.FrontendRouter
}

// MobileRouter Mobile 路由组
type MobileRouter struct {
	mobile.MobileLoginRouter
	mobile.MobileUserRouter
}

// SystemRouter System 路由组
type SystemRouter struct {
	system.ApiRouter
	system.GithubRouter
	system.AuthorityBtnRouter
	system.AuthorityRouter
	system.AutoCodeHistoryRouter
	system.AutoCodeRouter
	system.BaseRouter
	system.CasbinRouter
	system.DictionaryDetailRouter
	system.DictionaryRouter
	system.InitRouter
	system.JwtRouter
	system.MenuRouter
	system.OperationRecordRouter
	system.ProblemRouter
	system.SysRouter
	system.UserRouter
}

// 为了向后兼容，保留全局变量
var (
	AppRouterInstance      = &AppRouter{}
	ExampleRouterInstance = &ExampleRouter{}
	FrontendRouterInstance = &FrontendRouter{}
	MobileRouterInstance  = &MobileRouter{}
	SystemRouterInstance  = &SystemRouter{}
)
