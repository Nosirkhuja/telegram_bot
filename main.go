package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ipinfo/go/v2/ipinfo"
	"log"
	"net"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! ")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			ip := net.ParseIP(update.Message.Text)
			if ip == nil {
				log.Println("incorrect ip adress")
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неправильный ip адрес :( Повторите ввод ")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue
			}
			client := ipinfo.NewClient(nil, nil, "")
			info, err := client.GetIPInfo(ip)
			if err != nil {
				log.Println("incorrect ip adress")
			} else {
				infoIp := "Country: " + info.CountryName + " City: " + info.City + " Timezone: " + info.Timezone + " Org: " + info.Org
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, infoIp)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}

}
