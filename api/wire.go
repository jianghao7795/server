package api

import (
	"server-fiber/api/v1/app"
	"server-fiber/api/v1/example"
	"server-fiber/api/v1/frontend"
	"server-fiber/api/v1/mobile"
	"server-fiber/api/v1/system"
	appService "server-fiber/service/app"
	exampleService "server-fiber/service/example"
	frontendService "server-fiber/service/frontend"
	mobileService "server-fiber/service/mobile"
	systemService "server-fiber/service/system"

	"github.com/google/wire"
)

// AppApiSet App API 集合
var AppApiSet = wire.NewSet(
	ProvideArticleApi,
	ProvideLikeApi,
	ProvideBaseMessageApi,
	ProvideCommentApi,
	ProvideTagApi,
	ProvideFileUploadApi,
	ProvideUserApi,
)

// SystemApiSet System API 集合
var SystemApiSet = wire.NewSet(
	ProvideApiApi,
	ProvideJwtApi,
	ProvideAuthorityBtnApi,
	ProvideAuthorityApi,
	ProvideAutoCodeApi,
	ProvideAutoCodeHistoryApi,
	ProvideBaseMenuApi, // 注意：BaseMenuApi 和 MenuApi 都返回 AuthorityMenuApi，只保留一个
	ProvideCasbinApi,
	ProvideDictionaryDetailApi,
	ProvideDictionaryApi,
	ProvideGithubApi,
	ProvideOperationRecordApi,
	ProvideSystemConfigApi,
	ProvideProblemApi,
	ProvideSystemUserApi,
	ProvideInitDBApi,
)

// FrontendApiSet Frontend API 集合
var FrontendApiSet = wire.NewSet(
	ProvideFrontendArticleApi,
	ProvideFrontendCommentApi,
	ProvideFrontendTagApi,
	// 注意：User 和 Images 可能不需要单独的 API 类型
)

// MobileApiSet Mobile API 集合
var MobileApiSet = wire.NewSet(
	ProvideMobileLoginApi,
	ProvideMobileUserApi,
	ProvideMobileRegisterApi,
)

// ExampleApiSet Example API 集合
var ExampleApiSet = wire.NewSet(
	ProvideFileUploadAndDownloadApi,
	ProvideCustomerApi,
	ProvideExcelApi,
)

// ApiSet 所有 API 集合
var ApiSet = wire.NewSet(
	AppApiSet,
	SystemApiSet,
	FrontendApiSet,
	MobileApiSet,
	ExampleApiSet,
)

// ========== App APIs ==========

func ProvideArticleApi(articleService *appService.ArticleService) *app.ArticleApi {
	return &app.ArticleApi{}
}

func ProvideLikeApi(likeService *appService.LikeService) *app.LikeApi {
	return &app.LikeApi{}
}

func ProvideBaseMessageApi(baseMessageService *appService.BaseMessageService) *app.BaseMessageApi {
	return &app.BaseMessageApi{}
}

func ProvideCommentApi(commentService *appService.CommentService) *app.CommentApi {
	return &app.CommentApi{}
}

func ProvideTagApi(tagService *appService.TagService) *app.TagApi {
	return &app.TagApi{}
}

func ProvideFileUploadApi(fileUploadService *appService.FileUploadService) *app.FileUploadAndDownloadApi {
	return &app.FileUploadAndDownloadApi{}
}

func ProvideUserApi(userService *appService.UserService) *app.UserApi {
	return &app.UserApi{}
}

// ========== System APIs ==========

func ProvideApiApi(apiService *systemService.ApiService) *system.SystemApiApi {
	return &system.SystemApiApi{}
}

func ProvideJwtApi(jwtService *systemService.JwtService) *system.JwtApi {
	return &system.JwtApi{}
}

func ProvideAuthorityBtnApi(authorityBtnService *systemService.AuthorityBtnService) *system.AuthorityBtnApi {
	return &system.AuthorityBtnApi{}
}

func ProvideAuthorityApi(authorityService *systemService.AuthorityService) *system.AuthorityApi {
	return &system.AuthorityApi{}
}

func ProvideAutoCodeApi(autoCodeService *systemService.AutoCodeService) *system.AutoCodeApi {
	return &system.AutoCodeApi{}
}

func ProvideAutoCodeHistoryApi(autoCodeHistoryService *systemService.AutoCodeHistoryService) *system.AutoCodeHistoryApi {
	return &system.AutoCodeHistoryApi{}
}

func ProvideBaseMenuApi(baseMenuService *systemService.BaseMenuService) *system.AuthorityMenuApi {
	return &system.AuthorityMenuApi{}
}

func ProvideCasbinApi(casbinService *systemService.CasbinService) *system.CasbinApi {
	return &system.CasbinApi{}
}

func ProvideDictionaryDetailApi(dictionaryDetailService *systemService.DictionaryDetailService) *system.DictionaryDetailApi {
	return &system.DictionaryDetailApi{}
}

func ProvideDictionaryApi(dictionaryService *systemService.DictionaryService) *system.DictionaryApi {
	return &system.DictionaryApi{}
}

func ProvideGithubApi(githubService *systemService.GithubService) *system.SystemGithubApi {
	return &system.SystemGithubApi{}
}

// ProvideMenuApi 已移除，使用 ProvideBaseMenuApi 代替

func ProvideOperationRecordApi(operationRecordService *systemService.OperationRecordService) *system.OperationRecordApi {
	return &system.OperationRecordApi{}
}

func ProvideSystemConfigApi(systemConfigService *systemService.SystemConfigService) *system.SystemApi {
	return &system.SystemApi{}
}

func ProvideProblemApi(problemService *systemService.Problem) *system.UserProblem {
	return &system.UserProblem{}
}

func ProvideSystemUserApi(userService *systemService.UserService) *system.BaseApi {
	return &system.BaseApi{}
}

func ProvideInitDBApi(initDBService *systemService.InitDBService) *system.DBApi {
	return &system.DBApi{}
}

// ========== Frontend APIs ==========

func ProvideFrontendArticleApi(articleService *frontendService.Article) *frontend.ArticleApi {
	return &frontend.ArticleApi{}
}

func ProvideFrontendCommentApi(commentService *frontendService.Comment) *frontend.CommentApi {
	return &frontend.CommentApi{}
}

func ProvideFrontendTagApi(tagService *frontendService.Tag) *frontend.TagApi {
	return &frontend.TagApi{}
}

// ========== Mobile APIs ==========

func ProvideMobileLoginApi(loginService *mobileService.MobileLoginService) *mobile.LoginApi {
	return &mobile.LoginApi{}
}

func ProvideMobileUserApi(userService *mobileService.MobileUserService) *mobile.UserApi {
	return &mobile.UserApi{}
}

func ProvideMobileRegisterApi(registerService *mobileService.MobileRegisterService) *mobile.RegisterMobile {
	return &mobile.RegisterMobile{}
}

// ========== Example APIs ==========

func ProvideFileUploadAndDownloadApi(fileUploadService *exampleService.FileUploadAndDownloadService) *example.FileUploadAndDownloadApi {
	return &example.FileUploadAndDownloadApi{}
}

func ProvideCustomerApi(customerService *exampleService.CustomerService) *example.CustomerApi {
	return &example.CustomerApi{}
}

func ProvideExcelApi(excelService *exampleService.ExcelService) *example.ExcelApi {
	return &example.ExcelApi{}
}
