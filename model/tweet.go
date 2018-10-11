package model

import (
	"time"
)

type Tweet struct {
	Tweet_id   int       `gorm:"column:tweet_id;NOT NULL;PRIMARY KEY;"`
	Coin_id    int       `gorm:"column:coin_id;"`
	Account    string    `gorm:"column:account;"`
	Tweet      string    `gorm:"column:tweet;"`
	Tweet_date time.Time `gorm:"column:tweet_date;"`
}
