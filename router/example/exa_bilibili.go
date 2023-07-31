package example

import (
	"github.com/gin-gonic/gin"
	v1 "toolbox-server/api/v1"
	"toolbox-server/middleware"
)

type BiliRouter struct{}

func (s *BiliRouter) InitBilibiliRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	biliRouter := Router.Group("bilibili")
	biliRrRouter := biliRouter.Group("qrcode")
	collectRouter := biliRouter.Group("collect")
	biliApi := v1.ApiGroupApp.ExampleApiGroup.BilibiliApi
	{
		biliRouter.GET("download", biliApi.Download)
		biliRouter.GET("redownload", biliApi.ReDownload)
		biliRouter.GET("video", biliApi.GetVideoList)
		biliRouter.GET("videodetail", biliApi.VideoDetail)
		biliRouter.POST("delete", biliApi.DeleteVideo) //删除视频
		biliRouter.Use(middleware.DefaultLimit(5, 60)).GET("account", biliApi.Account)
		//biliRouter.GET("account", biliApi.Account)
	}
	{
		biliRrRouter.GET("generate", biliApi.Generate) //获取登录二维码
		biliRrRouter.GET("poll", biliApi.Poll)         //获取登录状态
	}
	{
		collectRouter.POST("add", biliApi.AddCollectUser)
		collectRouter.GET("list", biliApi.CollectUserList)
		collectRouter.GET("video", biliApi.CollectVideoList)

		collectRouter.GET("status", biliApi.CollectUserStatus)
	}
	return biliRouter
}
