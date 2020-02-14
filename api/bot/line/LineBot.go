package line

import (
	"../../../config"
	"../../../storage/maskStorage"
	"../command"
	"../utils"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/poolqa/log"
	"net/http"
	"time"
)

type Bot struct {
	port          string
	channelSecret string
	token         string
	debug         bool
	storage       *maskStorage.Storage
	botApi        *linebot.Client
}

func NewLineBot(conf *config.LineConfig, storage *maskStorage.Storage) (*Bot, error) {
	botApi, err := linebot.New(conf.ChannelSecret, conf.Token)
	if err != nil {
		log.Error("new telegram bot api error:", err)
		return nil, err
	}
	return &Bot{port: conf.Port,
		channelSecret: conf.ChannelSecret,
		token: conf.Token,
		debug: conf.Debug,
		storage: storage, botApi: botApi}, nil
}

func (mBot *Bot) Handler() {
	log.Infof("Authorized on line account %+v", mBot.botApi)
	http.HandleFunc("/line-web-hook", mBot.webHookHandler)
	port := mBot.port
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func (mBot *Bot) webHookHandler(w http.ResponseWriter, r *http.Request) {
	events, err := mBot.botApi.ParseRequest(r)
	if err != nil {
		log.Error("ParseRequest error:", err)
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		log.Warnf("event From:%v, Date:%v, message:%#v", event.Source.UserID, event.Timestamp, event.Message)
		strUserId := event.Source.UserID
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var msg linebot.SendingMessage
				commType, rText := commParser.ProcessMessage(message.Text, strUserId)
				switch commType {
				//case commParser.CommandTypeHelp:
				//	msg = tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "./start.jpg")
				case commParser.CommandTypeQuestion, commParser.CommandTypeAnswer:
					msg = linebot.NewTextMessage(rText)
				case commParser.CommandTypeFilter:
					template := linebot.NewButtonsTemplate(
						"", rText, rText,
						linebot.NewMessageAction(utils.FilterAdult, utils.FilterAdult),
						linebot.NewMessageAction(utils.FilterChild, utils.FilterChild),
						linebot.NewMessageAction(utils.FilterAdultAndChild, utils.FilterAdultAndChild),
						linebot.NewMessageAction(utils.FilterZero, utils.FilterZero),
					)
					msg = linebot.NewTemplateMessage(rText, template)
				default:
					msg = linebot.NewTextMessage("目前尚在開發階段，不定時抽風請見諒\n發送定位後，即可得到最近的五間藥局座標")
				}
				if _, err = mBot.botApi.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
					log.Error(err)
				}
			case *linebot.LocationMessage:
				filterData := commParser.GetFilterData(strUserId)
				maskDataList := utils.CalcLocation(mBot.storage, 5, filterData, message.Latitude, message.Longitude)
				if len(maskDataList) == 0 {
					msg := linebot.NewTextMessage("無符合條件目標，請重新查詢")
					_, err = mBot.botApi.ReplyMessage(event.ReplyToken, msg).Do()
					if err != nil {
						log.Error(err)
					}
				} else {
					msgList := []linebot.SendingMessage{}
					for idx := range maskDataList {
						msgList = append(msgList, mBot.CreateMaskDataMessage(event.ReplyToken, maskDataList[idx])...)
					}
					_, err = mBot.botApi.ReplyMessage(event.ReplyToken, msgList...).Do()
					if err != nil {
						log.Error(err)
					}
				}

			}
		}
	}
}

func (mBot *Bot) CreateMaskDataMessage(replyToken string, data utils.PharmacyDistance) []linebot.SendingMessage {
	text := ""
	now := time.Now()
	//today := now.Format("2006-01-02")
	if data.MaskData.SoldOut != 0 && data.MaskData.SoldOutDate != nil && data.MaskData.SoldOutDate.YearDay() == now.YearDay() {
		text = fmt.Sprintf("%v %v (網友回報已完售)\n", data.MaskData.Name, data.MaskData.Code)
	} else {
		text = fmt.Sprintf("%v %v\n", data.MaskData.Name, data.MaskData.Code)
	}
	text += fmt.Sprintf("距離:%.0f公尺\n", data.Distance)
	text += fmt.Sprintf("TEL:%v\n成人:%v, 兒童:%v\n", data.MaskData.Tel, data.MaskData.AdultCount, data.MaskData.ChildCount)
	if data.MaskData.UpdTime != nil {
		text += fmt.Sprintf("更新時間:%v\n", data.MaskData.UpdTime.Format("2006/01/02 15:04:05"))
	} else {
		text += fmt.Sprintf("更新時間:%v\n", "(無)")
	}
	text += fmt.Sprintf("銷售規則:%v\n", data.MaskData.SellRule)
	//textMsg := linebot.NewTextMessage(text)
	locationMsg := linebot.NewLocationMessage(text, data.MaskData.Addr, data.MaskData.Latitude, data.MaskData.Longitude)

	return []linebot.SendingMessage{locationMsg}
}