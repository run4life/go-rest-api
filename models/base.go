package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var conn *gorm.DB
	var err error

	conn, err = gorm.Open("mysql",
	"root:root@/gorest?charset=utf8")

	if err != nil {
		fmt.Print("Connect to db failed: " + err.Error())
	}

	db = conn

	// Database migration
	db.Debug().AutoMigrate(&Account{}, &Profile{})
}

/*
 GetDB Get a db instance
*/
func GetDB() *gorm.DB {
	return db
}
