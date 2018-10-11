package datasource

import (
	"encoding/json"
	"strings"
	"fmt"
	"net/http"
	"io/ioutil"
)

const apiUrl = "https://api.coinmarketcal.com"

type CoinMarketCal struct {
	access_token string
	Config       Config
}

type Config struct {
	ClientId     string
	ClientSecret string
}

var (
	response *http.Response
	err      error
	body     []byte
)

//GetCategories
func (coinmarketcal *CoinMarketCal) GetCategories() ([]byte) {
	return coinmarketcal.get("categories")
}

//GetEvents
func (coinmarketcal *CoinMarketCal) GetEvents() ([]byte) {
	return coinmarketcal.get("events")
}

//GetCoins
func (coinmarketcal *CoinMarketCal) GetCoins() ([]byte) {
	return coinmarketcal.get("coins")
}

//refreshToken
func (coinmarketcal *CoinMarketCal) refreshToken() string {
	response := coinmarketcal.get("access_token")

	var jsonResult map[string]interface{}
	_ = json.Unmarshal(response, &jsonResult)

	return jsonResult["access_token"].(string)
}

//tokenIsExpired
func (coinmarketcal *CoinMarketCal) tokenIsExpired(request []byte) (interface{}) {
	var errJson map[string]interface{}
	err = json.Unmarshal(request, &errJson)
	if errJson["error"] != nil && strings.EqualFold(errJson["error_description"].(string), "The access token provided is invalid.") {
		return coinmarketcal.refreshToken()
	} else {
		return nil
	}
}

//get all requests
func (coinmarketcal *CoinMarketCal) get(get string) ([]byte) {

	if coinmarketcal.access_token == "" {
		coinmarketcal.access_token = "access-token-default"
	}

	var queryString = "/v1/events?access_token=" + coinmarketcal.access_token

	switch get {
	case "coins":
		queryString = "/v1/coins?access_token=" + coinmarketcal.access_token
	case "events":
		queryString = "/v1/events?access_token=" + coinmarketcal.access_token
	case "access_token":
		queryString = fmt.Sprintf("/oauth/v2/token?grant_type=client_credentials&client_id=%s&client_secret=%s", coinmarketcal.Config.ClientId, coinmarketcal.Config.ClientSecret)
	case "categories":
		queryString = "/v1/categories?access_token=" + coinmarketcal.access_token
	default:
		queryString = "/v1/events?access_token=" + coinmarketcal.access_token
	}

	source := Request(apiUrl, queryString, "GET")

	if coinmarketcal.tokenIsExpired(source) != nil {
		coinmarketcal.access_token = coinmarketcal.refreshToken()
		return coinmarketcal.get(get)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Sprintf("recovered %s", r.(string)))
			return
		}
	}()

	var v interface{}
	if err = json.Unmarshal(source, &v); err != nil {
		panic(err.Error())
	}

	return source

}

//Request
func Request(site string, query string, method string) []byte {

	if method == "GET" {
		response, err = http.Get(site + query)
	} else {
	}

	if err != nil {
		return []byte(err.Error())
	}

	body, _ = ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	return body
}
