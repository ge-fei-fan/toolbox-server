package initialize

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"toolbox-server/global"
	"toolbox-server/model/example"
	"toolbox-server/model/system"
)

// Gorm 初始化数据库并产生数据库全局变量
// Author SliverHorn
func Gorm() *gorm.DB {
	switch global.TOOL_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "sqlite":
		return GormSqlite()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.TOOL_DB
	err := db.AutoMigrate(
		example.VideoInstance{},
		example.BilibiliCollect{},
		example.CollectVideo{},
		system.SysCrontab{},
		system.Weather{},
	)
	if err != nil {
		global.TOOL_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.TOOL_LOG.Info("register table success")
}
