package frontend

import (
	global "server-fiber/model"
	"server-fiber/model/common/response"
	"server-fiber/model/frontend"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
func (s *CommentApi) GetCommentByArticleId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("articleId")
	if err != nil {
		global.LOG.Error("获取articleId 失败!", zap.Error(err))
		return response.FailWithMessage("获取articleId 失败", c)
	}
	if articleComment, err := commentServiceApp.GetCommentByArticleId(id); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "获取失败", c)
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
func (s *CommentApi) CreatedComment(c *fiber.Ctx) error {
	var comment frontend.Comment
	if err := c.BodyParser(&comment); err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage(err.Error(), c)
	}
	if err := commentServiceApp.CreatedComment(&comment); err != nil {
		global.LOG.Error("评论失败!", zap.Error(err))
		return response.FailWithDetailed(fiber.Map{"msg": err.Error()}, "评论失败", c)
	} else {
		return response.OkWithId("评论成功", comment.ID, c)
	}
}
