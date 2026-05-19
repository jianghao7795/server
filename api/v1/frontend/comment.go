package frontend

import (
	"strconv"

	"server/model/common/response"
	"server/model/frontend"
	"server/utils"

	"github.com/gofiber/fiber/v3"
)

type CommentApi struct{}

// GetCommentByArticleId 获取文章评论
// @Tags Frontend Comment
// @Summary 获取文章评论
// @Description 根据文章ID获取相关评论列表
// @Produce application/json
// @Param articleId path integer true "文章ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /frontend/comment/{articleId} [get]
func (s *CommentApi) GetCommentByArticleId(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("articleId"))
	if err != nil {
		return response.FailWithMessage("获取articleId 失败", 3, err, c)
	}
	if articleComment, err := commentServiceApp.GetCommentByArticleId(id); err != nil {
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "获取失败", 3, err, c)
	} else {
		return response.OkWithDetailed(articleComment, "获取成功", c)
	}
}

// CreatedComment 创建评论
// @Tags Frontend Comment
// @Summary 创建评论
// @Description 为指定文章创建新评论
// @Accept application/json
// @Produce application/json
// @Param data body frontend.Comment true "评论信息"
// @Success 200 {object} response.Response{msg=string} "评论成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /frontend/comment [post]
func (s *CommentApi) CreatedComment(c fiber.Ctx) error {
	var comment frontend.Comment
	if err := c.Bind().Body(&comment); err != nil {
		return response.FailWithMessage(err.Error(), 3, err, c)
	}
	if err := commentServiceApp.CreatedComment(&comment); err != nil {
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "评论失败", 3, err, c)
	}
	return response.OkWithId("评论成功", comment.ID, c)
}

// LikeComment 点赞评论
func (s *CommentApi) LikeComment(c fiber.Ctx) error {
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取评论ID失败", 3, err, c)
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		return response.FailWithMessage401("请先登录", 3, err, c)
	}

	praise, err := commentServiceApp.LikeComment(uint(commentId), int64(claims.BaseClaims.ID))
	if err != nil {
		return response.FailWithDetailed(err.Error(), "点赞失败", 3, err, c)
	}
	return response.OkWithDetailed(praise, "点赞成功", c)
}

// UnlikeComment 取消点赞评论
func (s *CommentApi) UnlikeComment(c fiber.Ctx) error {
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取评论ID失败", 3, err, c)
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		return response.FailWithMessage401("请先登录", 3, err, c)
	}

	if err := commentServiceApp.UnlikeComment(uint(commentId), int64(claims.BaseClaims.ID)); err != nil {
		return response.FailWithDetailed(err.Error(), "取消点赞失败", 3, err, c)
	}
	return response.OkWithMessage("取消点赞成功", c)
}

// CheckCommentLiked 检查用户是否已点赞
func (s *CommentApi) CheckCommentLiked(c fiber.Ctx) error {
	commentId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.FailWithMessage("获取评论ID失败", 3, err, c)
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		return response.FailWithMessage401("请先登录", 3, err, c)
	}

	liked, praise, err := commentServiceApp.CheckCommentLiked(uint(commentId), int64(claims.BaseClaims.ID))
	if err != nil {
		return response.FailWithDetailed(err.Error(), "查询点赞状态失败", 3, err, c)
	}
	return response.OkWithDetailed(fiber.Map{"liked": liked, "praise": praise}, "查询成功", c)
}
