package main

import (
	"encoding/json"
	"fmt"
	"dehaa.com/core/datasource"
	"os"
	"dehaa.com/core/controller"
	"github.com/joho/godotenv"
	"log"
	"io/ioutil"
	"strings"
	"dehaa.com/core/model"
	"dehaa.com/core/migration"
)

var baseController *controller.BaseController
var twitterApi datasource.TwitterApi
var coinmarketcal datasource.CoinMarketCal
var coins map[string]datasource.Coin

//init function
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
	env := os.Getenv("env")

	file, _ := ioutil.ReadFile("config/" + env + "/twitter.json")
	twitterApi = datasource.TwitterApi{}
	_ = json.Unmarshal(file, &twitterApi.Config)
	twitterApi.CreateClient()

	file, _ = ioutil.ReadFile("config/" + env + "/coinmarketcal.json")
	coinmarketcal = datasource.CoinMarketCal{}
	_ = json.Unmarshal(file, &coinmarketcal.Config)
}

func main() {

	baseController := new(controller.BaseController)
	baseController.DBConnect()

	CoinMarketCapControl()

	defer baseController.Destruct()
}

func CoinMarketCapControl() {
	cmcController := controller.NewCoinMarketCapController()
	coins, _ := cmcController.GetAllCoins()

	for _, coin := range coins {

		if exists, coinModel := cmcController.IsCoinExists(coin.Symbol); exists {
			migration.CoinCheckTable(coinModel, cmcController.BaseController.DB)
		} else {
			added, coinModel := cmcController.AddCoin(coin)
			if added {
				migration.CoinCheckTable(coinModel, cmcController.BaseController.DB)
			} else {
				log.Println(coinModel.Name + " could not be added")
			}
		}
	}

	cmcController.BaseController.Destruct()
}

//WriteTweets
func WriteTweets(coin *model.Coin) {
	result := twitterApi.Search("$"+coin.Symbol, "en", "links")
	for _, sonuc := range result.Statuses {
		fmt.Println("fulltext", sonuc.FullText)
	}
}

//CoinMarketCalEventsControl
func CoinMarketCalEventsControl(coin *model.Coin) {
	a := datasource.CoinMarketCal{}

	var events = controller.CoinMarketCalEvents{}

	coinMarketCalEvents := a.GetEvents()
	err := json.Unmarshal(coinMarketCalEvents, &events)

	if err != nil {
		var errJson map[string]interface{}
		err = json.Unmarshal(coinMarketCalEvents, &errJson)
		if errJson["error"] != nil && strings.EqualFold(errJson["error_description"].(string), "The access token provided is invalid.") {
			fmt.Println("s")
		}

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for _, event := range events {

		category := model.CoinMarketCalCategory{Category: event.Categories[0].Name}

		//db.Create(&category)
		//db.First(&category, "category_id = ?", category.CategoryId)
		fmt.Print(category)

		//db.Model(&event).Related(&category, "Category")
		event.CategoryId = int64(category.CategoryId)
		//db.Create(&event)
		//db.Create(&event)
	}
}

func MarketControl(symbol string) {

}

func QuoteControl(symbol string) {

}
