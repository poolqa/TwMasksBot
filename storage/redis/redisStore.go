package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/poolqa/log"
	"time"
)

func init() {

}

type RedisStore struct {
	pool redis.Pool
}

type RedisConn struct {
	conn redis.Conn
}

func NewRedisStore(ip string, port string, password string, selectDb int) *RedisStore {

	if port == "" {
		port = ":6379"
	}

	log.Info("connect to redis : ", ip+":"+port, selectDb)

	pool := &redis.Pool{
		MaxIdle:     100,               //最大空闲数，数据库连接的最大空闲时间。超过空闲时间，数据库连接将被标记为不可用，然后被释放。设为0表示无限制
		MaxActive:   0,                 //最大连接数，0为不限制
		IdleTimeout: 240 * time.Second, //最大建立连接等待时间。
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ip+":"+port)
			if err != nil {
				log.Error("打开数据库连接失败：" + err.Error())
				panic(err)
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					log.Error("打开数据库连接失败：" + err.Error())
					panic(err)
				}
			}
			if selectDb > 0 {
				if _, err := c.Do("SELECT", selectDb); err != nil {
					c.Close()
					log.Error("打开数据库连接失败：" + err.Error())
					panic(err)
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			log.Error("打开数据库连接失败：" + err.Error())
			panic(err)
		},
	}

	return &RedisStore{*pool}
}

func NewRedisConn(rs *RedisStore) *RedisConn {

	return &RedisConn{
		conn: rs.pool.Get(),
	}
}

//订阅频道
func (rs *RedisConn) Subscribe(channel string, send func(receiveMsg string)) error {

	log.Info("订阅频道: " + channel)

	psc := redis.PubSubConn{rs.conn}
	err := psc.Subscribe(channel)

	if err != nil {
		return err
	}
	for {
		switch v := psc.Receive().(type) {

		case redis.Message:
			//收到的订阅信息

			send(string(v.Data))
			log.Info(v.Channel, v.Data, string(v.Data))

		case redis.Subscription:
			//订阅时返回的信息
			log.Info(v.Channel, v.Kind, v.Count)

		case error:
			log.Info("Subscribe2:", v)
			return v
		}
	}
}

//取消订阅
func (rs *RedisConn) Unsubscribe(channel string) error {
	log.Info("取消订阅频道: " + channel)
	psc := redis.PubSubConn{rs.conn}
	err := psc.Unsubscribe(channel)
	return err
}

//发布消息
func (rs *RedisConn) PublishMessage(channel string, message string) error {
	log.Info("发布消息: "+channel, message)
	psc := redis.PubSubConn{rs.conn}
	_, err := psc.Conn.Do("PUBLISH", channel, message)
	return err
}

func (rs *RedisConn) Get(key string) (value string, err error) {
	value, err = redis.String(rs.conn.Do("Get", key))
	return
}

func (rs *RedisConn) Set(key string, value interface{}, expire int64) error {
	_, err := rs.conn.Do("Set", key, value, "EX", expire)
	return err
}

func (rs *RedisConn) Del(key string) (value interface{}, err error) {
	value, err = rs.conn.Do("Del", key)
	return
}

func (rs *RedisConn) Pop(key string) (value interface{}, err error) {
	err = rs.conn.Send("MULTI")
	if err != nil {
		return
	}
	err = rs.conn.Send("Get", key)
	if err != nil {
		return
	}
	err = rs.conn.Send("Del", key)
	if err != nil {
		return
	}
	var values []interface{}
	values, err = redis.Values(rs.conn.Do("EXEC"))
	if err != nil {
		return
	}
	if values[0] != nil {
		value, _ = redis.String(values[0], nil)
	}
	return
}

func (rs *RedisConn) HSet(key string, member string, value interface{}) error {
	_, err := rs.conn.Do("HSET", key, member, value)
	return err
}

func (rs *RedisConn) HSetnx(key string, member string, value interface{}) (int64, error) {
	ret, err := rs.conn.Do("HSETNX", key, member, value)
	return ret.(int64), err
}

func (rs *RedisConn) HIncrBy(key string, member string) error {
	_, err := rs.conn.Do("HINCRBY", key, member, 1)
	return err
}

/**
Zadd 命令用于将一个或多个成员元素及其分数值加入到有序集当中
*/
func (rs *RedisConn) ZAdd(key string, sort int64, value interface{}) error {
	_, err := rs.conn.Do("ZADD", key, sort, value)
	return err
}

/**
计算集合中元素的数量。
*/
func (rs *RedisConn) ZCard(key string) (value interface{}, err error) {
	value, err = rs.conn.Do("ZCARD", key)
	return
}

func (rs *RedisConn) LPush(key string, values ...interface{}) error {
	params := []interface{}{key}
	params = append(params, values...)
	_, err := rs.conn.Do("lpush", params...)
	return err
}

func (rs *RedisConn) RPush(key string, values ...interface{}) error {
	params := []interface{}{key}
	params = append(params, values...)
	_, err := rs.conn.Do("rpush", params...)
	return err
}
