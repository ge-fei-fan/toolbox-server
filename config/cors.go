package config

//跨域配置
//需要配合 server/initialize/router.go -> `Router.Use(middleware.CorsByRules())` 使用

type CORS struct {
	//放行模式: allow-all, 放行全部; whitelist, 白名单模式, 来自白名单内域名的请求添加 cors 头; strict-whitelist 严格白名单模式, 白名单外的请求一律拒绝
	Mode      string          `mapstructure:"mode" json:"mode" yaml:"mode"`
	Whitelist []CORSWhitelist `mapstructure:"whitelist" json:"whitelist" yaml:"whitelist"`
}

type CORSWhitelist struct {
	AllowOrigin      string `mapstructure:"allow-origin" json:"allow-origin" yaml:"allow-origin"`
	AllowMethods     string `mapstructure:"allow-methods" json:"allow-methods" yaml:"allow-methods"`
	AllowHeaders     string `mapstructure:"allow-headers" json:"allow-headers" yaml:"allow-headers"`
	ExposeHeaders    string `mapstructure:"expose-headers" json:"expose-headers" yaml:"expose-headers"`
	AllowCredentials bool   `mapstructure:"allow-credentials" json:"allow-credentials" yaml:"allow-credentials"`
}

func initDefualtCors() CORS {
	return CORS{
		Mode: "allow-all",
		Whitelist: []CORSWhitelist{
			{
				AllowCredentials: true,
				AllowOrigin:      "example.com",
				AllowHeaders:     "content-type",
				AllowMethods:     "GET, POST",
				ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type",
			},
		},
	}
}
