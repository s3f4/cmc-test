package model

type CoinMarket struct {
	CoinMarketId int `gorm:"column:coin_market_id;NOT NULL;PRIMARY KEY;"`
	MarketId     int `gorm:"column:market_id;"`
	CoinId       int `gorm:"column:coin_id;"`
}

func (*CoinMarket) TableName() string {
	return "coin_market"
}
