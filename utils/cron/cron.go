package cron

import (
	"github.com/go-co-op/gocron"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
	"toolbox-server/global"
	"toolbox-server/model/system"
)

var Tasks = make(map[string]func(job gocron.Job), 20)

func AddJob(c *system.SysCrontab) error {
	f, has := Tasks[c.Func]
	if !has {
		return errors.New("未找到函数")
	}
	job, err := global.TOOL_SCHEDULER.Cron(c.Cron).Tag(c.Tag).DoWithJobDetails(f)
	//_, err = global.TOOL_SCHEDULER.Cron(c.Cron).Tag(c.Tag).Do(AutoCollect)
	if err != nil {
		return err
	}

	job.SetEventListeners(func() {
		global.TOOL_LOG.Info("定时任务开始执行：" + time.Now().String())
	}, func() {
		global.TOOL_LOG.Info("定时任务执行结束：" + time.Now().String())
	})
	return nil
}

func UpdateJob(c system.SysCrontab) error {
	var tmp system.SysCrontab
	err := global.TOOL_DB.Where("id = ?", c.ID).First(&tmp).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到任务")
	}
	if err != nil {
		return err
	}
	jobs, err := global.TOOL_SCHEDULER.FindJobsByTag(tmp.Tag)
	if err != nil {
		return err
	}
	_, err = global.TOOL_SCHEDULER.Job(jobs[0]).Cron(c.Cron).Update()
	if err != nil {
		_ = AddJob(&tmp)
		return err
	}
	if tmp.Tag != c.Tag {
		jobs[0].Untag(tmp.Tag)
		jobs[0].Tag(c.Tag)
	}
	return nil
}

//func InitTask() {
//	tasks["AutoCollect"] = AutoCollect
//	tasks["FreshToken"] = ReFreshToken
//}
