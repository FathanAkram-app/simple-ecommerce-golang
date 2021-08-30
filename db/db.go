package db

import (
	"depmod/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func Init() {
	conf := config.GetConfig()
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		conf.DB_HOST,
		conf.DB_PORT,
		conf.DB_USERNAME,
		conf.DB_NAME,
		conf.DB_PASSWORD)

	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

}

func CreateCon() *gorm.DB {
	return db
}
