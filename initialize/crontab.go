package initialize

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"time"
	"toolbox-server/global"
	"toolbox-server/model/system"
	"toolbox-server/utils/cron"
)

func InitCron() {

	global.TOOL_SCHEDULER = gocron.NewScheduler(time.Local)
	//global.TOOL_SCHEDULER.SingletonModeAll()
	global.TOOL_SCHEDULER.TagsUnique()
	gocron.SetPanicHandler(func(jobName string, recoverData interface{}) {
		global.TOOL_LOG.Error(fmt.Sprintf("Panic in job: %s %s", jobName, recoverData))
	})
	InitTask()
	crons := make([]*system.SysCrontab, 0)
	err := global.TOOL_DB.Where("status = ?", 0).Find(&crons).Error
	if err != nil {
		global.TOOL_LOG.Error("查询启用的定时任务失败", zap.Error(err))
		return
	}
	for _, c := range crons {
		err = cron.AddJob(c)
		if err != nil {
			global.TOOL_LOG.Error(c.Tag+"添加定时任务失败", zap.Error(err))
		}
	}
	global.TOOL_SCHEDULER.StartAsync()
}

func InitTask() {
	cron.Tasks["AutoCollect"] = cron.AutoCollect
	cron.Tasks["FreshToken"] = cron.ReFreshToken
	cron.Tasks["weather"] = cron.AutoGetWeather
}
