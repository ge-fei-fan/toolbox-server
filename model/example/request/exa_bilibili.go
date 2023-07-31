package request

import (
	"toolbox-server/model/common/request"
	"toolbox-server/model/example"
)

type ExaVideoSearch struct {
	example.VideoInstance
	//Downloading bool `json:"downloading" form:"downloading"`
	request.PageInfo
}

type ExaCollectVideoSearch struct {
	example.CollectVideo
	//Downloading bool `json:"downloading" form:"downloading"`
	request.PageInfo
}

type ExaCollectUserSearch struct {
	example.BilibiliCollect
	//Downloading bool `json:"downloading" form:"downloading"`
	request.PageInfo
}
