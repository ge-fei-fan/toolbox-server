package system

import (
	"context"
	"time"
	"toolbox-server/config"
	"toolbox-server/global"
)

type System struct {
}

func (s *System) ShutdownServer() {
	ctx, channel := context.WithTimeout(context.Background(), 5*time.Second)
	defer channel()
	if err := global.TOOL_SERVER.Shutdown(ctx); err != nil {
		global.TOOL_LOG.Warn("server shutdown error")
	}
}
func (s *System) UpdateSysConfig(conf *config.Server) error {
	//日志等级 日志路径 服务器端口
	//bilbili 粘贴板 采集 下载路径
	if conf.Zap.Level != "" {
		global.TOOL_CONFIG.Zap.Level = conf.Zap.Level
	}
	if conf.Zap.Director != "" {
		global.TOOL_CONFIG.Zap.Director = conf.Zap.Director
	}
	if conf.System.Addr != 0 {
		global.TOOL_CONFIG.System.Addr = conf.System.Addr
	}
	if conf.Bilibili.DownloadPath != "" {
		global.TOOL_CONFIG.Bilibili.DownloadPath = conf.Bilibili.DownloadPath
	}
	global.TOOL_CONFIG.Bilibili.Collect = conf.Bilibili.Collect
	global.TOOL_CONFIG.Bilibili.AutoDownload = conf.Bilibili.AutoDownload
	err := global.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
