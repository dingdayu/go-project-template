package router

import (
	"os"
	"time"

	"github.com/dingdayu/go-project-template/api/controller"
	"github.com/dingdayu/go-project-template/api/middleware"

	"github.com/dingdayu/go-project-template/pkg/logger"
	"github.com/dingdayu/go-project-template/pkg/otel"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var handle *gin.Engine

func Handler() *gin.Engine {
	handle = gin.New()
	handle.ForwardedByClientIP = true

	// 开启 Recover
	handle.Use(middleware.RecoveryWithZap(logger.WithNamespace("recovery"), true))
	// 开启 gzip
	handle.Use(gzip.Gzip(gzip.DefaultCompression))

	// handle.GET("/", func(c *gin.Context) {
	// 	file, _ := assets.IndexFS.ReadFile("index.html")
	// 	c.Data(http.StatusOK, "text/html; charset=utf-8", file)
	// })
	// assets, _ := fs.Sub(dist.AssetsFS, "assets")
	// handle.StaticFS("/assets", http.FS(assets))

	// 探活与采集接口
	handle.HEAD("/health", controller.Hello)
	handle.GET("/health", controller.Hello)
	handle.GET("/ping", controller.Ping)

	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		// otel 中间件
		handle.Use(otel.GinMiddleware())
		handle.GET("/metrics", otel.Prometheus)

		// 注册 http 指标
		handle.Use(otel.HttpRequestMetrics())
	}

	// 服务路由
	handle.GET("/version", controller.Version)

	// 根据配置决定是否启用 api 请求日志
	if viper.GetBool("log.request_log") {
		// 启用请求日志中间件
		handle.Use(middleware.WriterLog(logger.WithNamespace("http")))
	}

	if gin.Mode() == gin.DebugMode {
		handle.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://10.0.7.105:5666", "http://localhost:5666"},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
			AllowHeaders:     []string{"Origin", "Authorization", "X-Token", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	// oauth := handle.Group("/oauth2")
	// {
	// 	oauth.GET("/goto/:platform", oauth2.Goto)
	// 	oauth.GET("/callback", oauth2.Callback)
	// }

	// openapi := handle.Group("/openapi/v1", middleware.OpenAPIAuth())
	// {
	// 	openapi.POST("/token/verify", auth.VerifyToken)
	// 	openapi.GET("/token", auth.ApiToken)
	// 	openapi.POST("/project/:project_no/status-report", project.UploadStatusReport)
	// }

	api := handle.Group("/api")
	// api.POST("/auth/login", auth.Login)

	api.Use(middleware.Authorization())

	// CopilotKit 转发
	// api.POST("copilotkit", copilotkit.Forwarder)

	// 后端渲染菜单
	// api.GET("menu/all", auth.Menu)

	// adminGroup := api.Group("/admin")
	// {
	// 	adminGroup.GET("/navigations", admin.GetNavigationList)
	// 	adminGroup.POST("/navigation/:id/visit", admin.VisitNavigation)
	// }

	// handle.NoRoute(func(c *gin.Context) {
	// 	file, _ := assets.IndexFS.ReadFile("index.html")
	// 	c.Data(http.StatusOK, "text/html; charset=utf-8", file)
	// })
	return handle
}
