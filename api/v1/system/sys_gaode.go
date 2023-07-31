package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"toolbox-server/global"
	"toolbox-server/model/common/response"
)

type AmapApi struct{}

func (a *AmapApi) IpToCity(c *gin.Context) {
	err := amapService.IpToCity()
	if err != nil {
		global.TOOL_LOG.Error("获取城市失败!", zap.Error(err))
		response.FailWithMessage("获取城市失败", c)
		return
	}
	response.OkWithMessage("获取城市成功", c)
	return
}

func (a *AmapApi) GetWeather(c *gin.Context) {
	w, err := amapService.GetWeather()
	if err != nil {
		global.TOOL_LOG.Error("获取天气失败!", zap.Error(err))
		response.FailWithData(w, c)
		return
	}
	response.OkWithData(w, c)
	return
}
