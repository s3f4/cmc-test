package model

type CoinMarketData struct {
	MarketDataId  int     `gorm:"column:market_data_id;NOT NULL;PRIMARY KEY;"`
	CoinId        int     `gorm:"column:coin_id;"`
	MarketId      int     `gorm:"column:market_id;"`
	Base          string  `gorm:"column:base;"`
	Price_BTC     float64 `gorm:"column:price_btc;"`
	Price_USD     float64 `gorm:"column:price_usd;"`
	TimePeriod    string  `gorm:"column:time_period;"`
	Rank          int     `gorm:"column:rank;"`
	Volume_USD    float64 `gorm:"column:volume_usd;"`
	VolumePercent float64 `gorm:"column:volume_percent;"`
}

func (*CoinMarketData) TableName() string {
	return "coin_market_data"
}