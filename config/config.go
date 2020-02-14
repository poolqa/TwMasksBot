package config

import (
	"errors"
	"github.com/Unknwon/goconfig"
)

var GConfig *Config

func ParseConfig(configFile string) (*Config, error) {
	var err error
	cfg, err := goconfig.LoadConfigFile(configFile) //读取配置文件，并返回其Config
	if err != nil {
		err = errors.New("找不到配置文件！信息：" + err.Error())
		return nil, err
	}

	config := new(Config)

	//redis配置
	config.Redis, err = parseRedis(cfg)
	if err != nil {
		return nil, err
	}

	//mysql
	config.Mysql = &MysqlConfig{}
	if config.Mysql, err = parseMysql(cfg, "db"); err != nil {
		return nil, err
	}

	config.Telegram, err = parseTgConf(cfg)
	if err != nil {
		return nil, err
	}

	config.Line, err = parseLineConf(cfg)
	if err != nil {
		return nil, err
	}

	return config, err
}

func parseRedis(cfg *goconfig.ConfigFile) (*redis, error) {
	var err error
	redis := new(redis)
	if _, err := cfg.GetSection("redis"); err == nil {
		var err error
		redis.Ip, err = cfg.GetValue("redis", "ip")
		if err != nil {
			return nil, errors.New("配置文件中无法找到redis的IP配置信息！")
		}
		redis.Port, err = cfg.GetValue("redis", "port")
		if err != nil {
			return nil, errors.New("配置文件中无法找到redis的port配置信息！")
		}
		redis.Password, err = cfg.GetValue("redis", "password")
		if err != nil {
			return nil, errors.New("配置文件中无法找到redis的password配置信息！")
		}
		redis.SelectDb, err = cfg.Int("redis", "selectDb")
		if err != nil {
			redis.SelectDb = 1
		}
	} else {
		err = errors.New("配置文件中无法找到redis信息。")
	}
	return redis, err
}

func parseMysql(cfg *goconfig.ConfigFile, dbName string) (*MysqlConfig, error) {
	var err error
	mysql := new(MysqlConfig)
	if _, err := cfg.GetSection(dbName); err == nil {
		var err error
		mysql.Url, err = cfg.GetValue(dbName, "url")
		if err != nil {
			return nil, errors.New("配置文件中无法找到数据库:" + dbName + "的url配置信息！")
		}
		mysql.User, err = cfg.GetValue(dbName, "user")
		if err != nil {
			return nil, errors.New("配置文件中无法找到数据库:" + dbName + "的USER配置信息！")
		}
		mysql.Password, err = cfg.GetValue(dbName, "password")
		if err != nil {
			return nil, errors.New("配置文件中无法找到数据库:" + dbName + "的PASSWORD配置信息！")
		}
		mysql.Database, err = cfg.GetValue(dbName, "database")
		if err != nil {
			return nil, errors.New("配置文件中无法找到数据库:" + dbName + "的DATABASE配置信息！")
		}
		mysql.Charset, err = cfg.GetValue(dbName, "charset")
		if err != nil {
			return nil, errors.New("配置文件中无法找到数据库:" + dbName + "的CHARSET配置信息！")
		}
		mysql.MaxIdleConn, _ = cfg.Int(dbName, "max_idle_conn")

		mysql.MaxOpenConn, _ = cfg.Int(dbName, "max_open_conn")

		mysql.LogLevel, _ = cfg.GetValue(dbName, "log_level")

	} else {
		return nil, errors.New("配置文件中无法找到数据库:" + dbName + "的配置信息！")
	}
	return mysql, err
}

func parseTgConf(cfg *goconfig.ConfigFile) (*TgConfig, error) {
	section := "telegram"
	var err error
	conf := new(TgConfig)
	if _, err := cfg.GetSection(section); err == nil {
		var err error
		conf.Enable, _ = cfg.Bool(section, "enable")
		if conf.Enable {
			conf.Token, err = cfg.GetValue(section, "token")
			if err != nil {
				return nil, errors.New("配置文件中无法找到telegram的token配置信息！")
			}

			conf.Debug, err = cfg.Bool(section, "debug")
		}
	} else {
		conf.Enable = false
	}
	return conf, err
}

func parseLineConf(cfg *goconfig.ConfigFile) (*LineConfig, error) {
	section := "line"
	var err error
	conf := new(LineConfig)
	if _, err := cfg.GetSection(section); err == nil {
		var err error
		conf.Enable, _ = cfg.Bool(section, "enable")
		if conf.Enable {
			conf.Port, err = cfg.GetValue(section, "port")
			if err != nil {
				return nil, errors.New("配置文件中无法找到line的port配置信息！")
			}
			conf.ChannelSecret, err = cfg.GetValue(section, "channelSecret")
			if err != nil {
				return nil, errors.New("配置文件中无法找到line的channelSecret配置信息！")
			}
			conf.Token, err = cfg.GetValue(section, "token")
			if err != nil {
				return nil, errors.New("配置文件中无法找到line的token配置信息！")
			}

			conf.Debug, err = cfg.Bool(section, "debug")
		}
	} else {
		conf.Enable = false
	}
	return conf, err
}
