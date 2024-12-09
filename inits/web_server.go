package inits

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"local-audio-lib/handlers"
)

func router(e *echo.Echo) {
	// 健康状态检查
	e.GET("/", handlers.Health)

	// 媒体列表
	e.GET("/index", handlers.Index)

	// 封面图片
	e.GET("/cover/:id", handlers.Cover)

	// 音频文件
	e.GET("/audio/:id", handlers.Audio)
}

func WebServer(listen string) error {
	// 创建服务器
	e := echo.New()

	// 使用 CORS 中间件
	e.Use(middleware.CORS())

	// 绑定路由
	router(e)

	// 启动服务器
	return e.Start(listen)
}
