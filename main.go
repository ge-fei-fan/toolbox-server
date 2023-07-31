package main

import (
	"go.uber.org/zap"
	"toolbox-server/core"
	"toolbox-server/global"
	"toolbox-server/initialize"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	global.TOOL_VP = core.Viper()
	global.TOOL_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.TOOL_LOG)
	global.TOOL_DB = initialize.Gorm() // gorm连接数据库
	if global.TOOL_DB != nil {
		initialize.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.TOOL_DB.DB()
		defer func() {
			err := db.Close()
			if err != nil {
				global.TOOL_LOG.Error("close db err:", zap.Error(err))
			}
		}()
		initialize.InitCron()
		global.TOOL_LOG.Info("定时任务初始化完成")
		err := initialize.InitDownloading()
		defer func() {
			global.TOOL_LOG.Info("处理视频")
			initialize.DealDownloading()
		}()
		if err != nil {
			global.TOOL_LOG.Error("init Downloading err:", zap.Error(err))
		}
		err = initialize.CheckFFmpeg()
		if err != nil {
			global.TOOL_LOG.Error("check ffmpeg err:", zap.Error(err))

		}

		core.RunWindowsServer()
	}

}
