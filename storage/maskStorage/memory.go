package maskStorage

import (
	"../../entity/pharmacy"
	"../mysql"
	"github.com/poolqa/log"
	"strconv"
	"strings"
	"time"
)

var GMaskStorage *Storage

type Storage struct {
	maskMap *MaskMap
	db      *mysql.MysqlStore
}

func NewMaskStorage(db *mysql.MysqlStore) *Storage {
	return &Storage{
		maskMap: NewMap(),
		db:      db,
	}
}

func (st *Storage) InitData() error {
	st.maskMap.Clear()
	rows, err := pharmacy.GetAllList()
	if err != nil {
		log.Error("query data error:", err)
		return err
	}
	now := time.Now()
	for idx := range rows {
		row := rows[idx]
		//log.Debugf("%#v", row)
		if row.SoldOut != 0 {
			row.SoldOut = 1
		}
		if row.SoldOutDate != nil && row.SoldOutDate.YearDay() == now.YearDay() {
			row.SoldOutDate = nil
			row.SoldOut = 0
		}
		if row.Disabled != 0 {
			row.Disabled = 1
		}
		st.maskMap.Add(row.Code, row)
	}
	return nil
}

func (st *Storage) Set(code string, Name string, Addr string, Tel string, adultCount string, childCount string, updTime string) bool {
	if st.maskMap.Has(code) {
		iAdultCount, _ := strconv.ParseInt(adultCount, 10, 64)
		iChildCount, _ := strconv.ParseInt(childCount, 10, 64)
		tmpTime, err := time.ParseInLocation("2006/01/02 15:04:05", strings.TrimSpace(updTime), time.Local)
		if err != nil {
			log.Error("set mask data parse time error:" + err.Error())
			log.Errorf("code:%v, Name:%v, Addr:%v, Tel:%v, adultCount:%v, childCount:%v, updTime:%v",
				code, Name, Tel, Addr, adultCount, childCount, updTime)
			return st.maskMap.Upd(code, iAdultCount, iChildCount, nil)
		} else {
			return st.maskMap.Upd(code, iAdultCount, iChildCount, &tmpTime)
		}
	} else {
		record := pharmacy.Pharmacy{
			Id:        0,
			Code:      code,
			Name:      Name,
			Tel:       Tel,
			Addr:      Addr,
			SellRule:  "",
			Comment:   "",
			SoldOut:   0,
			Latitude:  0,
			Longitude: 0,
		}
		record.AdultCount, _ = strconv.ParseInt(adultCount, 10, 64)
		record.ChildCount, _ = strconv.ParseInt(childCount, 10, 64)
		tmpTime, err := time.ParseInLocation("2006/01/02 15:04:05", strings.TrimSpace(updTime), time.Local)
		if err != nil {
			log.Error("set mask data parse time error:", err.Error())
			log.Errorf("code:%v, Name:%v, Addr:%v, Tel:%v, adultCount:%v, childCount:%v, updTime:%v",
				code, Name, Tel, Addr, adultCount, childCount, updTime)
			record.UpdTime = nil
		} else {
			record.UpdTime = &tmpTime
		}
		log.Infof("new mask record:%+v", record)
		st.maskMap.Add(code, record)
	}
	return true
}

func (st *Storage) UpdSellRule(code string, sellRule string) bool {
	if st.maskMap.Has(code) {
		return st.maskMap.UpdSellRule(code, sellRule)
	} else {
		return false
	}
}

func (st *Storage) UpdSoldOut(code string, soldOutDate *time.Time) bool {
	if st.maskMap.Has(code) {
		return st.maskMap.UpdSoldOut(code, soldOutDate)
	} else {
		return false
	}
}

func (st *Storage) Flush() {
	//log.Debug("flush start")
	now := time.Now()
	//today := now.Format("2006-01-02")
	maskDataList := st.maskMap.ValList()
	//log.Info("maskDataList size:", len(maskDataList))
	for idx := range maskDataList {
		maskData := maskDataList[idx]
		//log.Infof("maskData:%+v", maskData)
		if maskData.SoldOutDate != nil && maskData.SoldOutDate.YearDay() != now.YearDay() {
			maskData.SoldOutDate = nil
			maskData.SoldOut = 0
		}
		_, err := pharmacy.InsertOrUpdate(&maskData)
		if err != nil {
			log.Error("update pharmacy data error:", err)
			return
		}
	}
	//log.Debug("finish flush")
}

func (st *Storage) GetByCode(code string) *pharmacy.Pharmacy {
	return st.maskMap.Get(code)
}

func (st *Storage) GetAllList() []pharmacy.Pharmacy {
	return st.maskMap.ValList()
}
