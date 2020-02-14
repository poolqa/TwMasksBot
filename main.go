package main

import (
	"./api/bot/line"
	"./api/bot/telegram"
	"./config"
	"./api/crawler"
	"./storage"
	"./storage/maskStorage"
	"flag"
	"github.com/jasonlvhit/gocron"
	"github.com/poolqa/log"
	"os"
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

	nhiMaskDataCrawler := crawler.NewNHICrawler("https://data.nhi.gov.tw/resource/mask/maskdata.csv", maskStorage.GMaskStorage)
	cron := gocron.NewScheduler()
	cron.Every(60).Seconds().Do(nhiMaskDataCrawler.Run)
	cron.Start()
	cron.RunAll()

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

//func loadFileToDb() {
//	fileName := "d:/drugstore-gps.csv"
//	f, err := os.Open(fileName)
//	if err != nil {
//	}
//	defer f.Close()
//
//	lineArr := []string{}
//	br := bufio.NewReader(f)
//	for {
//		line, _, err := br.ReadLine()
//		if err == io.EOF {
//			break
//		}
//		lineArr = append(lineArr, string(line))
//	}
//
//	sql := "INSERT INTO pharmacy (`code`,`name`,tel, addr, latitude, longitude) VALUES (?, ?, ?, ?, ?, ?) " +
//		"ON duplicate KEY UPDATE tel=?, addr=? "
//	dbSe := storage.GStorage.GetDB().NewSession()
//
//	for idx := 1; idx < len(lineArr); idx++ {
//		line := lineArr[idx]
//		row := strings.Split(line, ",")
//		log.Infof("idx:%v, %+v", idx, row)
//		_, err := dbSe.Execute(sql, row[0], row[1], row[2], row[3], row[5], row[6],
//			row[2], row[3])
//		if err != nil {
//			log.Error("update pharmacy data error:", err)
//			return
//		}
//	}
//}
