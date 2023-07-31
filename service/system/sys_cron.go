package system

import (
	"errors"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
	"toolbox-server/global"
	"toolbox-server/model/system"
	"toolbox-server/model/system/request"
	"toolbox-server/utils/cron"
)

type Crontab struct {
}

func (c *Crontab) GetCronList(info request.SysCrontabSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize                    //5  5 5
	offset := info.PageSize * (info.Page - 1) //0 5 10
	// 创建db
	db := global.TOOL_DB.Model(&system.SysCrontab{})
	var sysCrontabs []system.SysCrontab
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if info.ID != 0 {
		db = db.Where("`ID` =?", info.ID)
	}
	err = db.Limit(limit).Offset(offset).Order("ID desc").Find(&sysCrontabs).Error
	for i, _ := range sysCrontabs {
		if sysCrontabs[i].Status == 0 {
			job, _ := global.TOOL_SCHEDULER.FindJobsByTag(sysCrontabs[i].Tag)
			if job == nil {
				global.TOOL_LOG.Error(sysCrontabs[i].Tag + "任务在数据库启用，但是未加入定时任务管理中")
				continue
			}
			sysCrontabs[i].IsRunning = job[0].IsRunning()
			sysCrontabs[i].NextTime = job[0].NextRun().Format("2006-01-02 15:04:05")
		}
	}
	return sysCrontabs, total, err

}
func (c *Crontab) AddCron(sysCrontab system.SysCrontab) (err error) {
	err = cron.AddJob(&sysCrontab)
	if err != nil {
		return
	}
	err = global.TOOL_DB.Create(&sysCrontab).Error
	if err != nil {
		return
	}
	return
}

// 禁用定时任务
func (c *Crontab) DisabledCron(sysCrontab system.SysCrontab) (err error) {

	err = global.TOOL_DB.Where("id = ?", sysCrontab.ID).First(&sysCrontab).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到任务")
	}
	if err != nil {
		return err
	}
	err = global.TOOL_SCHEDULER.RemoveByTag(sysCrontab.Tag)
	if err != nil && errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		return errors.New("任务已禁用")
	}
	sysCrontab.Status = 1
	sysCrontab.NextTime = ""
	sysCrontab.IsRunning = false
	err = global.TOOL_DB.Save(&sysCrontab).Error
	if err != nil {
		return err
	}
	return
}
func (c *Crontab) EnabledCron(sysCrontab system.SysCrontab) (err error) {

	err = global.TOOL_DB.Where("id = ?", sysCrontab.ID).First(&sysCrontab).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到任务")
	}
	if err != nil {
		return err
	}
	_, err = global.TOOL_SCHEDULER.FindJobsByTag(sysCrontab.Tag)
	if err != nil && errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		err = cron.AddJob(&sysCrontab)
		if err != nil {
			return err
		}
		sysCrontab.Status = 0
		err = global.TOOL_DB.Save(&sysCrontab).Error
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	return
}

// 删除定时任务
func (c *Crontab) DeleteCron(sysCrontab system.SysCrontab) (err error) {
	err = global.TOOL_DB.Where("id = ?", sysCrontab.ID).First(&sysCrontab).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到任务")
	}
	if err != nil {
		return err
	}
	err = global.TOOL_DB.Delete(&sysCrontab).Error
	if err != nil {
		return err
	}
	return
}

// 更新定时任务
func (c *Crontab) UpdateCron(sysCrontab system.SysCrontab) (err error) {

	err = cron.UpdateJob(sysCrontab)
	if err != nil {
		return err
	}
	err = global.TOOL_DB.Model(&sysCrontab).
		Updates(system.SysCrontab{Tag: sysCrontab.Tag, Cron: sysCrontab.Cron}).Error
	if err != nil {
		return err
	}
	return nil
}
func (c *Crontab) RunByTag(cron system.SysCrontab) (err error) {
	job, err := global.TOOL_SCHEDULER.FindJobsByTag(cron.Tag)
	if err != nil && errors.Is(err, gocron.ErrJobNotFoundWithTag) {
		return errors.New("未找到任务")
	}
	if job[0].IsRunning() {
		return errors.New("任务正在运行")
	}
	err = global.TOOL_SCHEDULER.RunByTag(cron.Tag)
	if err != nil {
		return err
	}
	return
}
