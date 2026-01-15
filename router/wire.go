package router

import (
	"server-fiber/router/app"
	"server-fiber/router/example"
	"server-fiber/router/frontend"
	"server-fiber/router/mobile"
	"server-fiber/router/system"

	"github.com/google/wire"
)

// AppRouterSet App 路由集合
var AppRouterSet = wire.NewSet(
	ProvideArticleRouter,
	ProvideCommentRouter,
	ProvideBaseMessageRouter,
	ProvideUserRouter,
	ProvideTaskRouter,
	ProvideTagRouter,
	ProvideLikeRouter,
	ProvideAppGroup,
)

// ExampleRouterSet Example 路由集合
var ExampleRouterSet = wire.NewSet(
	ProvideCustomerRouter,
	ProvideExcelRouter,
	ProvideFileUploadAndDownloadRouter,
	ProvideExampleGroup,
)

// FrontendRouterSet Frontend 路由集合
var FrontendRouterSet = wire.NewSet(
	ProvideFrontendRouter,
	ProvideFrontendGroup,
)

// MobileRouterSet Mobile 路由集合
var MobileRouterSet = wire.NewSet(
	ProvideMobileLoginRouter,
	ProvideMobileUserRouter,
	ProvideMobileGroup,
)

// SystemRouterSet System 路由集合
var SystemRouterSet = wire.NewSet(
	ProvideApiRouter,
	ProvideGithubRouter,
	ProvideAuthorityBtnRouter,
	ProvideAuthorityRouter,
	ProvideAutoCodeHistoryRouter,
	ProvideAutoCodeRouter,
	ProvideBaseRouter,
	ProvideCasbinRouter,
	ProvideDictionaryDetailRouter,
	ProvideDictionaryRouter,
	ProvideInitRouter,
	ProvideJwtRouter,
	ProvideMenuRouter,
	ProvideOperationRecordRouter,
	ProvideProblemRouter,
	ProvideSysRouter,
	ProvideSystemUserRouter,
	ProvideSystemGroup,
)

// RouterSet 所有路由集合
var RouterSet = wire.NewSet(
	AppRouterSet,
	ExampleRouterSet,
	FrontendRouterSet,
	MobileRouterSet,
	SystemRouterSet,
)

// ========== App Routers ==========

func ProvideArticleRouter() *app.ArticleRouter {
	return &app.ArticleRouter{}
}

func ProvideCommentRouter() *app.CommentRouter {
	return &app.CommentRouter{}
}

func ProvideBaseMessageRouter() *app.BaseMessageRouter {
	return &app.BaseMessageRouter{}
}

func ProvideUserRouter() *app.UserRouter {
	return &app.UserRouter{}
}

func ProvideTaskRouter() *app.TaskRouter {
	return &app.TaskRouter{}
}

func ProvideTagRouter() *app.TagRouter {
	return &app.TagRouter{}
}

func ProvideLikeRouter() *app.LikeRouter {
	return &app.LikeRouter{}
}

func ProvideAppGroup(
	articleRouter *app.ArticleRouter,
	commentRouter *app.CommentRouter,
	baseMessageRouter *app.BaseMessageRouter,
	userRouter *app.UserRouter,
	taskRouter *app.TaskRouter,
	tagRouter *app.TagRouter,
	likeRouter *app.LikeRouter,
) *AppRouter {
	// 创建并返回 AppRouter
	return &AppRouter{
		ArticleRouter:     *articleRouter,
		CommentRouter:     *commentRouter,
		BaseMessageRouter: *baseMessageRouter,
		UserRouter:        *userRouter,
		TaskRouter:        *taskRouter,
		TagRouter:         *tagRouter,
		LikeRouter:        *likeRouter,
	}
}

// ========== Example Routers ==========

func ProvideCustomerRouter() *example.CustomerRouter {
	return &example.CustomerRouter{}
}

func ProvideExcelRouter() *example.ExcelRouter {
	return &example.ExcelRouter{}
}

func ProvideFileUploadAndDownloadRouter() *example.FileUploadAndDownloadRouter {
	return &example.FileUploadAndDownloadRouter{}
}

func ProvideExampleGroup(
	customerRouter *example.CustomerRouter,
	excelRouter *example.ExcelRouter,
	fileUploadRouter *example.FileUploadAndDownloadRouter,
) *ExampleRouter {
	// 创建并返回 ExampleRouter
	return &ExampleRouter{
		CustomerRouter:              *customerRouter,
		ExcelRouter:                 *excelRouter,
		FileUploadAndDownloadRouter: *fileUploadRouter,
	}
}

// ========== Frontend Routers ==========

func ProvideFrontendRouter() *frontend.FrontendRouter {
	return &frontend.FrontendRouter{}
}

