package system

import (
	v1 "server/api/v1/system"
	"server/middleware"

	"github.com/gofiber/fiber/v3"
)

type GithubRouter struct{}

func (g *GithubRouter) InitGithubRouter(Router fiber.Router) {
	githubRouter := Router.Group("github")
	githubRouterApi := new(v1.SystemGithubApi)

	githubRouter.Get("createGithub", middleware.OperationRecord, githubRouterApi.CreateGithub) // 创建github
	githubRouter.Get("getGithubList", githubRouterApi.GetGithubList)

}
