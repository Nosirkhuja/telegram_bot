package cache

import (
	"bot/model"
	"sync"
)

// Кэш для быстрого получения данных, если пользователь/админ захочет посмотреть историю запросов
//При рестарте сервера все данные из базы загружаются в кэш.
//При поиске IP данные добавляются в кэш и в базу

type Cache struct {
	Cache model.Data
}

func NewCache() *Cache {
	cc := make(map[string]map[string]string)
	return &Cache{Cache: cc}
}

var m sync.Mutex

func (c *Cache) Set(username string, ipAdress string, result string) {
	m.Lock()
	_, ok := c.Cache[username]
	if !ok {
		m := make(map[string]string)
		c.Cache[username] = m
	}

	c.Cache[username][ipAdress] = result
	m.Unlock()
}

func (c *Cache) Get(userId string) map[string]string {
	return c.Cache[userId]
}

func (cache *Cache) Checkcashe(username, ip string) bool {
	tmp := cache.Get(username)
	_, b := tmp[ip]
	if b {
		return true
	}
	return false
}
