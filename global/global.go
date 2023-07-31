package global

import (
	"bytes"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"toolbox-server/config"
)

var (
	TOOL_DB        *gorm.DB
	TOOL_VP        *viper.Viper
	TOOL_CONFIG    config.Server
	TOOL_LOG       *zap.Logger
	TOOL_SERVER    *http.Server
	TOOL_FFMPEG    string
	TOOL_SCHEDULER *gocron.Scheduler
)

const (
	UserAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	TOOL_APP_VERSION = "V1.0.2"
)

func WriteConfig() error {
	myConfig, err := json.Marshal(TOOL_CONFIG)
	if err != nil {
		return err
	}
	err = TOOL_VP.ReadConfig(bytes.NewBuffer(myConfig))
	if err != nil {
		return err
	}
	err = TOOL_VP.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
