package cron

import (
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"toolbox-server/global"
	"toolbox-server/model/system"
)

func AutoGetWeather(job gocron.Job) {
	if global.TOOL_CONFIG.Amap.CityCode == "" {
		global.TOOL_LOG.Error("城市代码为空")
		return
	}

	weatherUrl := "https://restapi.amap.com/v3/weather/weatherInfo"
	client := resty.New()
	sw := &system.WeatherLive{}
	_, err := client.R().SetResult(sw).SetQueryParams(
		map[string]string{
			"city":       global.TOOL_CONFIG.Amap.CityCode,
			"key":        global.TOOL_CONFIG.Amap.Key,
			"extensions": "base",
		}).Get(weatherUrl)
	w := system.Weather{
		WeatherInfo: sw.Lives[0],
	}
	if err != nil {
		global.TOOL_LOG.Error("查询实况天气数据出错", zap.Error(err))
		return
	}
	err = global.TOOL_DB.Create(&w).Error
	if err != nil {
		global.TOOL_LOG.Error("插入天气数据错误", zap.Error(err))
	}
	return

}
