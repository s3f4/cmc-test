package model

type GlobalMarketData struct {
	DataId                       int     `gorm:"column:data_id;NOT NULL;PRIMARY KEY;"`
	TotalMarketCapUSD            float64 `gorm:"column:total_market_cap_usd;" 			json:"total_market_cap_usd"`
	Total24HVolumeUSD            float64 `gorm:"column:total_24h_volume_usd;" 			json:"total_24h_volume_usd"`
	BitcoinPercentageOfMarketCap float64 `gorm:"column:bitcoin_percentage_of_market_cap;" json:"bitcoin_percentage_of_market_cap"`
	Active_currencies            int     `gorm:"column:active_currencies;" 				json:"active_currencies"`
	Active_assets                int     `gorm:"column:active_assets;" 					json:"active_assets"`
	Active_markets               int     `gorm:"column:active_markets;" 					json:"active_markets"`
}
