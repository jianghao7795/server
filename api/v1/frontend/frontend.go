package frontend

import (
	"server/service/frontend"
	systemService "server/service/system"
)

var articleServiceApp = frontend.ArticleServiceApp
var commentServiceApp = frontend.CommentServiceApp
var userServiceApp = frontend.UserServiceApp
var tagServiceApp = frontend.TagServiceApp
var imagesServiceApp = frontend.ImagesServiceApp

var userService = systemService.UserServiceApp
var jwtService = systemService.JwtServiceApp
