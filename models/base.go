package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var db *gorm.DB

func init() {
	var conn *gorm.DB
	var err error
	sqlite3DataPath := os.Getenv("go-jwt-simple-auth-sqlite3-path")
	if sqlite3DataPath == "" {
		sqlite3DataPath = "data.db"
	}

	fmt.Println("Use sqlite3: " + sqlite3DataPath)

	conn, err = gorm.Open("sqlite3", sqlite3DataPath)

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
