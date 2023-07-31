package system

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"toolbox-server/global"
	"toolbox-server/model/system"
)

type Amap struct {
}

func (a *Amap) IpToCity() (err error) {
	if global.TOOL_CONFIG.Amap.Key == "" {
		return errors.New("高德key为空")
	}
	ipUrl := fmt.Sprintf("https://restapi.amap.com/v3/ip?key=%s", global.TOOL_CONFIG.Amap.Key)
	client := resty.New()
	itc := &system.IpToCity{}
	_, err = client.R().SetHeader("User-Agent", global.UserAgent).SetResult(itc).Get(ipUrl)
	if err != nil {
		return err
	}
	global.TOOL_CONFIG.Amap.CityCode = itc.Adcode
	err = global.WriteConfig()
	if err != nil {
		global.TOOL_LOG.Error("WriteConfig err:", zap.Error(err))
		return
	}
	return nil
}
func (a *Amap) GetWeather() (weather *system.Weather, err error) {
	var w system.Weather
	err = global.TOOL_DB.Order("created_at DESC").First(&w).Error
	if err != nil {
		global.TOOL_LOG.Error("获取实况天气失败", zap.Error(err))
		return nil, err
	}
	return &w, nil
}
