package db

import (
	"github.com/jinzhu/gorm"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "data.sqlite3")
	if err != nil {
		panic(err)
	}
	return db
}
