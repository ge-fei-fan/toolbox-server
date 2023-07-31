package initialize

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"toolbox-server/global"
	"toolbox-server/middleware"
	"toolbox-server/router"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	systemRouter := router.RouterGroupApp.System
	exampleRouter := router.RouterGroupApp.Example
	Router.Use(middleware.Cors())                                                                     // 直接放行全部跨域请求
	Router.StaticFS(global.TOOL_CONFIG.Local.StorePath, http.Dir(global.TOOL_CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址
	PublicGroup := Router.Group(global.TOOL_CONFIG.System.RouterPrefix)
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	PrivateGroup := Router.Group(global.TOOL_CONFIG.System.RouterPrefix)
	//PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		exampleRouter.InitBilibiliRouter(PrivateGroup)
	}
	{
		systemRouter.InitSystemRouter(PrivateGroup)
	}
	return Router
}
