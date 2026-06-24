package init_load

import (
	"server/middleware"
	"server/router"
	"time"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func Routers(
	app *fiber.App,
	appRouter *router.AppRouter,
	systemRouter *router.SystemRouter,
	exampleRouter *router.ExampleRouter,
	frontendRouter *router.FrontendRouter,
	mobileRouter *router.MobileRouter,
) *fiber.App {
	staticConfig := static.Config{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		CacheDuration: 100 * time.Second,
		MaxAge:        36000,
	}

	app.Get("/api/uploads/*", static.New("uploads", staticConfig))
	app.Get("/backend/uploads/*", static.New("uploads", staticConfig))
	app.Get("/backend/public/*", static.New("public", staticConfig))
	app.Get("/mobile/uploads/*", static.New("uploads", staticConfig))
	app.Get("/backend/logs/*", static.New("logs", staticConfig))
	app.Get("/api/logs/*", static.New("logs", staticConfig))
	app.Get("/docs/*", static.New("./docs"))

	app.Get("/swagger/*", swaggo.New(swaggo.Config{
		URL:    "/docs/swagger.json",
		Title:  "server API documentation",
		Layout: "StandaloneLayout",
	}))

	app.Use("/ws", middleware.Ws)
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		for {
			messageType, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			if err = c.WriteMessage(messageType, msg); err != nil {
				break
			}
		}
	}))

	routers := app.Use(middleware.AppMiddlewares()...)
	routers.Get("/backend/form-generator/*", static.New("resource/page"))

	backendRouterNotLogin := routers.Group("/backend")
	systemRouter.InitBaseRouter(backendRouterNotLogin)
	systemRouter.InitInitRouter(backendRouterNotLogin)

	backendRouter := backendRouterNotLogin.Use(middleware.JWTAuth, middleware.CasbinHandler)
	{
		systemRouter.InitApiRouter(backendRouter)
		systemRouter.InitJwtRouter(backendRouter)
		systemRouter.InitUserRouter(backendRouter)
		systemRouter.InitMenuRouter(backendRouter)
		systemRouter.InitSystemRouter(backendRouter)
		systemRouter.InitCasbinRouter(backendRouter)
		systemRouter.InitAutoCodeRouter(backendRouter)
		systemRouter.InitAuthorityRouter(backendRouter)
		systemRouter.InitSysDictionaryRouter(backendRouter)
		systemRouter.InitAutoCodeHistoryRouter(backendRouter)
		systemRouter.InitSysOperationRecordRouter(backendRouter)
		systemRouter.InitSysDictionaryDetailRouter(backendRouter)
		systemRouter.InitAuthorityBtnRouterRouter(backendRouter)
		systemRouter.InitProblemRouter(backendRouter)
		systemRouter.InitGithubRouter(backendRouter)

		exampleRouter.InitExcelRouter(backendRouter)
		exampleRouter.InitCustomerRouter(backendRouter)
		exampleRouter.InitFileUploadAndDownloadRouter(backendRouter)

		appRouter.InitTagRouter(backendRouter)
		appRouter.InitArticleRouter(backendRouter)
		appRouter.InitCommentRouter(backendRouter)
		appRouter.InitBaseMessageRouter(backendRouter)
		appRouter.InitTaskRouter(backendRouter)
		appRouter.InitUserRouter(backendRouter)
		appRouter.InitLikeRouter(backendRouter)

		mobileRouter.InitMobileRouter(backendRouter)
	}

	publicGroup := routers.Group("api")
	frontendRouter.InitFrontendRouter(publicGroup)
	mobileRouter.InitMobileLoginRouter(routers)
	InstallPlugin(backendRouter)

	return app
}
