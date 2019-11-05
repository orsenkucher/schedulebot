package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/fbclient"
)

// SubEvent represents subscription event
type SubEvent struct {
	ChatID int64
	Action SubAction
}

// SubAction represents user action
type SubAction int

// Add is when user Success
// Del is when user unsubbed
const (
	_ SubAction = iota
	Add
	Del
)

func sendOnChan(ch chan SubEvent, e SubEvent) {
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
		user := user(update.Message.Chat.ID)
		msg := tgbotapi.NewMessage(int64(user), "Выбери свое расписание👇🏻") // ⬇️ 🎓 👇🏻
		currentRoutes[user] = Routes
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
	chans map[string]chan SubEvent) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	switch {
	case strings.Contains(data, "route"):
		fmt.Println(data)
		childName := strings.Split(data, ":")[1]
		user := user(chatID)
		if route, ok := currentRoutes[user]; ok {
			msg := tgbotapi.NewMessage(int64(user), route.name+" выбери свое расписание👇🏻") // ⬇️ 🎓 👇🏻
			childRoute, err := route.Select(childName)
			if err == nil {
				msg.ReplyMarkup = GenFor(childRoute)
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
			currentRoutes[user] = childRoute
		}
	case strings.Contains(data, "back"):
		fmt.Println(data)
		user := user(chatID)
		if route, ok := currentRoutes[user]; ok {
			parent := route.parent
			msg := tgbotapi.NewMessage(int64(user), parent.name+" выбери свое расписание👇🏻") // ⬇️ 🎓 👇🏻
			msg.ReplyMarkup = GenFor(parent)
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
		}
	case strings.Contains(data, "reset"):
		scheduleName := strings.Split(data, ":")[1]

		ch, ok := chans[scheduleName]
		if ok {
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
