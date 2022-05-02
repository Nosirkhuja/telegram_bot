package database

import (
	"bot/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Newdb() *gorm.DB {
	dsn := "host=localhost user=postgres password=123 dbname=postgres port=5434 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(model.Hystory{})
	db.AutoMigrate(model.Admin{})
	return db
}
