package app

import (
	appService "server/service/app"
)

// ArticleApi article
type ArticleApi struct{}

var articleService = appService.ArticleServer

// BaseMessageApi base message
type BaseMessageApi struct{}

var baseMessageService = appService.BaseMessageServer

// CommentApi comment
type CommentApi struct{}

var commentService = appService.CommentServer

// TagApi tag
type TagApi struct{}

var appTabService = appService.TagServer

// TaskNameApi task
type TaskNameApi struct{}

// FileUploadAndDownloadApi file upload
type FileUploadAndDownloadApi struct{}

var fileUploadService = appService.FileUploadServer

// UserApi
type UserApi struct{}

var userService = appService.UserServer