func ProvideFrontendGroup(frontendRouter *frontend.FrontendRouter) *FrontendRouter {
	// 创建并返回 FrontendRouter
	return &FrontendRouter{
		FrontendRouter: *frontendRouter,
	}
}

// ========== Mobile Routers ==========

func ProvideMobileLoginRouter() *mobile.MobileLoginRouter {
	return &mobile.MobileLoginRouter{}
}

func ProvideMobileUserRouter() *mobile.MobileUserRouter {
	return &mobile.MobileUserRouter{}
}

func ProvideMobileGroup(
	loginRouter *mobile.MobileLoginRouter,
	userRouter *mobile.MobileUserRouter,
) *MobileRouter {
	// 创建并返回 MobileRouter
	return &MobileRouter{
		MobileLoginRouter: *loginRouter,
		MobileUserRouter:  *userRouter,
	}
}

// ========== System Routers ==========

func ProvideApiRouter() *system.ApiRouter {
	return &system.ApiRouter{}
}

func ProvideGithubRouter() *system.GithubRouter {
	return &system.GithubRouter{}
}

func ProvideAuthorityBtnRouter() *system.AuthorityBtnRouter {
	return &system.AuthorityBtnRouter{}
}

func ProvideAuthorityRouter() *system.AuthorityRouter {
	return &system.AuthorityRouter{}
}

func ProvideAutoCodeHistoryRouter() *system.AutoCodeHistoryRouter {
	return &system.AutoCodeHistoryRouter{}
}

func ProvideAutoCodeRouter() *system.AutoCodeRouter {
	return &system.AutoCodeRouter{}
}

func ProvideBaseRouter() *system.BaseRouter {
	return &system.BaseRouter{}
}

func ProvideCasbinRouter() *system.CasbinRouter {
	return &system.CasbinRouter{}
}

func ProvideDictionaryDetailRouter() *system.DictionaryDetailRouter {
	return &system.DictionaryDetailRouter{}
}

func ProvideDictionaryRouter() *system.DictionaryRouter {
	return &system.DictionaryRouter{}
}

func ProvideInitRouter() *system.InitRouter {
	return &system.InitRouter{}
}

func ProvideJwtRouter() *system.JwtRouter {
	return &system.JwtRouter{}
}

func ProvideMenuRouter() *system.MenuRouter {
	return &system.MenuRouter{}
}

func ProvideOperationRecordRouter() *system.OperationRecordRouter {
	return &system.OperationRecordRouter{}
}

func ProvideProblemRouter() *system.ProblemRouter {
	return &system.ProblemRouter{}
}

func ProvideSysRouter() *system.SysRouter {
	return &system.SysRouter{}
}

func ProvideSystemUserRouter() *system.UserRouter {
	return &system.UserRouter{}
}

func ProvideSystemGroup(
	apiRouter *system.ApiRouter,
	githubRouter *system.GithubRouter,
	authorityBtnRouter *system.AuthorityBtnRouter,
	authorityRouter *system.AuthorityRouter,
	autoCodeHistoryRouter *system.AutoCodeHistoryRouter,
	autoCodeRouter *system.AutoCodeRouter,
	baseRouter *system.BaseRouter,
	casbinRouter *system.CasbinRouter,
	dictionaryDetailRouter *system.DictionaryDetailRouter,
	dictionaryRouter *system.DictionaryRouter,
	initRouter *system.InitRouter,
	jwtRouter *system.JwtRouter,
	menuRouter *system.MenuRouter,
	operationRecordRouter *system.OperationRecordRouter,
	problemRouter *system.ProblemRouter,
	sysRouter *system.SysRouter,
	userRouter *system.UserRouter,
) *SystemRouter {
	// 创建并返回 SystemRouter
	return &SystemRouter{
		ApiRouter:              *apiRouter,
		GithubRouter:           *githubRouter,
		AuthorityBtnRouter:     *authorityBtnRouter,
		AuthorityRouter:        *authorityRouter,
		AutoCodeHistoryRouter:  *autoCodeHistoryRouter,
		AutoCodeRouter:         *autoCodeRouter,
		BaseRouter:             *baseRouter,
		CasbinRouter:           *casbinRouter,
		DictionaryDetailRouter: *dictionaryDetailRouter,
		DictionaryRouter:       *dictionaryRouter,
		InitRouter:             *initRouter,
		JwtRouter:              *jwtRouter,
		MenuRouter:             *menuRouter,
		OperationRecordRouter:  *operationRecordRouter,
		ProblemRouter:          *problemRouter,
		SysRouter:              *sysRouter,
		UserRouter:             *userRouter,
	}
}
