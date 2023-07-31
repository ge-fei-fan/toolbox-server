package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"path/filepath"
	myconfig "toolbox-server/config"
	"toolbox-server/core/internal"
	"toolbox-server/global"
	"toolbox-server/utils"
)

func Viper(path ...string) *viper.Viper {
	var config string
	switch gin.Mode() {
	case gin.DebugMode:
		config = internal.ConfigDebugFile
		fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigDebugFile)
	case gin.ReleaseMode:
		config = internal.ConfigDefaultFile
		fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigDefaultFile)
	}
	v := viper.New()
	//v.SetConfigFile(config)
	v.SetConfigName(config)
	v.SetConfigType("yaml")
	appPath, _ := utils.AppConfigPath()
	configPath := filepath.Join(appPath, myconfig.AppName)
	//v.AddConfigPath(".")
	v.AddConfigPath(configPath)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			//创建默认配置文件
			defaultConfig, err := json.Marshal(myconfig.InitDefaultServer())
			if err != nil {
				panic(fmt.Errorf("Fatal error Marshal default config file: %s \n", err))
			}
			err = v.ReadConfig(bytes.NewBuffer(defaultConfig))
			if err != nil {
				panic(fmt.Errorf("Fatal error read default config file: %s \n", err))
			}
			err = v.SafeWriteConfig()
			if err != nil {
				panic(fmt.Errorf("Fatal error write default config file: %s \n", err))
			}
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.TOOL_CONFIG); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&global.TOOL_CONFIG); err != nil {
		fmt.Println(err)
	}

	return v
}
