package datasource

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"fmt"
)

type TwitterApi struct {
	client *twitter.Client
	Config TwitterConfig
}

type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

//CreateClient
func (t *TwitterApi) CreateClient() {

	config := oauth1.NewConfig(
		t.Config.ConsumerKey,
		t.Config.ConsumerSecret,
	)
	token := oauth1.NewToken(t.Config.Token, t.Config.TokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	t.client = twitter.NewClient(httpClient)
}

//TimeLine
func (t *TwitterApi) TimeLine(count int) []twitter.Tweet {
	tweets, _, err := t.client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: count,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	return tweets
}

//GetTweet
func (t *TwitterApi) GetTweet(tweetId int64) *twitter.Tweet {
	tweet, _, err := t.client.Statuses.Show(tweetId, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	return tweet
}

/*
	filterString = -filter:links
	lang = en
 */
//Search
func (t *TwitterApi) Search(searchString, langString, filterString string) *twitter.Search {
	search, _, err := t.client.Search.Tweets(&twitter.SearchTweetParams{
		Query:     searchString + "-filter:" + filterString + " lang:" + langString,
		TweetMode: "extended",
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	return search
}

//GetUser
func (t *TwitterApi) GetUser(userName string) *twitter.User {
	user, _, err := t.client.Users.Show(&twitter.UserShowParams{
		ScreenName: userName,
	})

	if err != nil {
		fmt.Println(err.Error())
	}

	return user
}

//GetFollowerList
func (t *TwitterApi) GetFollowerList() *twitter.Followers {
	followers, _, err := t.client.Followers.List(&twitter.FollowerListParams{})
	if err != nil {
		fmt.Println(err.Error())
	}

	return followers
}

//NewTweet
func (t *TwitterApi) NewTweet(tweetString string) *twitter.Tweet {
	tweet, _, err := t.client.Statuses.Update(tweetString, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	return tweet
}
