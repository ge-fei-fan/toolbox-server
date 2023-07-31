package request

import (
	"toolbox-server/model/common/request"
	"toolbox-server/model/system"
)

type SysCrontabSearch struct {
	system.SysCrontab
	request.PageInfo
}
