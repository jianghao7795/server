package frontend

import (
	global "server/model"
	"server/model/common/response"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// GetImages 获取图片列表
// @Tags Frontend Image
// @Summary 获取图片列表
// @Description 获取前台展示的图片列表
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "获取成功"
// @Failure 500 {object} response.Response "服务器错误"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Router /frontend/getImages [get]
func (u *User) GetImages(c fiber.Ctx) error {
	imageList, err := imagesServiceApp.GetImagesList()
	if err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(imageList, "获取成功", c)
	}
}
