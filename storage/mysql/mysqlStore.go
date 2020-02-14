package mysql

import (
	"../../config"
	//"../../entity"
	_ "github.com/go-sql-driver/mysql"
	"github.com/poolqa/log"
	"github.com/xormplus/xorm"
)

func init() {

}

type MysqlStore struct {
	Engine *xorm.Engine
}

func NewMysqlStore(conf *config.MysqlConfig) *MysqlStore {
	link := conf.User + ":" + conf.Password + "@tcp(" + conf.Url + ")/" + conf.Database +
		"?charset=" + conf.Charset + "&multiStatements=true"

	engine, err := xorm.NewEngine("mysql", link)
	if err != nil {
		log.Error("打开数据库连接失败：" + err.Error())
		panic(err)
	}
	//连接池设置
	engine.SetMaxOpenConns(conf.MaxOpenConn) //用于设置最大打开的连接数，默认值为0表示不限制
	engine.SetMaxIdleConns(conf.MaxIdleConn) //用于设置闲置的连接数。

	// log level

	//sync entity
	//engine.Sync2(new(entity.Subscription))

	return &MysqlStore{
		Engine: engine,
	}
}
