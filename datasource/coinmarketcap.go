package datasource

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	soup "dehaa.com/core/lib"
	"reflect"
)

// Interface interface
type Interface interface {
	GetGlobalMarketData() (GlobalMarketData, error)
	GetGlobalMarketGraphData(start int64, end int64) (MarketGraph, error)
	GetAltcoinMarketGraphData(start int64, end int64) (MarketGraph, error)
	GetCoinData(coin string) (Coin, error)
	GetAllCoinData(limit int) (map[string]Coin, error)
	GetCoinGraphData(coin string, start int64, end int64) (CoinGraph, error)
	GetCoinPriceUSD(coin string) (float64, error)
	GetCoinMarkets(coin string) ([]Market, error)
}

// Coin struct
type Coin struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             int     `json:"rank,string"`
	PriceUSD         float64 `json:"price_usd,string"`
	PriceBTC         float64 `json:"price_btc,string"`
	USD24HVolume     float64 `json:"24h_volume_usd,string"`
	MarketCapUSD     float64 `json:"market_cap_usd,string"`
	AvailableSupply  float64 `json:"available_supply,string"`
	TotalSupply      float64 `json:"total_supply,string"`
	PercentChange1H  float64 `json:"percent_change_1h,string"`
	PercentChange24H float64 `json:"percent_change_24h,string"`
	PercentChange7D  float64 `json:"percent_change_7d,string"`
	LastUpdated      string  `json:"last_updated"`
}

