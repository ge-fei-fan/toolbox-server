package internal

import (
	"fmt"
	"gorm.io/gorm/logger"
	"toolbox-server/global"
)

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
// Author [SliverHorn](https://github.com/SliverHorn)
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志
// Author [SliverHorn](https://github.com/SliverHorn)
func (w *writer) Printf(message string, data ...interface{}) {
	var logZap bool
	switch global.TOOL_CONFIG.System.DbType {
	case "mysql":
		logZap = global.TOOL_CONFIG.Mysql.LogZap
	}
	if logZap {
		global.TOOL_LOG.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
