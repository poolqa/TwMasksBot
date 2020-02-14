package soldOutReport

import (
	"../../storage"
	"github.com/poolqa/log"
	"github.com/xormplus/xorm"
	"time"
)

type SoldOutReport struct {
	Id           int64      `xorm:"'id' BIGINT(20) notnull pk autoincr"`
	Code         string     `xorm:"'code' VARCHAR(20) DEFAULT ''"`
	ReportDate   *time.Time `xorm:"'report_date' DATE notnull "`
	ReportUserId string     `xorm:"'report_user_id' VARCHAR(40) notnull"`
	CreateTime   *time.Time `xorm:"'create_time' timestamp NULL DEFAULT NULL "`
}

func updateSoldOut(se *xorm.Session, userId string, code string, date time.Time) error {
	sql := "INSERT INTO sold_out_report (`code`, report_date, report_user_id) VALUES (?, ?, ?) " +
		"ON duplicate KEY UPDATE id=id "
	_, err := se.Exec(sql, code, date.Format("2006-01-02"), userId)
	return err
}

func UpdateAndReturnCnt(userId string, code string, date time.Time) (int64, error) {
	engine := storage.GStorage.GetDB().Engine
	se := engine.NewSession()
	err := updateSoldOut(se, userId, code, date)
	if err != nil {
		log.Errorf("insert sold_out_report error:%v, pharmacy code:%+v, userId:%v", err, code, userId)
		return 0, err
	}
	qSql := "SELECT COUNT(DISTINCT report_user_id) as cnt FROM sold_out_report WHERE code = ? AND report_date = ? "
	var cnt int64
	_, err = se.SQL(qSql, code, date.Format("2006-01-02")).Get(&cnt)
	if err != nil {
		log.Errorf("query sold_out_report error:%v, pharmacy code:%+v, report_date:%v", err, code, date.Format("2006-01-02"))
	}
	return cnt, nil
}