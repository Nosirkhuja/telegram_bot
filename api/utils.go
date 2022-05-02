package api

import (
	"bot/cache"
	"bot/model"
	"gorm.io/gorm"
)

func GetUsernames(c *cache.Cache) (usernames []string) {
	for k, _ := range c.Cache {
		usernames = append(usernames, k)
	}
	return
}

func GetAllUserRequests(c *cache.Cache, username string) map[string]string {
	return c.Cache[username]
}

func DeleteUserIPRequest(username, request string, db *gorm.DB) {
	db.Delete(&model.Hystory{
		Username: username,
		Ip:       request,
	})
}
