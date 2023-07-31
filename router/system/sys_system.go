package system

import (
	"github.com/gin-gonic/gin"
	v1 "toolbox-server/api/v1"
)

type SystemRouter struct{}

func (s *SystemRouter) InitSystemRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	systemRouter := Router.Group("system")
	//服务相关接口
	systemApi := v1.ApiGroupApp.SystemApiGroup.SystemApi
	{
		systemRouter.GET("shutdown", systemApi.Shutdown)
		systemRouter.GET("version", systemApi.Version)
		systemRouter.GET("mode", systemApi.Mode)
		systemRouter.GET("test", systemApi.Test)
	}
	//配置相关接口
	configRouter := systemRouter.Group("config")
	{
		configRouter.GET("list", systemApi.GetConfig)
		configRouter.POST("update", systemApi.UpdateConfig)
	}
	//定时任务相关接口
	crontabRouter := systemRouter.Group("crontab")
	crontabApi := v1.ApiGroupApp.SystemApiGroup.CrontabApi
	{
		crontabRouter.GET("list", crontabApi.GetCronList)
		crontabRouter.POST("run", crontabApi.RunByTag)
		crontabRouter.POST("add", crontabApi.AddCron)
		crontabRouter.POST("disable", crontabApi.DisabledCron)
		crontabRouter.POST("enable", crontabApi.EnabledCron)
		crontabRouter.POST("delete", crontabApi.DeleteCron)
		crontabRouter.POST("update", crontabApi.UpdateCron)

	}
	//高德相关接口
	amapRouter := systemRouter.Group("amap")
	amapApi := v1.ApiGroupApp.SystemApiGroup.AmapApi
	{
		amapRouter.GET("ip", amapApi.IpToCity)
		amapRouter.GET("weather", amapApi.GetWeather)
	}
	return systemRouter
}
