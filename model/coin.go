package model

type Coin struct {
	/**
		Associations
	 */
	FounderTwitter                           `gorm:"foreignkey:coin_id;association_key:coin_Id;"`
	FounderTweet                             `gorm:"foreignkey:coin_id;association_key:coin_Id;"`
	Tweets              []Tweet              `gorm:"foreignkey:coin_id;association_key:coin_Id;"`
	CoinMarketCalEvents []CoinMarketCalEvent `gorm:"foreignkey:coin_id;association_key:coin_Id;"`
	CoinMarkets         []CoinMarket         `gorm:"foreignkey:coin_id;association_key:coin_Id;"`

	CoinId          int     `gorm:"column:coin_id;primary_key;AUTO_INCREMENT"`
	Name            string  `gorm:"column:name;"`
	Symbol          string  `gorm:"column:symbol;index:symbol"`
	AvailableSupply float64 `gorm:"column:available_supply;"`
	TotalSupply     float64 `gorm:"column:total_supply;"`
	Rank            int     `gorm:"column:rank;"`
	LastUpdated     string  `gorm:"column:last_updated;"`
}

func (Coin) TableName() string {
	return "coin"
}