// GlobalMarketData struct
type GlobalMarketData struct {
	TotalMarketCapUSD            float64 `json:"total_market_cap_usd"`
	Total24HVolumeUSD            float64 `json:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap float64 `json:"bitcoin_percentage_of_market_cap"`
	ActiveCurrencies             int     `json:"active_currencies"`
	ActiveAssets                 int     `json:"active_assets"`
	ActiveMarkets                int     `json:"active_markets"`
}

// CoinGraph struct
type CoinGraph struct {
	MarketCapByAvailableSupply [][]float64 `json:"market_cap_by_available_supply"`
	PriceBTC                   [][]float64 `json:"price_btc"`
	PriceUSD                   [][]float64 `json:"price_usd"`
	VolumeUSD                  [][]float64 `json:"volume_usd"`
}

// Market struct
type Market struct {
	Rank          int
	Exchange      string
	Pair          string
	VolumeUSD     float64
	Price         float64
	VolumePercent float64
	Updated       string
}

// MarketGraph struct
type MarketGraph struct {
	MarketCapByAvailableSupply [][]float64 `json:"market_cap_by_available_supply"`
	VolumeUSD                  [][]float64 `json:"volume_usd"`
}

var (
	baseURL               = "https://api.coinmarketcap.com/v1"
	coinGraphURL          = "https://graphs2.coinmarketcap.com/currencies"
	globalMarketGraphURL  = "https://graphs2.coinmarketcap.com/global/marketcap-total"
	altcoinMarketGraphURL = "https://graphs2.coinmarketcap.com/global/marketcap-altcoin"
)

// GetGlobalMarketData get information about the global market data of the cryptocurrencies
func GetGlobalMarketData() (GlobalMarketData, error) {
	url := fmt.Sprintf(baseURL + "/global/")

	resp, err := makeReq(url)

	var data GlobalMarketData
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return GlobalMarketData{}, err
	}

	return data, nil
}

// GetGlobalMarketGraphData get graph data points of global market
func GetGlobalMarketGraphData(start int64, end int64) (MarketGraph, error) {
	url := fmt.Sprintf("%s/%d/%d", globalMarketGraphURL, start*1000, end*1000)
	resp, err := makeReq(url)
	if err != nil {
		return MarketGraph{}, err
	}
	var data MarketGraph
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return MarketGraph{}, err
	}

	return data, nil
}

// GetAltcoinMarketGraphData get graph data points of altcoin market
func GetAltcoinMarketGraphData(start int64, end int64) (MarketGraph, error) {
	url := fmt.Sprintf("%s/%d/%d", altcoinMarketGraphURL, start*1000, end*1000)
	resp, err := makeReq(url)
	if err != nil {
		return MarketGraph{}, err
	}
	var data MarketGraph
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return MarketGraph{}, err
	}

	return data, nil
}

// GetCoinData get information about a crypto currency
func GetCoinData(coin string) (Coin, error) {
	coin = strings.ToLower(coin)
	url := fmt.Sprintf("%s/ticker/%s", baseURL, coin)
	resp, err := makeReq(url)
	if err != nil {
		return Coin{}, err
	}
	var data []Coin
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return Coin{}, err
	}

	return data[0], nil
}

// GetAllCoinData get information about all coins listed in Coin Market Cap
func GetAllCoinData(limit int) (map[string]Coin, error) {
	var l string
	if limit >= 0 {
		l = fmt.Sprintf("?limit=%v", limit)
	}
	url := fmt.Sprintf("%s/ticker/%s", baseURL, l)

	resp, err := makeReq(url)

	var data []Coin
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	// creating map from the array
	allCoins := make(map[string]Coin)
	for i := 0; i < len(data); i++ {
		allCoins[data[i].ID] = data[i]
	}

	return allCoins, nil
}

// GetCoinGraphData get graph data points for a crypto currency
func GetCoinGraphData(coin string, start int64, end int64) (CoinGraph, error) {
	url := fmt.Sprintf("%s/%s/%d/%d", coinGraphURL, strings.ToLower(coin), start*1000, end*1000)
	resp, err := makeReq(url)
	if err != nil {
		return CoinGraph{}, err
	}
	var data CoinGraph
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return CoinGraph{}, err
	}

	return data, nil
}

// GetCoinPriceUSD get USD price of crypto currency
func GetCoinPriceUSD(coin string) (float64, error) {
	data, err := GetCoinData(strings.ToLower(coin))
	if err != nil {
		return float64(0), nil
	}
	return data.PriceUSD, nil
}

// GetCoinMarkets get market data for a coin name
func GetCoinMarkets(coin string) ([]Market, error) {
	url := fmt.Sprintf("https://coinmarketcap.com/currencies/%s/#markets", strings.ToLower(coin))
	var markets []Market
	response, err := soup.Get(url)
	row := soup.HTMLParse(response).Find("table", "id", "markets-table")

	var rows []soup.Root
	if !reflect.DeepEqual(row, soup.Root{}) {
		row = row.Find("tbody")
		if !reflect.DeepEqual(row, soup.Root{}) {
			rows = row.FindAll("tr")
		}
	}

	if r := recover(); r != nil {
		fmt.Println("birşeyler oldu")
		return nil, err
	}
	if rows == nil || err != nil {
		return nil, err
	}
	for _, row := range rows {
		var data []string
		for _, column := range row.FindAll("td") {
			attrs := column.Attrs()
			if attrs["data-sort"] != "" {
				data = append(data, attrs["data-sort"])
			} else {
				data = append(data, column.Text())
			}
		}
		markets = append(markets, Market{
			Rank:          toInt(data[0]),
			Exchange:      data[1],
			Pair:          data[2],
			VolumeUSD:     toFloat(data[3]),
			Price:         toFloat(data[4]),
			VolumePercent: toFloat(data[5]),
			Updated:       data[6],
		})
	}
	return markets, nil
}

// doReq HTTP client
func doReq(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

// makeReq HTTP request helper
func makeReq(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := doReq(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

// toInt helper for parsing strings to int
func toInt(rawInt string) int {
	parsed, _ := strconv.Atoi(strings.Replace(strings.Replace(rawInt, "$", "", -1), ",", "", -1))
	return parsed
}

// toFloat helper for parsing strings to float
func toFloat(rawFloat string) float64 {
	parsed, _ := strconv.ParseFloat(strings.Replace(strings.Replace(strings.Replace(rawFloat, "$", "", -1), ",", "", -1), "%", "", -1), 64)
	return parsed
}
