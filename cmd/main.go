package main

import "C"
import (
	"bot/api"
	"bot/cache"
	"bot/database"
	"bot/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ipinfo/go/v2/ipinfo"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {

	db := database.Newdb()
	cache := cache.NewCache()

	fmt.Println("Введите username для админа")
	var defAdmin string
	fmt.Fscan(os.Stdin, &defAdmin)
	database.AddAdmin(db, defAdmin)

	database.Dbtocache(db, cache)

	client := ipinfo.NewClient(nil, nil, "c16b8e4834c983")

	bot, err := tgbotapi.NewBotAPI("5383835265:AAHApOOvkFUsiz7SWcy9X_uSfRYMiqi2eX0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	api.Api()
	for update := range updates {
		if update.Message == nil {
			continue
		}
		isAdmin := database.IsAdmin(db, update.Message.From.UserName)
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать!")
			msg.ReplyToMessageID = update.Message.MessageID
			if isAdmin {
				msg.ReplyMarkup = numericKeyboardAdmin
			} else {
				msg.ReplyMarkup = numericKeyboardUser
			}
			bot.Send(msg)
		case "История":
			all := history(cache, update.Message.From.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, all)
			msg.ReplyToMessageID = update.Message.MessageID
			if isAdmin {
				msg.ReplyMarkup = numericKeyboardAdmin
			} else {
				msg.ReplyMarkup = numericKeyboardUser
			}
			bot.Send(msg)
		case "Поиск IP":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите Ip адрес (X.X.X.X)")
			msg.ReplyToMessageID = update.Message.MessageID
			if isAdmin {
				msg.ReplyMarkup = numericKeyboardAdmin
			} else {
				msg.ReplyMarkup = numericKeyboardUser
			}
			bot.Send(msg)
			ipRequest := get_message(updates)
			infoIp, err := getInfo(client, ipRequest)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, infoIp)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			if err == nil {
				userId := fmt.Sprint(update.Message.From.ID)
				if !cache.Checkcashe(update.Message.From.UserName, ipRequest) {
					cache.Set(update.Message.From.UserName, ipRequest, infoIp)
					database.Addtodb(
						db,
						userId,
						ipRequest,
						infoIp,
						update.Message.From.UserName,
					)
				}
			}
		case "Добавить админа":
			if isAdmin {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, " Введите username пользователя")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

				var addUsername = get_message(updates)
				database.AddAdmin(db, addUsername)

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Пользователь %s назначен админом!", addUsername)))
			}
		case "Удалить админа":
			if isAdmin {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, " Введите username админа для удаления")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

				var RemoveUsername = get_message(updates)
				database.RemoveAdmin(db, RemoveUsername)
				if isAdmin {
					msg.ReplyMarkup = numericKeyboardAdmin
				} else {
					msg.ReplyMarkup = numericKeyboardUser
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Пользователь %s больше не админ!", RemoveUsername)))
			}
		case "Проверить пользователя":
			if isAdmin {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, " Введите username для проверки")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)

				all := historyAdmin(cache, get_message(updates))
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, all)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		case "Рассылка":
			if isAdmin {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите сообщение для отправки пользователям:")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				MassageToUsers := get_message(updates)
				var result model.Hystory
				rows, _ := db.Model(&model.Hystory{}).Rows()
				users := make(map[string]int64)
				for rows.Next() {
					db.ScanRows(rows, &result)
					varInt, _ := strconv.ParseInt(result.Id, 0, 64)
					_, exist := users[result.Id]
					if !exist && update.Message.From.ID != varInt {
						bot.Send(tgbotapi.NewMessage(varInt, MassageToUsers))
						users[result.Id] = varInt
					}

				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Сообщение отправлено всем пользователям!")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}

		}

	}
}
