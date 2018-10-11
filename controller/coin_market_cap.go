package controller

import (
	"dehaa.com/core/datasource"
	"fmt"
	"time"
	"dehaa.com/core/model"
	"reflect"
	"strings"
	"dehaa.com/core/migration"
)

type CoinMarketCapController struct {
	BaseController *BaseController
}

var coinMarketCapController *CoinMarketCapController

//NewCoinMarketCapController Constructor
func NewCoinMarketCapController() *CoinMarketCapController {

	if coinMarketCapController == nil {
		coinMarketCapController = new(CoinMarketCapController)
	}

	return coinMarketCapController
}

//IsCoinExists
func (c *CoinMarketCapController) IsCoinExists(symbol string) (bool, *model.Coin) {
	var coinModel model.Coin
	c.BaseController.DB.Where("symbol = ?", symbol).First(&coinModel)
	if !reflect.DeepEqual(coinModel, model.Coin{}) {
		return true, &coinModel
	}
	return false, &model.Coin{}
}

//addCoin from coinmarketcap.com if not added
func (c *CoinMarketCapController) AddCoin(coinData datasource.Coin) (bool, *model.Coin) {

	coinModel := new(model.Coin)
	coinModel.Name = coinData.ID
	coinModel.Rank = coinData.Rank
	coinModel.Symbol = coinData.Symbol
	coinModel.AvailableSupply = coinData.AvailableSupply
	coinModel.TotalSupply = coinData.TotalSupply
	coinModel.LastUpdated = c.BaseController.timestampToTime(coinData.LastUpdated).Format("2006-01-02 15:04:05")
	c.BaseController.DB.Create(&coinModel)

	return true, coinModel
}

//updateCoin
func (c *CoinMarketCapController) updateCoin(coinData datasource.Coin, coinModel *model.Coin) (bool, *model.Coin) {

	if coinModel.Rank != coinData.Rank ||
		coinModel.AvailableSupply != coinData.AvailableSupply ||
		coinModel.TotalSupply != coinData.TotalSupply {
		var coin model.Coin
		c.BaseController.DB.Model(&coin).Where("symbol = ?", coinData.Symbol).Updates(map[string]interface{}{
			"rank":             coinData.Rank,
			"available_supply": coinData.AvailableSupply,
			"total_supply":     coinData.TotalSupply,
			"last_updated":     c.BaseController.timestampToTime(coinData.LastUpdated).Format("2006-01-02 15:04:05"),
		})

		var returnCoin *model.Coin
		returnCoin = c.GetCoin(coinData.Symbol)
		return true, returnCoin
	}

	return true, coinModel
}

//GetCoin from database
func (c *CoinMarketCapController) GetCoin(symbol string) *model.Coin {
	var coins []*model.Coin
	c.BaseController.DB.Model(&model.Coin{}).Where("symbol=?", symbol).First(&coins)
	return coins[0]
}

//GlobalData
func (c *CoinMarketCapController) GlobalData(timeInterval ...int64) (*model.GlobalMarketData, *datasource.MarketGraph) {
	switch len(timeInterval) {
	case 0:
		data, _ := datasource.GetGlobalMarketData()

		globalMarketData := new(model.GlobalMarketData)
		globalMarketData.Active_assets = data.ActiveAssets
		globalMarketData.Active_currencies = data.ActiveCurrencies
		globalMarketData.Active_markets = data.ActiveMarkets
		globalMarketData.BitcoinPercentageOfMarketCap = data.BitcoinPercentageOfMarketCap
		globalMarketData.Total24HVolumeUSD = data.Total24HVolumeUSD
		globalMarketData.TotalMarketCapUSD = data.TotalMarketCapUSD

		c.BaseController.DB.Create(&globalMarketData)
		return globalMarketData, nil
	case 2:
		globalMarketGraphData, _ := datasource.GetGlobalMarketGraphData(timeInterval[0], timeInterval[1])
		for _, globalMarketData := range globalMarketGraphData.MarketCapByAvailableSupply {
			fmt.Println(globalMarketData)
			fmt.Println(globalMarketGraphData)
		}
		return nil, &globalMarketGraphData
	}

	return nil, nil

}

func (c *CoinMarketCapController) GlobalDataCaller() {
	then := time.Date(
		2009, 01, 01, 0, 0, 0, 0, time.UTC)
	c.GlobalData(then.Unix(), time.Now().Unix())
}

//isMarketExists
func (c *CoinMarketCapController) IsMarketExists(marketName string) (bool, *model.Market) {
	var market model.Market
	c.BaseController.DB.Where("market_code = ?", strings.ToLower(marketName)).First(&market)
	if market != (model.Market{}) {
		return true, &market
	}
	return false, &model.Market{}
}

//addMarket
func (c *CoinMarketCapController) addMarket(marketName string) (bool, *model.Market) {
	var market model.Market
	market.MarketName = marketName
	market.MarketCode = strings.ToLower(marketName)
	c.BaseController.DB.Create(&market)
	return true, &market
}

