package system

import "toolbox-server/global"

// 根据ip获取定位
type IpToCity struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

type WeatherInfo struct {
	Province         string `json:"province" gorm:"column:province;comment:省份名"`
	City             string `json:"city" gorm:"column:city;comment:城市名"`
	Adcode           string `json:"adcode" gorm:"column:adcode;comment:区域编码"`
	Weather          string `json:"weather" gorm:"column:weather;comment:天气现象"`
	Temperature      string `json:"temperature" gorm:"column:temperature;comment:实时气温"`
	Winddirection    string `json:"winddirection" gorm:"column:winddirection;comment:风向描述"`
	Windpower        string `json:"windpower" gorm:"column:windpower;comment:风力级别"`
	Humidity         string `json:"humidity" gorm:"column:humidity;comment:空气湿度"`
	Reporttime       string `json:"reporttime" gorm:"column:reporttime;comment:数据发布的时间"`
	TemperatureFloat string `json:"temperature_float" gorm:"column:temperature_float;comment:实时气温float"`
	HumidityFloat    string `json:"humidity_float" gorm:"column:humidity_float;comment:空气湿度float"`
}

// 根据位置获取实况天气
type WeatherLive struct {
	Status   string        `json:"status"`
	Count    string        `json:"count"`
	Info     string        `json:"info"`
	Infocode string        `json:"infocode"`
	Lives    []WeatherInfo `json:"lives"`
}

// ======================================= dbmodel===================================

// 实况天气数据表
type Weather struct {
	global.TOOL_MODEL
	WeatherInfo `gorm:"embedded"`
}
