package initialize

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"toolbox-server/config"
	"toolbox-server/global"
	"toolbox-server/initialize/internal"
)

func GormSqlite() *gorm.DB {
	m := global.TOOL_CONFIG.Sqlite
	if m.Dbname == "" {
		return nil
	}

	if db, err := gorm.Open(sqlite.Open(m.Dsn()), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

func GormSqliteByConfig(m config.Sqlite) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	if db, err := gorm.Open(sqlite.Open(m.Dsn()), internal.Gorm.Config(m.Prefix, m.Singular)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}
