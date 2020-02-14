package telegram

import (
	"../../../config"
	"../../../storage/maskStorage"
	"../command"
	"../utils"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/poolqa/log"
	"strconv"
	"time"
)

type Bot struct {
	token   string
	debug   bool
	storage *maskStorage.Storage
	botApi  *tgbotapi.BotAPI
}

func NewTelegramBot(conf *config.TgConfig, storage *maskStorage.Storage) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Error("new telegram bot api error:", err)
		return nil, err
	}
	return &Bot{token: conf.Token, debug: conf.Debug, storage: storage, botApi: botApi}, nil
}

func (mBot *Bot) Handler() {
	mBot.botApi.Debug = mBot.debug

	log.Infof("Authorized on account %s", mBot.botApi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := mBot.botApi.GetUpdatesChan(u)
	if err != nil {
		log.Error("GetUpdatesChan error:", err)
		return
	}

	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			go mBot.ProcessMessage(update)
		} else if update.CallbackQuery != nil {
			go mBot.ProcessCallbackQuery(update)
		}
	}
}

func (mBot *Bot) ProcessMessage(update tgbotapi.Update) error {
	log.Warnf("MessageID:%v, From:%v, Date:%v, Text:%+v, Location:%+v, Entities:%+v", update.Message.MessageID, update.Message.From, update.Message.Date, update.Message.Text, update.Message.Location, update.Message.Entities)

	strUserId := strconv.Itoa(update.Message.From.ID)
	if update.Message.Location != nil {
		filterData := commParser.GetFilterData(strUserId)
		maskDataList := utils.CalcLocation(mBot.storage, 5, filterData, update.Message.Location.Latitude, update.Message.Location.Longitude)
		if len(maskDataList) == 0 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "無符合條件目標，請重新查詢")
			mBot.botApi.Send(msg)
		} else {
			for idx := range maskDataList {
				mBot.SendMaskData(update.Message.Chat.ID, maskDataList[idx])
			}
		}
	} else if update.Message.Entities != nil {
		var msg tgbotapi.Chattable
		reqMsgEntry := (*update.Message.Entities)[0]
		command := update.Message.Text[reqMsgEntry.Offset:reqMsgEntry.Length]
		commType, rText := commParser.ProcessMessage(command, strUserId)
		switch commType {
		case commParser.CommandTypeHelp:
			msg = tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "./start.jpg")
		case commParser.CommandTypeQuestion:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, rText)
		case commParser.CommandTypeFilter:
			srcMsg := tgbotapi.NewMessage(update.Message.Chat.ID, rText)
			srcMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(utils.FilterAdult, utils.FilterAdult),
					tgbotapi.NewInlineKeyboardButtonData(utils.FilterChild, utils.FilterChild),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(utils.FilterAdultAndChild, utils.FilterAdultAndChild),
					tgbotapi.NewInlineKeyboardButtonData(utils.FilterZero, utils.FilterZero),
				),
			)
			msg = srcMsg
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "目前尚在開發階段，不定時抽風請見諒\n發送定位後，即可得到最近的五間藥局座標")
		}
		mBot.botApi.Send(msg)
	} else if update.Message.Text != "" {
		var msg tgbotapi.Chattable
		commType, rText := commParser.ProcessMessage(update.Message.Text, strUserId)
		switch commType {
		case commParser.CommandTypeHelp:
			msg = tgbotapi.NewPhotoUpload(update.Message.Chat.ID, "./start.jpg")
		case commParser.CommandTypeQuestion, commParser.CommandTypeAnswer:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, rText)
		default:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "目前尚在開發階段，不定時抽風請見諒\n發送定位後，即可得到最近的五間藥局座標")
		}
		mBot.botApi.Send(msg)
	} else {

	}
	return nil
}

func (mBot *Bot) ProcessCallbackQuery(update tgbotapi.Update) error {
	message := update.CallbackQuery.Message
	chatId := update.CallbackQuery.Message.Chat.ID
	log.Warnf("CallbackQuery MessageID:%v, From:%v, Date:%v, Text:%+v, CallbackQuery.Data:%+v", message.MessageID, update.CallbackQuery.From, message.Date, message.Text, update.CallbackQuery.Data)
	strUserId := strconv.Itoa(update.CallbackQuery.From.ID)
	rText := commParser.ProcessFilterAnswer(update.CallbackQuery.Data, strUserId)

	mBot.botApi.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
	mBot.botApi.Send(tgbotapi.NewDeleteMessage(chatId, message.MessageID))
	msg := tgbotapi.NewMessage(chatId, rText)
	mBot.botApi.Send(msg)
	return nil
}

func (mBot *Bot) SendMaskData(chatId int64, data utils.PharmacyDistance) error {
	text := ""
	now := time.Now()
	//today := now.Format("2006-01-02")
	if data.MaskData.SoldOut != 0 && data.MaskData.SoldOutDate != nil && data.MaskData.SoldOutDate.YearDay() == now.YearDay() {
		text = fmt.Sprintf("%v %v (網友回報已完售)\n", data.MaskData.Name, data.MaskData.Code)
	} else {
		text = fmt.Sprintf("%v %v\n", data.MaskData.Name, data.MaskData.Code)
	}
	text += fmt.Sprintf("距離:%.0f公尺\n", data.Distance)
	text += fmt.Sprintf("TEL:%v\n地址:%v\n成人:%v, 兒童:%v\n", data.MaskData.Tel, data.MaskData.Addr, data.MaskData.AdultCount, data.MaskData.ChildCount)
	if data.MaskData.UpdTime != nil {
		text += fmt.Sprintf("更新時間:%v\n", data.MaskData.UpdTime.Format("2006/01/02 15:04:05"))
	} else {
		text += fmt.Sprintf("更新時間:%v\n", "(無)")
	}
	text += fmt.Sprintf("銷售規則:%v\n", data.MaskData.SellRule)
	textMsg := tgbotapi.NewMessage(chatId, text)
	locationMsg := tgbotapi.NewLocation(chatId, data.MaskData.Latitude, data.MaskData.Longitude)
	mBot.botApi.Send(textMsg)
	mBot.botApi.Send(locationMsg)
	return nil
}
