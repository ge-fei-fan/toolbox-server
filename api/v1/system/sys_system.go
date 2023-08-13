package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"toolbox-server/config"
	"toolbox-server/global"
	"toolbox-server/model/common/response"
)

type SystemApi struct{}

func (s *SystemApi) Shutdown(c *gin.Context) {
	go func() {
		systemService.ShutdownServer()
	}()
	response.Ok(c)
}
func (s *SystemApi) Version(c *gin.Context) {
	data := struct {
		Version string `json:"version"`
	}{
		Version: global.TOOL_APP_VERSION,
	}
	response.OkWithData(data, c)
}
func (s *SystemApi) Mode(c *gin.Context) {
	data := struct {
		Mode string `json:"mode"`
	}{
		Mode: gin.Mode(),
	}
	response.OkWithData(data, c)
}
func (s *SystemApi) Test(c *gin.Context) {
	response.OkWithMessage("哈哈哈", c)
}

func (s *SystemApi) GetConfig(c *gin.Context) {
	response.OkWithData(global.TOOL_CONFIG, c)
}

func (s *SystemApi) UpdateConfig(c *gin.Context) {
	var conf config.Server
	err := c.ShouldBind(&conf)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemService.UpdateSysConfig(&conf)
	if err != nil {
		global.TOOL_LOG.Error("更新配置文件错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
