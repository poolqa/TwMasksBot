package commParser

import (
	"../../../entity/pharmacy"
	"../../../entity/soldOutReport"
	"../../../storage"
	"../../../storage/maskStorage"
	"../utils"
	"fmt"
	"github.com/poolqa/log"
	"strings"
	"time"
)

const (
	CommandStart    = "/start"
	CommandHelp     = "/help"
	CommandSellRule = "/sell_rule"
	CommandSoldOut  = "/sold_out"
	CommandFilter  = "/filter"
)

type CommandType int
const (
	CommandTypeNone CommandType = iota
	CommandTypeHelp
	CommandTypeQuestion
	CommandTypeAnswer
	CommandTypeFilter

)

func ProcessMessage(text string, userId string) (CommandType, string) {
	conn := storage.GStorage.GetRedisConn()
	switch text {
	case CommandStart, CommandHelp:
		return CommandTypeHelp, ""
	case CommandSellRule:
		resp := "請輸入藥局代碼 欲補充規則(盡量簡短)\n例:5901012203 早上9點排隊領取號碼牌"
		err := conn.Set(userId, text, 180)
		if err != nil {
			resp = "500 伺服器發生錯誤"
		}
		return CommandTypeQuestion, resp
	case CommandSoldOut:
		resp := "如要回報藥局已完售，請輸入藥局代碼，有兩人以上回報才會進行標示\n例:5901012203"
		err := conn.Set(userId, text, 180)
		if err != nil {
			resp = "500 伺服器發生錯誤"
		}
		return CommandTypeQuestion, resp
	case CommandFilter:
		resp := "請選擇想要設定的過濾條件"
		err := conn.Set(userId, text, 180)
		if err != nil {
			resp = "500 伺服器發生錯誤"
		}
		return CommandTypeFilter, resp
	default:
		value, err := conn.Pop(userId)
		if err != nil {
			return CommandTypeAnswer, "500 伺服器發生錯誤"

		} else if value != nil {
			// question
			command := value.(string)
			log.Infof("question:%v", command)
			switch command {
			case CommandSellRule:
				resp := updateSellRule(text, userId)
				return CommandTypeAnswer, resp
			case CommandSoldOut:
				resp := updateSoldOut(text, userId)
				return CommandTypeAnswer, resp
			case CommandFilter:
				resp := ProcessFilterAnswer(text, userId)
				return CommandTypeAnswer, resp
			default:
				return CommandTypeAnswer, "500 伺服器發生錯誤"
			}
		}
		return CommandTypeNone, ""
	}
	return CommandTypeNone, ""
}

func GetFilterData(userId string) string {
	conn := storage.GStorage.GetRedisConn()
	val, err := conn.Pop(userId+CommandFilter)
	if err != nil {
		return ""
	} else if val != nil {
		// question
		return val.(string)
	}
	return ""
}

func SetFilterData(userId string, filter string) error {
	conn := storage.GStorage.GetRedisConn()
	err := conn.Set(userId+CommandFilter, filter, 180)
	return err
}

func updateSellRule(answer string, userId string) string {
	ans := strings.Split(answer, " ")
	if len(ans) < 2 {
		return "藥局代號及規則請用半形空格隔開"
	}
	maskData := maskStorage.GMaskStorage.GetByCode(ans[0])
	if maskData == nil {
		return "找不到藥局代號，請確認後在重新輸入命令"
	}
	text := ""
	if len(ans) > 50 {
		text = ans[1][0:50]
	} else {
		text = ans[1]
	}
	ok := maskStorage.GMaskStorage.UpdSellRule(maskData.Code, text)
	if !ok {
		log.Errorf("update memory data error, maskData:%+v, rule:%v", maskData, text)
		return "更新失敗"
	}
	err := pharmacy.UpdateSellRule(maskData.Code, text)
	if err != nil {
		log.Errorf("update db data error, maskData:%+v, error:%v", maskData, err)
		return "更新資料庫失敗"
	}
	return "回報成功"
}

func updateSoldOut(answer string, userId string) string {
	now := time.Now()
	today := now.Format("2006-01-02")
	maskData := maskStorage.GMaskStorage.GetByCode(answer)
	if maskData == nil {
		return "找不到藥局代號，請確認後在重新輸入命令"
	}

	cnt, err := soldOutReport.UpdateAndReturnCnt(userId, maskData.Code, now)
	if err != nil {
		log.Errorf("insert sold_out_report error:%v, maskData:%+v, userId:%v", err, maskData, userId)
		return "更新資料庫失敗"
	}
	if cnt >= 2 {
		ok := maskStorage.GMaskStorage.UpdSoldOut(maskData.Code, &now)
		if !ok {
			log.Errorf("update memory data error, maskData:%+v, today:%v", maskData, today)
			return "更新失敗"
		}
		err = pharmacy.UpdateSoldOut(maskData.Code, now)
		if err != nil {
			log.Errorf("update db data error, maskData:%+v, error:%v", maskData, err)
			return "更新資料庫失敗"
		}
	}
	return "回報成功"
}

func ProcessFilterAnswer(text, strUserId string) string {
	rText := ""
	switch text {
	case utils.FilterAdult:
		rText = fmt.Sprintf("你已經選擇\"%v\",請於三分鐘內發送定位", text)
		SetFilterData(strUserId, text)
	case utils.FilterChild:
		rText = fmt.Sprintf("你已經選擇\"%v\",請於三分鐘內發送定位", text)
		SetFilterData(strUserId, text)
	case utils.FilterAdultAndChild:
		rText = fmt.Sprintf("你已經選擇\"%v\",請於三分鐘內發送定位", text)
		SetFilterData(strUserId, text)
	case utils.FilterZero:
		rText = fmt.Sprintf("你已經選擇\"%v\",請於三分鐘內發送定位", text)
		SetFilterData(strUserId, text)
	default:
		rText = "選擇錯誤，請重新執行過濾命令"
	}
	return rText
}