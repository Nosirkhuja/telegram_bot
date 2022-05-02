package api

import (
	cache "bot/cache"
	"bot/database"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
)

func Api() {
	router := gin.Default()
	router.GET("/get_users", get_users)
	router.GET("/get_user/:username", get_user_by_id)
	router.DELETE("/remove/:username/:request")
	router.Run("localhost:8080")
}

func get_users(c *gin.Context) {
	usernames := GetUsernames(cache.NewCache())
	c.IndentedJSON(http.StatusOK, usernames)
}

func get_user_by_id(c *gin.Context) {
	username := c.Param("username")
	requests := GetAllUserRequests(cache.NewCache(), username)
	c.IndentedJSON(http.StatusOK, requests)
}
func delete_user_request(c *gin.Context) {
	username := c.Param("username")
	request := c.Param("request")
	DeleteUserIPRequest(username, request, database.Newdb())
	c.IndentedJSON(http.StatusOK, "Removed")
}
