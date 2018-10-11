package model

type Market struct {
	MarketId   int    `gorm:"column:market_id;NOT NULL;PRIMARY KEY;"`
	MarketName string `gorm:"column:market_name;"`
	MarketCode string `gorm:"column:market_code;"`
}

func (*Market) TableName() string {
	return "market"
}