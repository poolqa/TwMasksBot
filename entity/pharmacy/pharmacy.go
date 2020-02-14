package pharmacy

import (
	"../../storage"
	"time"
)

type Pharmacy struct {
	Id          int64  `xorm:"'id' BIGINT(11) notnull pk autoincr"`
	Code        string `xorm:"'code' VARCHAR(20) DEFAULT ''"`
	Name        string `xorm:"'name' VARCHAR(64) DEFAULT ''"`
	Tel         string `xorm:"'tel' VARCHAR(20) DEFAULT ''"`
	Addr        string `xorm:"'addr' VARCHAR(200) DEFAULT ''"`
	Latitude    float64 `xorm:"'latitude' DECIMAL(14,7) DEFAULT '0.0'"`
	Longitude   float64 `xorm:"'longitude' DECIMAL(14,7) DEFAULT '0.0'"`
	AdultCount  int64  `xorm:"'adult_count' BIGINT(11) notnull DEFAULT '0'"`
	ChildCount  int64  `xorm:"'child_count' BIGINT(11) notnull DEFAULT '0'"`
	UpdTime     *time.Time `xorm:"'upd_time' timestamp NULL DEFAULT NULL "`
	SellRule    string `xorm:"'sell_rule' VARCHAR(200) DEFAULT ''"`
	Comment     string `xorm:"'comment' VARCHAR(200) DEFAULT ''"`
	SoldOut     int `xorm:"'sold_out' TINYINT(4) notnull DEFAULT '0'"`
	SoldOutDate *time.Time `xorm:"'sold_out_date' DATE DEFAULT ''"`
	Disabled    int `xorm:"'disabled' TINYINT(4) notnull DEFAULT '0'"`
}

func GetAllList() ([]Pharmacy, error) {
	var list = []Pharmacy{}
	engine := storage.GStorage.GetDB().Engine
	se := engine.NewSession()
	err := se.Find(&list)
	//cmd, opt := se.LastSQL()
	//log.Info("sql cmd:", cmd, ", opt:", opt)
	if err != nil {
		return []Pharmacy{}, err
	}
	return list, nil
}

func InsertOrUpdate(maskData *Pharmacy) (int64, error) {
	engine := storage.GStorage.GetDB().Engine
	se := engine.NewSession()
	sql := "INSERT INTO pharmacy (`code`,`name`,tel, addr, adult_count, child_count, upd_time) VALUES (?, ?, ?, ?, ?, ?, ?) " +
		"ON duplicate KEY UPDATE adult_count=?, child_count=?, upd_time=?, sold_out=?, sold_out_date=? "

	res, err := se.Exec(sql, maskData.Code, maskData.Name, maskData.Tel, maskData.Addr, maskData.AdultCount, maskData.ChildCount, maskData.UpdTime,
		maskData.AdultCount, maskData.ChildCount, maskData.UpdTime, maskData.SoldOut, maskData.SoldOutDate)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func UpdateSellRule(code string, sellRule string) error {
	engine := storage.GStorage.GetDB().Engine
	se := engine.NewSession()
	sql := "UPDATE pharmacy SET sell_rule = ? WHERE code = ?"
	_, err := se.Exec(sql, sellRule, code)
	return err
}

func UpdateSoldOut(code string, date time.Time) error {
	engine := storage.GStorage.GetDB().Engine
	se := engine.NewSession()
	sql := "UPDATE pharmacy SET sold_out = ?, sold_out_date = ? WHERE code = ?"
	_, err := se.Exec(sql, 1, date.Format("2006-01-02"), code)
	return err
}

