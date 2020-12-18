package conf

import (
	"sync"
)

type Config struct {
	Language string
	Token    string
	Super    string
	RedisPre string
	Host     string
	Routers  []string
	Env      string
}

var (
	Cfg     Config
	mutex   sync.Mutex //互斥锁
	declare sync.Once  //只执行一次
)

func Set(cfg Config) {
	mutex.Lock()
	Cfg.RedisPre = setDefault(cfg.RedisPre, "", "go.admin.redis")
	Cfg.Language = setDefault(cfg.Language, "", "cn")
	Cfg.Token = setDefault(cfg.Token, "", "token")
	Cfg.Host = setDefault(cfg.Host, "", "http://localhost:8080")
	Cfg.Super = setDefault(cfg.Super, "", "admin")
	Cfg.Env = setDefault(cfg.Env, "", "dev")
	Cfg.Routers = cfg.Routers
	mutex.Unlock()
}

func setDefault(value, def, defValue string) string {
	switch value == def {
	case true:
		return defValue
	default:
		return value
	}
}
