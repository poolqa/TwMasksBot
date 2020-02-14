package config

type Config struct {
	Redis    *redis
	Mysql    *MysqlConfig
	Telegram *TgConfig
	Line     *LineConfig
}

type MysqlConfig struct {
	Url         string
	User        string
	Password    string
	Database    string
	Charset     string
	MaxIdleConn int
	MaxOpenConn int
	LogLevel    string
}

type redis struct {
	Ip       string
	Port     string
	Password string
	SelectDb int
}

type TgConfig struct {
	Enable bool
	Token  string
	Debug  bool
}

type LineConfig struct {
	Enable        bool
	Port          string
	ChannelSecret string
	Token         string
	Debug         bool
}
