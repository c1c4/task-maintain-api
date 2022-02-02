package database

import (
	"api/app/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() {
	var err error
	Database, err = gorm.Open(mysql.Open(config.DBURL), &gorm.Config{})

	if err != nil {
		log.Fatal("it's not possible to connect with the database", err)
	}
}
