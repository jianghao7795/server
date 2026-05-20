package app

import (
	"strconv"

	appService "server/service/app"

	"server/model/app"
	commentReq "server/model/app/request"
	"server/model/common/request"
	"server/model/common/response"
	"server/utils"

	"github.com/gofiber/fiber/v3"
)

var praiseService = appService.PraiseServer

// CreateComment 创建评论
// @Tags Comment
// @Summary 创建评论
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body app.Comment true "评论信息"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /comment/createComment [post]
func (commentApi *CommentApi) CreateComment(c fiber.Ctx) error {
	var commentData app.Comment
	err := c.Bind().Body(&commentData)
	if err != nil {
		return response.FailWithMessage("获取数据失败", 3, err, c)
	}
	if err := commentService.CreateComment(&commentData); err != nil {
		return response.FailWithMessage("创建失败"+err.Error(), 3, err, c)
	}
	return response.OkWithMessage("创建成功", c)
}

// DeleteComment 删除评论
// @Tags Comment
// @Summary 删除评论
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "评论ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /comment/deleteComment/{id} [delete]
func (commentApi *CommentApi) DeleteComment(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取id失败", 3, err, c)
	}
	if err := commentService.DeleteComment(uint(id)); err != nil {
		return response.FailWithDetailed(err.Error(), "删除失败", 3, err, c)
	}
	return response.OkWithMessage("删除成功", c)
}

// DeleteCommentByIds 批量删除评论
// @Tags Comment
// @Summary 批量删除评论
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /comment/deleteCommentByIds [delete]
func (commentApi *CommentApi) DeleteCommentByIds(c fiber.Ctx) error {
	var IDS request.IdsReq
	err := c.Bind().Body(&IDS)
	if err != nil {
		return response.FailWithMessage("获取id组失败", 3, err, c)
	}
	if err := commentService.DeleteCommentByIds(IDS); err != nil {
		return response.FailWithMessage("批量删除失败", 3, err, c)
	}
	return response.OkWithMessage("批量删除成功", c)
}

// UpdateComment 更新评论
// @Tags Comment
// @Summary 更新评论
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "评论ID"
// @Param data body app.Comment true "评论信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /comment/updateComment/{id} [put]
func (commentApi *CommentApi) UpdateComment(c fiber.Ctx) error {
	var comment2 app.Comment
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取id失败", 3, err, c)
	}
	comment2.ID = uint(id)
	err = c.Bind().Body(&comment2)
	if err != nil {
		return response.FailWithMessage("获取数据失败", 3, err, c)
	}
	if err = commentService.UpdateComment(&comment2); err != nil {
		return response.FailWithMessage("更新失败"+err.Error(), 3, err, c)
	}
	return response.OkWithMessage("更新成功", c)
}

// FindComment 用id查询Comment
// @Tags Comment
// @Summary 用id查询Comment
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "评论ID"
// @Success 200 {object} response.Response{msg=string} "查询成功"
// @Router /comment/getComment/{id} [get]
func (commentApi *CommentApi) FindComment(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("未获取到id参数", 3, err, c)
	}
	if comment, err := commentService.GetComment(id); err != nil {
		return response.FailWithMessage("查询失败"+err.Error(), 3, err, c)
	} else {
		return response.OkWithData(comment, c)
	}
}

// GetCommentList 分页获取Comment列表
// @Tags Comment
// @Summary 分页获取Comment列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query commentReq.CommentSearch true "搜索条件"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /comment/getCommentList [get]
func (commentApi *CommentApi) GetCommentList(c fiber.Ctx) error {
	var pageInfo commentReq.CommentSearch
	_ = c.Bind().Query(&pageInfo)
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}

	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 10
	}
	if list, total, err := commentService.GetCommentInfoList(&pageInfo); err != nil {
		return response.FailWithMessage("获取失败"+err.Error(), 3, err, c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetCommentTreeList 树状获取Comment列表
// @Tags Comment
// @Summary 树状获取Comment列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query commentReq.CommentSearch true "搜索条件"
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Router /comment/getCommentTreeList [get]
func (*CommentApi) GetCommentTreeList(c fiber.Ctx) error {
	var pageInfo commentReq.CommentSearch
	_ = c.Bind().Query(&pageInfo)

	if list, total, err := commentService.GetCommentTreeList(&pageInfo); err != nil {
		return response.FailWithMessage("获取失败"+err.Error(), 3, err, c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// LikeComment 点赞评论
// @Tags Comment
// @Summary 点赞评论
// @Security ApiKeyAuth
// @Produce application/json
// @Param id path integer true "评论ID"
// @Success 200 {object} response.Response{msg=string} "点赞成功"
// @Router /comment/{id}/like [post]
func (*CommentApi) LikeComment(c fiber.Ctx) error {
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取评论ID失败", 3, err, c)
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		return response.FailWithMessage401("请先登录", 3, err, c)
	}

	praise, err := praiseService.LikeComment(uint(commentId), claims.BaseClaims.ID)
	if err != nil {
		return response.FailWithDetailed(err.Error(), "点赞失败", 3, err, c)
	}
	return response.OkWithDetailed(praise, "点赞成功", c)
}

// UnlikeComment 取消点赞评论
// @Tags Comment
// @Summary 取消点赞评论
// @Security ApiKeyAuth
// @Produce application/json
// @Param id path integer true "评论ID"
// @Success 200 {object} response.Response{msg=string} "取消点赞成功"
// @Router /comment/{id}/like [delete]
func (*CommentApi) UnlikeComment(c fiber.Ctx) error {
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取评论ID失败", 3, err, c)
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		return response.FailWithMessage401("请先登录", 3, err, c)
	}

	if err := praiseService.UnlikeComment(uint(commentId), claims.BaseClaims.ID); err != nil {
		return response.FailWithDetailed(err.Error(), "取消点赞失败", 3, err, c)
	}
	return response.OkWithMessage("取消点赞成功", c)
}

// CheckCommentLiked 检查用户是否已点赞
// @Tags Comment
// @Summary 检查用户是否已点赞
// @Security ApiKeyAuth
// @Produce application/json
// @Param id path integer true "评论ID"
// @Success 200 {object} response.Response{msg=string} "查询成功"
// @Router /comment/{id}/like [get]
func (*CommentApi) CheckCommentLiked(c fiber.Ctx) error {
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取评论ID失败", 3, err, c)
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		return response.FailWithMessage401("请先登录", 3, err, c)
	}

	liked, praise, err := praiseService.CheckCommentLiked(uint(commentId), claims.BaseClaims.ID)
	if err != nil {
		return response.FailWithDetailed(err.Error(), "查询点赞状态失败", 3, err, c)
	}
	return response.OkWithDetailed(fiber.Map{"liked": liked, "praise": praise}, "查询成功", c)
}
