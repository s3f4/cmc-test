package controller

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
	"strconv"
)

type BaseController struct {
	DB                 *gorm.DB
	DBConnectionStatus bool
}

//DBConnect
func (baseController *BaseController) DBConnect() {
	var err error
	baseController.DB, err = gorm.Open("mysql", "root:@/crypto_currency?charset=utf8&parseTime=True&loc=Local")
	baseController.DB.LogMode(true)
	if err != nil {
		fmt.Println("connection error")
	} else {
		baseController.DBConnectionStatus = true
	}

}

//DBCloseConnection
func (baseController *BaseController) close() {
	baseController.DBConnectionStatus = false
	baseController.DB.Close()
}

//Destruct
func (baseController *BaseController) Destruct() {
	baseController.close()
}

//timestampToTime
func (*BaseController) timestampToTime(tsStr string) time.Time {

	if tsStr == "" {
		return time.Unix(time.Now().Unix(), 0)

	}

	timestamp, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(timestamp, 0)
	return tm
}
