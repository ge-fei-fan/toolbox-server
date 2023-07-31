package config

type Server struct {
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	System   System   `mapstructure:"system" json:"system" yaml:"system"`
	Mysql    Mysql    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Sqlite   Sqlite   `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
	Local    Local    `mapstructure:"local" json:"local" yaml:"local"`
	Bilibili Bilibili `mapstructure:"bilibili" json:"bilibili" yaml:"bilibili"`
	// 跨域配置
	Cors CORS `mapstructure:"cors" json:"cors" yaml:"cors"`
	Amap Amap `mapstructure:"amap" json:"amap" yaml:"amap"`
}

func InitDefaultServer() *Server {
	return &Server{
		Zap:      initDefaultZap(),
		System:   initDefaultSystem(),
		Mysql:    initMysql(),
		Sqlite:   initSqlite(),
		Local:    initDefaultLocal(),
		Bilibili: initDefaultBilibili(),
		Cors:     initDefualtCors(),
		Amap:     Amap{},
	}
}
