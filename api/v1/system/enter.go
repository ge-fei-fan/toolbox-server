package system

import "toolbox-server/service"

type ApiGroup struct {
	SystemApi
	CrontabApi
	AmapApi
}

var (
	systemService  = service.ServiceGroupApp.SystemServiceGroup.System
	crontabService = service.ServiceGroupApp.SystemServiceGroup.Crontab
	amapService    = service.ServiceGroupApp.SystemServiceGroup.Amap
)
