package database

import (
	"bot/cache"
	"bot/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Addtodb(db *gorm.DB, id, ip, result, username string) error {
	if err := db.Create(&model.Hystory{IpId: ip + id, Ip: ip, Result: result, Id: id, Username: username}).Error; err != nil {
		return err
	}
	log.Println("Данные добавлены в базу данных")
	return nil
}

func Dbtocache(db *gorm.DB, cache *cache.Cache) {
	var hystories []model.Hystory
	db.Find(&hystories)
	for _, hystories := range hystories {
		cache.Set(hystories.Username, hystories.Ip, hystories.Result)
	}
}

func IsAdmin(db *gorm.DB, username string) bool {
	count := int64(0)
	db.Model(&model.Admin{}).Where("Username = ? ", username).Count(&count)
	return count > 0
}

func AddAdmin(db *gorm.DB, username string) {
	db.Create(&model.Admin{Username: username})
}

func RemoveAdmin(db *gorm.DB, username string) {
	db.Delete(model.Admin{}, "Username LIKE ?", username)
}
