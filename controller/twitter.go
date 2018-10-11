package controller

import (
	"fmt"
	"dehaa.com/core/datasource"
)

type TwitterController struct {
	*BaseController
}

func (tc *TwitterController) test() {

}

func TwitterTest() {
	t := new(datasource.TwitterApi)
	t.CreateClient()
	user := t.GetUser("x")
	fmt.Println(user)
	followerList := t.GetFollowerList()
	fmt.Println(followerList)
	timeLine := t.TimeLine(20)
	fmt.Println(timeLine)
}
