package core

import (
	"fmt"
	"toolbox-server/initialize"

	"go.uber.org/zap"
	"time"
	"toolbox-server/global"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	//if global.TOOL_CONFIG.System.UseMultipoint || global.TOOL_CONFIG.System.UseRedis {
	//	// 初始化redis服务
	//	initialize.Redis()
	//}
	Router := initialize.Routers()
	//Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.TOOL_CONFIG.System.Addr)
	//s := initServer(address, Router)
	global.TOOL_SERVER = initServer(address, Router)

	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.TOOL_LOG.Info("server run success on ", zap.String("address", address))

	global.TOOL_LOG.Error(global.TOOL_SERVER.ListenAndServe().Error())
}
