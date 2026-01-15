package service

import (
	"server-fiber/service/app"
	"server-fiber/service/example"
	"server-fiber/service/frontend"
	"server-fiber/service/mobile"
	"server-fiber/service/system"

	"github.com/google/wire"
	"gorm.io/gorm"
)

// AppServiceSet App 服务集合
var AppServiceSet = wire.NewSet(
	ProvideArticleService,
	ProvideLikeService,
	ProvideBaseMessageService,
	ProvideCommentService,
	ProvideTagService,
	ProvideFileUploadService,
	ProvideUserService,
)

// SystemServiceSet System 服务集合
var SystemServiceSet = wire.NewSet(
	ProvideApiService,
	ProvideJwtService,
	ProvideAuthorityBtnService,
	ProvideAuthorityService,
	ProvideAutoCodeService,
	ProvideAutoCodeHistoryService,
	ProvideBaseMenuService,
	ProvideCasbinService,
	ProvideDictionaryDetailService,
	ProvideDictionaryService,
	ProvideGithubService,
	ProvideMenuService,
	ProvideOperationRecordService,
	ProvideSystemConfigService,
	ProvideProblemService,
	ProvideSystemUserService,
	ProvideInitDBService,
)

// FrontendServiceSet Frontend 服务集合
var FrontendServiceSet = wire.NewSet(
	ProvideFrontendArticleService,
	ProvideFrontendCommentService,
	ProvideFrontendUserService,
	ProvideFrontendTagService,
	ProvideFrontendImagesService,
)

// MobileServiceSet Mobile 服务集合
var MobileServiceSet = wire.NewSet(
	ProvideMobileLoginService,
	ProvideMobileUserService,
	ProvideMobileRegisterService,
)

// ExampleServiceSet Example 服务集合
var ExampleServiceSet = wire.NewSet(
	ProvideFileUploadAndDownloadService,
	ProvideCustomerService,
	ProvideExcelService,
)

// ServiceSet 所有服务集合
var ServiceSet = wire.NewSet(
	AppServiceSet,
	SystemServiceSet,
	FrontendServiceSet,
	MobileServiceSet,
	ExampleServiceSet,
)

// ========== App Services ==========

func ProvideArticleService(db *gorm.DB) *app.ArticleService {
	// ArticleService 使用 global.DB，但保留 db 参数以便未来扩展
	return &app.ArticleService{}
}

func ProvideLikeService() *app.LikeService {
	return &app.LikeService{}
}

func ProvideBaseMessageService(db *gorm.DB) *app.BaseMessageService {
	return &app.BaseMessageService{}
}

func ProvideCommentService(db *gorm.DB) *app.CommentService {
	return &app.CommentService{}
}

func ProvideTagService(db *gorm.DB) *app.TagService {
	return &app.TagService{}
}

func ProvideFileUploadService(db *gorm.DB) *app.FileUploadService {
	return &app.FileUploadService{}
}

func ProvideUserService(db *gorm.DB) *app.UserService {
	return &app.UserService{}
}

// ========== System Services ==========

func ProvideApiService() *system.ApiService {
	return &system.ApiService{}
}

func ProvideJwtService() *system.JwtService {
	return &system.JwtService{}
}

func ProvideAuthorityBtnService() *system.AuthorityBtnService {
	return &system.AuthorityBtnService{}
}

func ProvideAuthorityService() *system.AuthorityService {
	return &system.AuthorityService{}
}

func ProvideAutoCodeService() *system.AutoCodeService {
	return &system.AutoCodeService{}
}

func ProvideAutoCodeHistoryService() *system.AutoCodeHistoryService {
	return &system.AutoCodeHistoryService{}
}

func ProvideBaseMenuService() *system.BaseMenuService {
	return &system.BaseMenuService{}
}

func ProvideCasbinService() *system.CasbinService {
	return &system.CasbinService{}
}

func ProvideDictionaryDetailService() *system.DictionaryDetailService {
	return &system.DictionaryDetailService{}
}

func ProvideDictionaryService() *system.DictionaryService {
	return &system.DictionaryService{}
}

func ProvideGithubService() *system.GithubService {
	return &system.GithubService{}
}

func ProvideMenuService() *system.MenuService {
	return &system.MenuService{}
}

func ProvideOperationRecordService() *system.OperationRecordService {
	return &system.OperationRecordService{}
}

func ProvideSystemConfigService() *system.SystemConfigService {
	return &system.SystemConfigService{}
}

func ProvideProblemService() *system.Problem {
	return &system.Problem{}
}

func ProvideSystemUserService() *system.UserService {
	return &system.UserService{}
}

func ProvideInitDBService() *system.InitDBService {
	return &system.InitDBService{}
}

// ========== Frontend Services ==========

func ProvideFrontendArticleService() *frontend.Article {
	return &frontend.Article{}
}

func ProvideFrontendCommentService() *frontend.Comment {
	return &frontend.Comment{}
}

func ProvideFrontendUserService() *frontend.FrontendUser {
	return &frontend.FrontendUser{}
}

func ProvideFrontendTagService() *frontend.Tag {
	return &frontend.Tag{}
}

func ProvideFrontendImagesService() *frontend.Images {
	return &frontend.Images{}
}

// ========== Mobile Services ==========

func ProvideMobileLoginService() *mobile.MobileLoginService {
	return &mobile.MobileLoginService{}
}

func ProvideMobileUserService() *mobile.MobileUserService {
	return &mobile.MobileUserService{}
}

func ProvideMobileRegisterService() *mobile.MobileRegisterService {
	return &mobile.MobileRegisterService{}
}

// ========== Example Services ==========

func ProvideFileUploadAndDownloadService() *example.FileUploadAndDownloadService {
	return &example.FileUploadAndDownloadService{}
}

func ProvideCustomerService() *example.CustomerService {
	return &example.CustomerService{}
}

func ProvideExcelService() *example.ExcelService {
	return &example.ExcelService{}
}

