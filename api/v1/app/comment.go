package app

import (
	global "server/model"
	"server/model/app"
	commentReq "server/model/app/request"
	"server/model/common/request"
	"server/model/common/response"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// CreateComment 创建评论
// @Tags Comment
// @Summary 创建评论
// @Description 创建新的评论
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body app.Comment true "评论信息"
// @Success 200 {object} response.Response{msg=string} "创建评论成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/createComment [post]
func (commentApi *CommentApi) CreateComment(c *fiber.Ctx) error {
	var commentData app.Comment
	err := c.BodyParser(&commentData)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err := commentService.CreateComment(&commentData); err != nil {
		global.LOG.Error("创建失败!", zap.Error(err))
		return response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		return response.OkWithMessage("创建成功", c)
	}
}

// DeleteComment 删除评论
// @Tags Comment
// @Summary 删除评论
// @Description 根据评论ID删除指定评论
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "评论ID" minimum(1)
// @Success 200 {object} response.Response{msg=string} "删除评论成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/deleteComment/{id} [delete]
func (commentApi *CommentApi) DeleteComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败!", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	if err := commentService.DeleteComment(uint(id)); err != nil {
		global.LOG.Error("删除失败!", zap.Error(err))
		return response.FailWithDetailed(err.Error(), "删除失败", c)
	} else {
		return response.OkWithMessage("删除成功", c)
	}
}

// DeleteCommentByIds 批量删除评论
// @Tags Comment
// @Summary 批量删除评论
// @Description 根据ID列表批量删除评论
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "评论ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除评论成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/deleteCommentByIds [delete]
func (commentApi *CommentApi) DeleteCommentByIds(c *fiber.Ctx) error {
	var IDS request.IdsReq
	err := c.BodyParser(&IDS)
	if err != nil {
		global.LOG.Error("获取id组失败", zap.Error(err))
		return response.FailWithMessage("获取id组失败", c)
	}
	if err := commentService.DeleteCommentByIds(IDS); err != nil {
		global.LOG.Error("批量删除失败!", zap.Error(err))
		return response.FailWithMessage("批量删除失败", c)
	} else {
		return response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateComment 更新评论
// @Tags Comment
// @Summary 更新评论
// @Description 根据评论ID更新评论信息
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param id path integer true "评论ID" minimum(1)
// @Param data body app.Comment true "评论信息"
// @Success 200 {object} response.Response{msg=string} "更新评论成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/updateComment/{id} [put]
func (commentApi *CommentApi) UpdateComment(c *fiber.Ctx) error {
	var comment2 app.Comment
	id, err := c.ParamsInt("id")
	if err != nil {
		global.LOG.Error("获取id失败!", zap.Error(err))
		return response.FailWithMessage("获取id失败", c)
	}
	comment2.ID = uint(id)
	err = c.BodyParser(&comment2)
	if err != nil {
		global.LOG.Error("获取数据失败!", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}
	if err = commentService.UpdateComment(&comment2); err != nil {
		global.LOG.Error("更新失败!", zap.Error(err))
		return response.FailWithMessage("更新失败"+err.Error(), c)
	} else {
		return response.OkWithMessage("更新成功", c)
	}
}

// FindComment 用id查询Comment
// @Tags Comment
// @Summary 用id查询Comment
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id path number true "用id查询Comment"
// @Success 200 {object} response.Response{msg=string} "查询成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/getComment/:id [get]
func (commentApi *CommentApi) FindComment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		global.LOG.Error("未获取到id参数!", zap.Error(err))
		return response.FailWithMessage("未获取到id参数", c)
	}
	if comment, err := commentService.GetComment(id); err != nil {
		global.LOG.Error("查询失败!", zap.Error(err))
		return response.FailWithMessage("查询失败"+err.Error(), c)
	} else {
		return response.OkWithData(comment, c)
	}
}

// GetCommentList 分页获取Comment列表
// @Tags Comment
// @Summary 分页获取Comment列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query commentReq.CommentSearch true "分页获取Comment列表"
// @Success 200 {object} response.Response{msg=string,data=response.PageResult{list=[]app.Comment,total=number,page=number,pageSize=number},code=number} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/getCommentList [get]
func (commentApi *CommentApi) GetCommentList(c *fiber.Ctx) error {
	var pageInfo commentReq.CommentSearch
	_ = c.QueryParser(&pageInfo)
	if pageInfo.Page == 0 {
		pageInfo.Page = 1
	}

	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 10
	}
	if list, total, err := commentService.GetCommentInfoList(&pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetCommentList 树状获取Comment列表
// @Tags Comment
// @Summary 树状获取Comment列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query commentReq.CommentSearch true "分页获取Comment列表"
// @Success 200 {object} response.Response{msg=string,data=response.PageResult{list=[]app.Comment,total=number,page=number,pageSize=number},code=number} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/getCommentTreeList [get]
func (*CommentApi) GetCommentTreeList(c *fiber.Ctx) error {
	var pageInfo commentReq.CommentSearch
	_ = c.QueryParser(&pageInfo)

	if list, total, err := commentService.GetCommentTreeList(&pageInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// PutLikeItOrDislike 点赞
// @Tags Comment
// @Summary 点赞
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body app.Praise true "点赞"
// @Success 200 {object} response.Response{msg=string} "点赞成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /comment/pariseComment [put]
func (*CommentApi) PutLikeItOrDislike(c *fiber.Ctx) error {
	var likeIt app.Praise
	err := c.BodyParser(&likeIt)
	if err != nil {
		global.LOG.Error("获取数据失败", zap.Error(err))
		return response.FailWithMessage("获取数据失败", c)
	}

	if err := commentService.PutLikeItOrDislike(&likeIt); err != nil {
		global.LOG.Error("点赞失败!", zap.Error(err))
		return response.FailWithDetailed(err, "点赞失败", c)
	} else {
		return response.OkWithDetailed(likeIt, "点赞成功", c)
	}
}
