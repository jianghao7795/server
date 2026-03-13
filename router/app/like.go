package app

import (
	v1 "server/api/v1/app"
	"server/middleware"

	"github.com/gofiber/fiber/v2"
)

type LikeRouter struct{}

// InitLikeRouter 初始化 like 路由信息
func (s *LikeRouter) InitLikeRouter(Router fiber.Router) {
	likeApi := new(v1.LikeApi)
	likeRouter := Router.Group("like")

	likeRouter.Post("likePost/:post_id", middleware.OperationRecord, likeApi.LikePost).Name("likePost") // 点赞帖子
	likeRouter.Delete("unlikePost/:post_id", middleware.OperationRecord, likeApi.UnlikePost)            // 取消点赞帖子
	likeRouter.Get("getPostLikes/:post_id", likeApi.GetPostLikes)                                       // 获取帖子点赞列表
	likeRouter.Get("checkUserLiked/:post_id", likeApi.CheckUserLiked)                                   // 检查用户是否点赞了帖子
	likeRouter.Get("getUserLikedPosts", likeApi.GetUserLikedPosts)                                      // 获取用户点赞的帖子列表
	likeRouter.Get("getPostLikeCount/:post_id", likeApi.GetPostLikeCount)                               // 获取帖子点赞数
}
