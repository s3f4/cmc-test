package migration

import (
	"dehaa.com/core/model"
	"github.com/jinzhu/gorm"
	"strconv"
)

func CoinCheckTable(coin *model.Coin, db *gorm.DB) (bool, *model.Coin) {

	newCoinTable := model.CoinMarketPeriodData{}
	newCoinTable.Table = "coin_market_data_" + strconv.Itoa(coin.CoinId)
	if !db.HasTable(&newCoinTable) {
		db.CreateTable(&newCoinTable)
		return true, coin
	}
	return false, coin
}
