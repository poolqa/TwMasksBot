package crawler

import (
	"../../storage/maskStorage"
	"crypto/tls"
	"github.com/pkg/errors"
	"github.com/poolqa/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type NHICrawler struct {
	url     string
	storage *maskStorage.Storage
}

func NewNHICrawler(url string, storage *maskStorage.Storage) *NHICrawler {
	return &NHICrawler{url: url, storage: storage}
}

func (mCrawler *NHICrawler) InitData() error {
	return mCrawler.storage.InitData()
}

func (mCrawler *NHICrawler) Run() error {
	csvData, err := mCrawler.QueryMasksData()
	if err != nil {
		return err
	}
	if len(csvData) == 0 {
		return errors.New("not got any csv data")
	}
	err = mCrawler.Parse(csvData)
	if err != nil {
		return err
	}
	return nil
}

func (mCrawler *NHICrawler) Print() {
	maskData := mCrawler.storage.GetAllList()
	log.Debugf("maskData:%+v", maskData)
}
func (mCrawler *NHICrawler) Parse(csvData string) error {
	// 醫事機構代碼,醫事機構名稱,醫事機構地址,醫事機構電話,成人口罩總剩餘數,兒童口罩剩餘數,來源資料時間
	// 0           1          2           3          4              5            6
	//log.Debug("csv:", csvData)
	csvLines := strings.Split(csvData, "\n")
	for idx := 1; idx < len(csvLines); idx++ {
		line := csvLines[idx]
		row := strings.Split(line, ",")
		if len(row) < 7 {
			continue
		}
		//log.Infof("code:%v, name:%v, adult:%v, child:%v, upd_time:%v", row[0], row[1], row[4], row[5], row[6])
		if !mCrawler.storage.Set(row[0], row[1], row[2], row[3], row[4], row[5], row[6]) {
			log.Error("set mask data error:", row)
		}
	}
	go mCrawler.storage.Flush()
	return nil
}

func (mCrawler *NHICrawler) QueryMasksData() (string, error) {
	// link := "https://data.nhi.gov.tw/resource/mask/maskdata.csv"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Get(mCrawler.url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	return string(content), nil
}