//addCoinMarketData
func (c *CoinMarketCapController) addCoinMarketData(coinMarketData datasource.Market) (bool) {
	/*
		pair := strings.Split(coinMarketData.Pair, "/")

		coinName := pair[0]
		base := pair[1]

		var coin model.Coin
		db.Where("symbol = ?", coinName).First(&coin)

		var market model.Market
		db.Where("market_code = ?", strings.ToLower(coinMarketData.Exchange)).First(&market)

		if !reflect.DeepEqual(coin, model.Coin{}) && !reflect.DeepEqual(market, model.Market{}) {

			var markets []model.CoinMarket
			var coinMarket model.CoinMarket

			db.Where("coin_id = ? and market_id = ?", coin.CoinId, market.MarketId).Find(&markets)
			if len(markets) == 0 {
				coinMarket.CoinId = coin.CoinId
				coinMarket.MarketId = market.MarketId
				db.Create(&coinMarket)
			}

			marketData := new(model.CoinMarketData)
			marketData.CoinId = coin.CoinId
			marketData.MarketId = market.MarketId
			marketData.Rank = coinMarketData.Rank
			marketData.Base = base
			marketData.Price_USD = coinMarketData.Price
			marketData.Volume_USD = coinMarketData.VolumeUSD
			marketData.VolumePercent = coinMarketData.VolumePercent

			db.Create(&marketData)
			return true
		}
	*/
	return false

}

//GetAllCoins
func (c *CoinMarketCapController) GetAllCoins() (map[string]datasource.Coin, error) {
	datas, err := datasource.GetAllCoinData(0)

	return datas, err
}

func (c *CoinMarketCapController) DataTableControl(datas map[string]datasource.Coin) {

	var coinModels []*model.Coin
	var newCoinModel *model.Coin

	for _, coinData := range datas {

		if isExists, coinModel := c.IsCoinExists(coinData.Symbol); isExists {
			_, newCoinModel = c.updateCoin(coinData, coinModel)
		} else {
			_, newCoinModel = c.AddCoin(coinData)
			newCoinModel = c.GetCoin(coinData.Symbol)
		}

		migration.CoinCheckTable(newCoinModel, c.BaseController.DB)
		coinModels = append(coinModels, newCoinModel)

		//b, _ := datasource.GetAltcoinMarketGraphData(time.Now().AddDate(-3, 0, 0).Unix(), time.Now().Unix())
		//c, _ := datasource.GetCoinData(coin)
	}
}

func (c *CoinMarketCapController) AddCoinMarket(coin *model.Coin) {

	var markets []datasource.Market
	markets, err := datasource.GetCoinMarkets(coin.Symbol)
	if err != nil {
		fmt.Println(err.Error())
	}

	for _, market := range markets {
		if isExists, _ := c.IsMarketExists(market.Exchange); !isExists {
			c.addMarket(market.Exchange)
		}
		c.addCoinMarketData(market)
	}
}

//GlobalMarketData
func (c *CoinMarketCapController) GlobalMarketData() {

}

func MarketList() {
	exchanges := []string{
		"_1broker",
		"_1btcxe",
		"acx",
		"allcoin",
		"anxpro",
		"bibox",
		"binance",
		"bit2c",
		"bitbank",
		"bitbay",
		"bitfinex",
		"bitfinex2",
		"bitflyer",
		"bithumb",
		"bitkk",
		"bitlish",
		"bitmarket",
		"bitmex",
		"bitso",
		"bitstamp",
		"bitstamp1",
		"bittrex",
		"bitz",
		"bl3p",
		"bleutrade",
		"braziliex",
		"btcbox",
		"btcchina",
		"btcexchange",
		"btcmarkets",
		"btctradeim",
		"btctradeua",
		"btcturk",
		"btcx",
		"bxinth",
		"ccex",
		"cex",
		"chbtc",
		"chilebit",
		"cobinhood",
		"coinbase",
		"coincheck",
		"coinegg",
		"coinex",
		"coinexchange",
		"coinfloor",
		"coingi",
		"coinmarketcap",
		"coinmate",
		"coinnest",
		"coinone",
		"coinsecure",
		"coinspot",
		"cointiger",
		"coolcoin",
		"cryptopia",
		"dsx",
		"ethfinex",
		"exmo",
		"exx",
		"flowbtc",
		"foxbit",
		"fybse",
		"fybsg",
		"gatecoin",
		"gateio",
		"gdax",
		"gemini",
		"getbtc",
		"hadax",
		"hitbtc",
		"hitbtc2",
		"huobi",
		"huobicny",
		"huobipro",
		"ice3x",
		"independentreserve",
		"indodax",
		"itbit",
		"jubi",
		"kraken",
		"kucoin",
		"kuna",
		"lakebtc",
		"lbank",
		"liqui",
		"livecoin",
		"luno",
		"lykke",
		"mercado",
		"mixcoins",
		"negociecoins",
		"nova",
		"okcoincny",
		"okcoinusd",
		"okex",
		"paymium",
		"poloniex",
		"qryptos",
		"quadrigacx",
		"quoinex",
		"southxchange",
		"surbitcoin",
		"therock",
		"tidebit",
		"tidex",
		"urdubit",
		"vaultoro",
		"vbtc",
		"virwox",
		"wex",
		"xbtce",
		"yobit",
		"yunbi",
		"zaif",
		"zb",
	}

	fmt.Println(exchanges)
}
