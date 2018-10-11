package model

type FounderTwitter struct {
	CoinTwitterId int    `gorm:"column:coin_twitter_id;NOT NULL;PRIMARY KEY;"`
	CoinId        int    `gorm:"column:coin_id;"`
	Founder       string `gorm:"column:founder;"`
}
