package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/fbclient"
)

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// // Do nothing
	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// if _, err := bot.Send(msg); err != nil {
	// 	log.Panic(err)
	// }
}

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "sub", "start", "go":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери свое расписание👇🏻") // ⬇️ 🎓 👇🏻
		msg.ReplyMarkup = inlineKeyboard
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	case "reset", "unsub":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Варианты отписки ("+update.Message.Chat.FirstName+")") // Unsub options ("+update.Message.Chat.FirstName+")"
		msg.ReplyMarkup = inlineResetKeyboard
		fmt.Println("Doing reset for user", update.Message.Chat.ID)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	default:
		return
	}
}

func handleCallback(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	chans map[string]chan SubEvent) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	scheduleName := strings.Split(data, ":")[1]
	ch, ok := chans[scheduleName]
	if ok {
		switch {
		case strings.Contains(data, "sub"):
			fmt.Println(data)
			go sendOnChan(ch, SubEvent{Action: Add, ChatID: chatID})
			fbclient.AddSubscriber(chatID, scheduleName)
			// snackMsg := "Our congrats 🥂. We handled your sub!"
			// Поздравляю! Ты подписался на бота! До скорых встреч на паре!
			snackMsg := "Поздравляю! Ты подписался на \"" + cmdMapping[data] + "\". Увидимся на паре 🥂"
			// snackMsg := "Ваша регистрация обработана 🥂 (" + cmdMapping[data] + ")"
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
			msg := tgbotapi.NewMessage(chatID, snackMsg)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
		case strings.Contains(data, "reset"):
			fmt.Println(data)
			go sendOnChan(ch, SubEvent{Action: Del, ChatID: chatID})
			fbclient.DeleteSubscriber(chatID, scheduleName)
			// snackMsg := "Un️subscribed ♻️" // ☠️
			snackMsg := "Отписка проведена ♻️ (" + cmdMapping[data] + ")"
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
			msg := tgbotapi.NewMessage(chatID, snackMsg)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}
