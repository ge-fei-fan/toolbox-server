package config

type Amap struct {
	Key      string `mapstructure:"key" json:"key" yaml:"key"`
	CityCode string `mapstructure:"city_code" json:"city_code" yaml:"city_code"`
}
