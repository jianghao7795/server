package system

import (
	// json "github.com/bytedance/sonic"
	"encoding/json"
	"io"
	"net/http"
	global "server/model"
	"server/model/common/request"
	"server/model/common/response"
	"server/model/system"
	"server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type SystemGithubApi struct{}

// @Tags SysGithub
// @Summary 获取GitHub提交列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Success 200 {object} response.Response{data=response.PageResult{list=[]system.SysGithub},msg=string} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "获取失败"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /github/getGithubList [get]
func (g *SystemGithubApi) GetGithubList(c *fiber.Ctx) error {
	var searchInfo request.PageInfo
	_ = c.QueryParser(&searchInfo)
	if list, err := githubService.GetGithubList(searchInfo); err != nil {
		global.LOG.Error("获取失败!", zap.Error(err))
		return response.FailWithMessage("获取失败", c)
	} else {
		return response.OkWithDetailed(response.PageResult{
			List:     list,
			Page:     searchInfo.Page,
			PageSize: searchInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags SysGithub
// @Summary 创建GitHub提交记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Description 从GitHub API获取最新的提交记录并保存到数据库
// @Success 200 {object} response.Response{data=fiber.Map{total=int},msg=string} "创建成功"
// @Failure 400 {object} response.Response{msg=string} "创建失败或网络错误"
// @Failure 401 {object} response.Response{msg=string} "未授权"
// @Failure 500 {object} response.Response{msg=string} "服务器错误"
// @Router /github/createGithub [get]
func (g *SystemGithubApi) CreateGithub(c *fiber.Ctx) error {
	isCheck := utils.NetWorkStatus("https://api.github.com")
	if !isCheck {
		global.LOG.Error("网络错误：networking not work")
		return response.FailWithMessage("网络错误： networking not work", c)
	}
	data := make([]system.SysGithub, 1)

	page := "1"
	per_page := "5"
	resp, err := http.Get("https://api.github.com/repos/JiangHaoCode/server-web/commits?page=" + page + "&per_page=" + per_page)
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		global.LOG.Error("请求Commit错误", zap.Error(err))
		return response.FailWithMessage("请求Commit错误", c)
	}
	body, _ := io.ReadAll(resp.Body)
	// respData := new([]GithubCommit)
	var respData []system.GithubCommit
	json.Unmarshal(body, &respData)
	time.LoadLocation("Asia/Shanghai")
	for _, val := range respData {
		var temp system.SysGithub
		temp.Author = val.Commit.Author.Name
		temp.CommitTime = val.Commit.Author.Date.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
		temp.Message = val.Commit.Message
		data = append(data, temp)
	}
	if total, err := githubService.CreateGithub(&data); err != nil {
		global.LOG.Error("创建commit有错误!", zap.Error(err))
		return response.FailWithMessage("创建commit有错误!", c)
	} else {
		return response.OkWithData(map[string]int{"total": total}, c)
	}
}
