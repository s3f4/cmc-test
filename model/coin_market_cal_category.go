package model

type CoinMarketCalCategory struct {
	CategoryId         uint                 `gorm:"column:category_id;NOT NULL;PRIMARY KEY;AUTO_INCREMENT";`
	Category           string               `gorm:"column:category;"`
}

func (CoinMarketCalCategory) TableName() string {
	return "coin_market_cal_category"
}
