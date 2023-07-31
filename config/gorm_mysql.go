package config

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func initMysql() Mysql {
	return Mysql{
		GeneralDB: initGeneralDB(),
	}
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func (m *Mysql) GetLogMode() string {
	return m.LogMode
}
