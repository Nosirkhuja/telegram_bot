package main

import "C"
import (
	"bot/cache"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ipinfo/go/v2/ipinfo"
	log "github.com/sirupsen/logrus"
	"net"
)

func history(cache *cache.Cache, userId string) string {
	result := "История поиска: "
	idinfo := make(map[string]string)
	idinfo = cache.Get(userId)
	for ipAdress, info := range idinfo {
		result = result + ipAdress + " ---> " + info + "\n"
	}
	return result
}

func historyAdmin(cache *cache.Cache, userId string) string {
	result := "История поиска пользователя @" + userId + " :\n"
	idinfo := make(map[string]string)
	idinfo = cache.Get(userId)
	for ipAdress, _ := range idinfo {
		result = result + ipAdress + "\n"
	}
	return result
}

func get_message(updates tgbotapi.UpdatesChannel) (msg string) {

	for update := range updates {
		if update.Message == nil {
			// ignore any non-Message Updates
			continue
		}
		msg = update.Message.Text
		break
	}

	return msg
}

func getInfo(client *ipinfo.Client, ipAdress string) (infoIp string, err error) {
	ip := net.ParseIP(ipAdress)
	if ipAdress == "/start" {
		return "Добро пожаловать!", nil
	}
	if ip == nil {
		return "Неправильный ip адрес :( Повторите ввод ", errors.New("Invalid IP")
	}
	info, err := client.GetIPInfo(ip)
	if err != nil {
		log.Println("не удалось распознать адрес")
		return "не удалось распознать адрес", err
	}

	infoIp = "Country: " + info.CountryName + " City: " + info.City + " Timezone: " + info.Timezone + " Org: " + info.Org

	return infoIp, nil

}

var numericKeyboardUser = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Поиск IP"),
		tgbotapi.NewKeyboardButton("История"),
	),
)

var numericKeyboardAdmin = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Поиск IP"),
		tgbotapi.NewKeyboardButton("История"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить админа"),
		tgbotapi.NewKeyboardButton("Удалить админа"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Рассылка"),
		tgbotapi.NewKeyboardButton("Проверить пользователя"),
	),
)
