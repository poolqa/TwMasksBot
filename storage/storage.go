package storage

import (
	"../config"
	"./mysql"
	"./redis"
)

var GStorage *Storage

type Storage struct {
	config *config.Config
	db     *mysql.MysqlStore
	redis  *redis.RedisStore
}

func NewStorage(config *config.Config) *Storage {
	server := new(Storage)
	server.config = config
	server.db = mysql.NewMysqlStore(config.Mysql)

	server.redis = redis.NewRedisStore(config.Redis.Ip, config.Redis.Port, config.Redis.Password, config.Redis.SelectDb)
	redis.NewRedisConn(server.redis)

	return server
}

func (st *Storage) GetConfig() *config.Config {
	return st.config
}

/**
返回redis链接
*/
func (st *Storage) GetRedisConn() (redisConn *redis.RedisConn) {
	if st.config.Redis.Ip == "" {
		return
	}
	redisConn = redis.NewRedisConn(st.redis)
	return
}

func (st *Storage) GetDB() *mysql.MysqlStore {
	return st.db
}
