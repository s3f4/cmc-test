package model

type FounderTweet struct {
	FounderTweetId int       `gorm:"column:founder_tweet_id;NOT NULL;PRIMARY KEY;"`
	CoinId         int       `gorm:"column:coin_id;"`
	Founder        string    `gorm:"column:founder;"`
	Tweet          string    `gorm:"column:tweet;"`
}
