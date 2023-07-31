package config

import (
	"github.com/gin-gonic/gin"
	"path/filepath"
	"toolbox-server/utils"
)

type Sqlite struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func initSqlite() Sqlite {
	appPath, _ := utils.AppConfigPath()
	var dbPath string
	switch gin.Mode() {
	case gin.DebugMode:
		dbPath = filepath.Join(appPath, AppName, AppName+"-debug.db")
	case gin.ReleaseMode:
		dbPath = filepath.Join(appPath, AppName, AppName+".db")
	}

	return Sqlite{
		GeneralDB: GeneralDB{
			Path:         dbPath,
			Dbname:       "toolbox",
			MaxIdleConns: 10,
			MaxOpenConns: 100,
			LogMode:      "info",
			Config:       "charset=utf8mb4&parseTime=True&loc=Local",
		},
	}
}

// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func (m *Sqlite) Dsn() string {
	return m.Path
}

func (m *Sqlite) GetLogMode() string {
	return m.LogMode
}
