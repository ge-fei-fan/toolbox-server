package config

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type System struct {
	Env    string `mapstructure:"env" json:"env" yaml:"env"`             // 环境值
	Addr   int    `mapstructure:"addr" json:"addr" yaml:"addr"`          // 端口值
	DbType string `mapstructure:"db-type" json:"db-type" yaml:"db-type"` // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	//OssType       string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`                   // Oss类型
	//UseRedis      bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                // 使用redis
	LimitCountIP int    `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"` //IP限制次数
	LimitTimeIP  int    `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`    //  IP限制时间
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
}

func initDefaultSystem() System {
	var port int
	switch gin.Mode() {
	case gin.DebugMode:
		port = 28888
	case gin.ReleaseMode:
		port = getRandPort()
		//for {
		//	port = getRandPort()
		//	if isPortAvailable(port) {
		//		port = getRandPort()
		//	} else {
		//		break
		//	}
		//}
	}

	return System{
		Env:          "develop",
		Addr:         port,
		DbType:       "sqlite",
		LimitCountIP: 15000,
		RouterPrefix: "",
	}
}
func getRandPort() int {
	rand.Seed(time.Now().UnixNano())
	// 生成随机数
	min := 25000
	max := 65000
	return rand.Intn(max-min+1) + min
}
func isPortAvailable(port int) bool {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}
