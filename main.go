package main

import (
	"./api/bot/line"
	"./api/bot/telegram"
	"./api/crawler"
	"./config"
	"./storage"
	"./storage/maskStorage"
	"encoding/csv"
	"flag"
	"github.com/jasonlvhit/gocron"
	"github.com/poolqa/log"
	"io/ioutil"
	"os"
	"strings"
)

var (
	//支持命令行输入格式为-config=name, 默认为config.ini
	//配置文件一般获取到都是类型
	configFile = flag.String("config", "./config.ini", "General configuration file")
)

func main() {
	log.InitByConfigFile("./log.conf")

	flag.Parse()

	var err error
	config.GConfig, err = config.ParseConfig(*configFile) //读取配置文件，并返回其Config
	if err != nil {
		log.Error("找不到配置文件！信息：" + err.Error())
		os.Exit(0)
	}

	storage.GStorage = storage.NewStorage(config.GConfig)

	maskStorage.GMaskStorage = maskStorage.NewMaskStorage(storage.GStorage.GetDB())
	err = maskStorage.GMaskStorage.InitData()
	if err != nil {
		log.Error("nhiMaskDataCrawler init data error:", err)
		os.Exit(0)
	}


	if true {
		nhiMaskDataCrawler := crawler.NewNHICrawler("https://data.nhi.gov.tw/resource/mask/maskdata.csv", maskStorage.GMaskStorage)
		cron := gocron.NewScheduler()
		cron.Every(60).Seconds().Do(nhiMaskDataCrawler.Run)
		cron.Start()
		cron.RunAll()
	} else {
		loadFileToDb()
	}

	if config.GConfig.Telegram.Enable {
		tgBot, _ := telegram.NewTelegramBot(config.GConfig.Telegram, maskStorage.GMaskStorage)
		go tgBot.Handler()
	}
	if config.GConfig.Line.Enable {
		tgBot, _ := line.NewLineBot(config.GConfig.Line, maskStorage.GMaskStorage)
		go tgBot.Handler()
	}
	select {}

}

func loadFileToDb() {
	fileName := "d:/gps.csv"
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}


	br := csv.NewReader(strings.NewReader(string(f)))

	csvData, _ := br.ReadAll()
	csvLineSize := len(csvData)
	se := storage.GStorage.GetDB().Engine.NewSession()
	sql := "UPDATE pharmacy SET latitude = ?, longitude = ? WHERE `code` =? AND latitude = 0 AND longitude = 0 "
	for idx := 1; idx < csvLineSize; idx++ {
		row := csvData[idx]
		log.Infof("idx:%v, %+v, %v, %v", idx, row[13], row[12], row[0])
		_, err := se.Exec(sql, row[13], row[12], row[0])
		if err != nil {
			log.Error("update pharmacy data error:", err)
			return
		}
	}
}
