package app

import (
	"strconv"

	global "server/model"
	modelApp "server/model/app"
	"server/model/common/response"
	appService "server/service/app"

	"server/model/common/request"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// LikeApi 点赞API
type LikeApi struct{}

var likeService = appService.LikeServer

// LikePost 点赞帖子
// @Tags Like
// @Summary 点赞帖子
// @Description 用户点赞指定帖子
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param post_id path integer true "帖子ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "点赞成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /like/likePost/{post_id} [post]
func (l *LikeApi) LikePost(c fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("post_id"))
	if err != nil {
		global.LOG.Error("获取帖子ID失败", zap.Error(err))
		return response.FailWithMessage("获取帖子ID失败", c)
	}

	// 从JWT中获取用户ID
	userID := c.Locals("user_id")
	if userID == nil {
		return response.FailWithMessage("用户未登录", c)
	}
	userId := userID.(uint)

	if err := likeService.LikePost(uint(postId), userId); err != nil {
		global.LOG.Error("点赞失败", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}

	return response.OkWithMessage("点赞成功", c)
}

// UnlikePost 取消点赞帖子
// @Tags Like
// @Summary 取消点赞帖子
// @Description 用户取消点赞指定帖子
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param post_id path integer true "帖子ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "取消点赞成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /like/unlikePost/{post_id} [delete]
func (l *LikeApi) UnlikePost(c fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("post_id"))
	if err != nil {
		global.LOG.Error("获取帖子ID失败", zap.Error(err))
		return response.FailWithMessage("获取帖子ID失败", c)
	}

	// 从JWT中获取用户ID
	userID := c.Locals("user_id")
	if userID == nil {
		return response.FailWithMessage("用户未登录", c)
	}
	userId := userID.(uint)

	if err := likeService.UnlikePost(uint(postId), userId); err != nil {
		global.LOG.Error("取消点赞失败", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}

	return response.OkWithMessage("取消点赞成功", c)
}

// GetPostLikes 获取帖子点赞列表
// @Tags Like
// @Summary 获取帖子点赞列表
// @Description 分页获取指定帖子的点赞用户列表
// @Produce application/json
// @Param post_id path integer true "帖子ID" minimum(1)
// @Param page query integer false "页码" default(1) minimum(1)
// @Param page_size query integer false "每页数量" default(10) minimum(1) maximum(100)
// @Success 200 {object} response.Response{msg=string,data=response.PageResult{list=[]modelApp.Like,total=integer,page=integer,pageSize=integer},code=integer} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /like/getPostLikes/{post_id} [get]
func (l *LikeApi) GetPostLikes(c fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("post_id"))
	if err != nil {
		global.LOG.Error("获取帖子ID失败", zap.Error(err))
		return response.FailWithMessage("获取帖子ID失败", c)
	}

	var page request.Empty

	_ = c.Bind().Query(page)

	if page.Page < 1 {
		page.Page = 1
	}
	if page.PageSize < 1 || page.PageSize > 100 {
		page.PageSize = 10
	}
	var likes []modelApp.Like
	var total int64
	likes, total, err = likeService.GetPostLikes(uint(postId), page.Page, page.PageSize)
	if err != nil {
		global.LOG.Error("获取点赞列表失败", zap.Error(err))
		return response.FailWithMessage("获取点赞列表失败", c)
	}

	return response.OkWithDetailed(response.PageResult{
		List:     likes,
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, "获取成功", c)
}

// CheckUserLiked 检查用户是否点赞了帖子
// @Tags Like
// @Summary 检查用户是否点赞了帖子
// @Description 检查当前用户是否点赞了指定帖子
// @Security ApiKeyAuth
// @Produce application/json
// @Param post_id path integer true "帖子ID" minimum(1)
// @Success 200 {object} response.Response{msg=string,data=object{liked=boolean},code=integer} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /like/checkUserLiked/{post_id} [get]
func (l *LikeApi) CheckUserLiked(c fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("post_id"))
	if err != nil {
		global.LOG.Error("获取帖子ID失败", zap.Error(err))
		return response.FailWithMessage("获取帖子ID失败", c)
	}

	// 从JWT中获取用户ID
	userID := c.Locals("user_id")
	if userID == nil {
		return response.FailWithMessage("用户未登录", c)
	}
	userId := userID.(uint)

	liked, err := likeService.CheckUserLiked(uint(postId), userId)
	if err != nil {
		global.LOG.Error("检查点赞状态失败", zap.Error(err))
		return response.FailWithMessage("检查点赞状态失败", c)
	}

	return response.OkWithDetailed(fiber.Map{"liked": liked}, "获取成功", c)
}

// GetUserLikedPosts 获取用户点赞的帖子列表
// @Tags Like
// @Summary 获取用户点赞的帖子列表
// @Description 分页获取当前用户点赞的帖子列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param page query integer false "页码" default(1) minimum(1)
// @Param page_size query integer false "每页数量" default(10) minimum(1) maximum(100)
// @Success 200 {object} response.Response{msg=string,data=response.PageResult{list=[]modelApp.Like,total=integer,page=integer,pageSize=integer},code=integer} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /like/getUserLikedPosts [get]
func (l *LikeApi) GetUserLikedPosts(c fiber.Ctx) error {
	// 从JWT中获取用户ID
	userID := c.Locals("user_id")
	if userID == nil {
		return response.FailWithMessage("用户未登录", c)
	}
	userId := userID.(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	likes, total, err := likeService.GetUserLikedPosts(userId, page, pageSize)
	if err != nil {
		global.LOG.Error("获取用户点赞列表失败", zap.Error(err))
		return response.FailWithMessage("获取用户点赞列表失败", c)
	}

	return response.OkWithDetailed(response.PageResult{
		List:     likes,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, "获取成功", c)
}

// GetPostLikeCount 获取帖子点赞数
// @Tags Like
// @Summary 获取帖子点赞数
// @Description 获取指定帖子的点赞总数
// @Produce application/json
// @Param post_id path integer true "帖子ID" minimum(1)
// @Success 200 {object} response.Response{msg=string,data=object{like_count=integer},code=integer} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /like/getPostLikeCount/{post_id} [get]
func (l *LikeApi) GetPostLikeCount(c fiber.Ctx) error {
	postId, err := strconv.Atoi(c.Params("post_id"))
	if err != nil {
		global.LOG.Error("获取帖子ID失败", zap.Error(err))
		return response.FailWithMessage("获取帖子ID失败", c)
	}

	count, err := likeService.GetPostLikeCount(uint(postId))
	if err != nil {
		global.LOG.Error("获取点赞数失败", zap.Error(err))
		return response.FailWithMessage("获取点赞数失败", c)
	}

	return response.OkWithDetailed(fiber.Map{"like_count": count}, "获取成功", c)
}
