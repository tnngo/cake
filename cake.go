package cake

import (
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

type profiles struct {
	Active string
}

type profilesConfig struct {
	Profiles *profiles
}

type serverConfig struct {
	Port int
}

type mysqlConfig struct {
	Database     string
	IP           string
	Port         int
	User         string
	Password     string
	MaxOpenConns int `toml:"max_open_conns"`
	MaxIdleConns int `toml:"max_idle_conns"`
}

type redisConfig struct {
	DB       int
	Addr     string
	Password string

	PoolSize int `toml:"pool_size"`
}

type tomlPromise struct {
	Profiles *profiles
	Server   *serverConfig
	Mysql    []*mysqlConfig
	Redis    *redisConfig
}

var (
	_defaultTomlNames = []string{"cake.toml", "app.toml"}
)

func init() {
	zapConsole()
	var defaultTomlName string
	for _, v := range _defaultTomlNames {
		_, err := os.Stat(v)
		if err == nil {
			defaultTomlName = v
			break
		}
	}

	if len(defaultTomlName) == 0 {
		zap.L().Panic("未找到toml配置文件，默认toml配置文件名为: cake.toml, app.toml")
	}

	parseToml(defaultTomlName)
}

var (
	ServerConfig *serverConfig
	RedisConfig  *redisConfig
)

func parseToml(filename string) {
	pc := &profilesConfig{}
	_, err := toml.DecodeFile(filename, pc)
	if err != nil {
		zap.L().Error(err.Error())
		return
	}

	tp := new(tomlPromise)

	if pc.Profiles == nil || pc.Profiles.Active == "" {
		_, err := toml.DecodeFile(filename, tp)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
	} else {
		strs := strings.Split(filename, ".")
		_, err := toml.DecodeFile(strs[0]+"_"+pc.Profiles.Active+".toml", tp)
		if err != nil {
			zap.L().Error(err.Error())
			return
		}
	}

	ServerConfig = tp.Server

	if len(tp.Mysql) != 0 {
		loadGrom(tp.Mysql)
	}

	if tp.Redis != nil {
		loadRedis(tp.Redis)
	}
}
