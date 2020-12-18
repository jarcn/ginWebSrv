package conf

type MysqlConfig struct {
	DriverName   string
	Dsn          string
	ShowSql      bool
	ShowExecTime bool
	MaxIdle      int
	MaxOpen      int
}

var Db = map[string]MysqlConfig{
	"mysql": {
		DriverName:   "mysql",
		Dsn:          "root:mysql@tcp(127.0.0.1:3306)/systemdb?charset=utf8mb4&parseTime=true&loc=Local",
		ShowSql:      false,
		ShowExecTime: false,
		MaxIdle:      10,
		MaxOpen:      200,
	},
}
