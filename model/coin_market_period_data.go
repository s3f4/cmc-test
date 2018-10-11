package model

type CoinMarketPeriodData struct {
	DataId      int     `gorm:"column:data_id;NOT NULL;PRIMARY KEY;"`
	QuoteId     int     `gorm:"column:quote_id;index:coin_quote_period"`
	Open        float64 `gorm:"column:open;"`
	Close       float64 `gorm:"column:close;"`
	High        float64 `gorm:"column:high;"`
	Low         float64 `gorm:"column:low;"`
	Volume      float64 `gorm:"column:volume;"`
	Period      string  `gorm:"column:period;index:coin_quote_period"`
	LastUpdated string  `gorm:"column:last_updated;"`
	Table       string
}

//coin_market_data_marketId_coinSymbol
func (dmpc CoinMarketPeriodData) SetTableName(tableName string) {
	dmpc.Table = tableName
}

func (dmpc CoinMarketPeriodData) TableName() string {
	return dmpc.Table
}
