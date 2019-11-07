package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/route"
	"github.com/orsenkucher/schedulebot/subs"
	"github.com/orsenkucher/schedulebot/user"
)

var currentRoutes = make(map[user.User]*route.Tree)

func sendOnChan(ch chan subs.SubEvent, e subs.SubEvent) {
	ch <- e
}

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
		user := user.User(update.Message.Chat.ID)
		msg := tgbotapi.NewMessage(int64(user), "Выбери свое расписание👇🏻") // ⬇️ 🎓 👇🏻
		currentRoutes[user] = route.Routes
		msg.ReplyMarkup = GenFor(currentRoutes[user])
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
	chans map[string]chan subs.SubEvent) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	switch {
	case strings.Contains(data, "route"):
		fmt.Println(data)
		childName := strings.Split(data, ":")[1]
		user := user.User(chatID)
		if route, ok := currentRoutes[user]; ok {
			if childRoute, ok := route.Select(childName); ok {
				msg := tgbotapi.NewEditMessageText(int64(user), messageID, fmt.Sprintf("%s👇🏻", childRoute))
				mkp := GenFor(childRoute)
				msg.ReplyMarkup = &mkp
				// msg := tgbotapi.NewEditMessageReplyMarkup(int64(user), messageID, GenFor(childRoute))
				if _, err := bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
					log.Panic(err)
				}
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
				currentRoutes[user] = childRoute
			}
		}
	case strings.Contains(data, "back"):
		fmt.Println(data)
		user := user.User(chatID)
		if route, ok := currentRoutes[user]; ok {
			parent := route.Parent
			msg := tgbotapi.NewEditMessageText(int64(user), messageID, fmt.Sprintf("%s👇🏻", parent))
			mkp := GenFor(parent)
			msg.ReplyMarkup = &mkp
			if _, err := bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
				log.Panic(err)
			}
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}
			currentRoutes[user] = parent
		}
	case strings.Contains(data, "sub"):
		scheduleName := strings.Split(data, ":")[1]
		ch, ok := chans[scheduleName]
		if ok {
			fmt.Println(data)
			go sendOnChan(ch, subs.SubEvent{Action: subs.Add, ChatID: chatID})
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
		}
	case strings.Contains(data, "reset"):
		scheduleName := strings.Split(data, ":")[1]

		ch, ok := chans[scheduleName]
		if ok {
			fmt.Println(data)
			go sendOnChan(ch, subs.SubEvent{Action: subs.Del, ChatID: chatID})
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
