package system

import "toolbox-server/global"

type SysCrontab struct {
	global.TOOL_MODEL
	Func      string `json:"func" form:"func" gorm:"column:func;comment:任务"`
	Tag       string `json:"tag" form:"tag" gorm:"column:tag"`
	Status    int    `json:"status" form:"status" gorm:"column:status; comment:启用0 禁用1" `
	NextTime  string `json:"nextTime" form:"nextTime" gorm:"column:next_time; comment:下次运行时间"`
	Cron      string `json:"cron" form:"cron" gorm:"column:cron;comment:执行间隔"`
	IsRunning bool   `json:"isRunning" gorm:"column:-; comment:0 未执行 1 正在执行"`
}
