package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"toolbox-server/global"
	"toolbox-server/model/common/response"
	"toolbox-server/model/system"
	"toolbox-server/model/system/request"
)

type CrontabApi struct{}

func (s *CrontabApi) GetCronList(c *gin.Context) {
	var pageInfo request.SysCrontabSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := crontabService.GetCronList(pageInfo)
	if err != nil {
		global.TOOL_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}
func (s *CrontabApi) RunByTag(c *gin.Context) {
	var cron system.SysCrontab
	err := c.ShouldBind(&cron)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cron.Tag == "" {
		response.FailWithMessage("tag为空", c)
		return
	}
	err = crontabService.RunByTag(cron)
	if err != nil {
		global.TOOL_LOG.Error("运行任务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("执行定时任务成功", c)

}
func (s *CrontabApi) AddCron(c *gin.Context) {
	var cron system.SysCrontab
	err := c.ShouldBind(&cron)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cron.Tag == "" {
		response.FailWithMessage("tag为空", c)
		return
	}
	err = crontabService.AddCron(cron)
	if err != nil {
		global.TOOL_LOG.Error("创建任务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
func (s *CrontabApi) DisabledCron(c *gin.Context) {
	var cron system.SysCrontab
	err := c.ShouldBind(&cron)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cron.ID == 0 {
		response.FailWithMessage("ID为空", c)
		return
	}
	err = crontabService.DisabledCron(cron)
	if err != nil {
		global.TOOL_LOG.Error("禁用任务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
func (s *CrontabApi) EnabledCron(c *gin.Context) {
	var cron system.SysCrontab
	err := c.ShouldBind(&cron)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cron.ID == 0 {
		response.FailWithMessage("ID为空", c)
		return
	}
	err = crontabService.EnabledCron(cron)
	if err != nil {
		global.TOOL_LOG.Error("启用任务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
func (s *CrontabApi) UpdateCron(c *gin.Context) {
	var cron system.SysCrontab
	err := c.ShouldBind(&cron)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cron.ID == 0 {
		response.FailWithMessage("ID为空", c)
		return
	}
	err = crontabService.UpdateCron(cron)
	if err != nil {
		global.TOOL_LOG.Error("更新任务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
func (s *CrontabApi) DeleteCron(c *gin.Context) {
	var cron system.SysCrontab
	err := c.ShouldBind(&cron)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if cron.ID == 0 {
		response.FailWithMessage("ID为空", c)
		return
	}
	err = crontabService.DeleteCron(cron)
	if err != nil {
		global.TOOL_LOG.Error("删除任务失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.Ok(c)
}
