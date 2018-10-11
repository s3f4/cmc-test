package model

import (
	"time"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type CoinMarketCalEvent struct {
	EventId uint `gorm:"primary_key;AUTO_INCREMENT;"`
	Title   string
	Coins []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"coins"`
	DateEvent         time.Time `json:"date_event"`
	CreatedDate       time.Time `json:"created_date"`
	Description       string    `json:"description"`
	Proof             string    `json:"proof"`
	Source            string    `json:"source"`
	IsHot             bool      `json:"is_hot"`
	VoteCount         int       `json:"vote_count"`
	PositiveVoteCount int       `json:"positive_vote_count"`
	Percentage        int       `json:"percentage"`
	Categories []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"categories" gorm:"-"`
	CategoryId int64 `gorm:"foreignkey:CategoryId;"`
	TipSymbol      string `json:"tip_symbol" gorm:"-"`
	TipAdress      string `json:"tip_adress" gorm:"-"`
	TwitterAccount string `json:"twitter_account"`
	CanOccurBefore bool   `json:"can_occur_before"`
}

func (CoinMarketCalEvent) TableName() string {
	return "coin_market_cal_event"
}
