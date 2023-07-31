package example

import "toolbox-server/service"

type ApiGroup struct {
	BilibiliApi
}

var (
	biliQrcodeService = service.ServiceGroupApp.ExampleServiceGroup.BilibiliQrcode
	biliService       = service.ServiceGroupApp.ExampleServiceGroup.Bilibili
)
