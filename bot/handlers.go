package bot

import (
	"fmt"
	"log"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/root"
	"github.com/orsenkucher/schedulebot/route"
)

func (b *Bot) handleMessage(update tgbotapi.Update) {
	// // Do nothing
	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// if _, err := bot.Send(msg); err != nil {
	// 	log.Panic(err)
	// }
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "sub", "start", "go":
		user := b.userByID(update.Message.Chat.ID)
		msg := tgbotapi.NewMessage(user.ID, "Выбери свое расписание👇🏻") // ⬇️ 🎓 👇🏻
		user.Route = route.Routes
		msg.ReplyMarkup = GenFor(user.Route)
		if _, err := b.api.Send(msg); err != nil {
			log.Panic(err)
		}
	// case "reset", "unsub":
	// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
	// 		"Варианты отписки ("+update.Message.Chat.FirstName+")") // Unsub options ("+update.Message.Chat.FirstName+")"
	// 	msg.ReplyMarkup = inlineResetKeyboard
	// 	fmt.Println("Doing reset for user", update.Message.Chat.ID)
	// 	if _, err := b.api.Send(msg); err != nil {
	// 		log.Panic(err)
	// 	}
	default:
		return
	}
}

func (b *Bot) handleCallback(
	update tgbotapi.Update,
	chans map[string]chan root.SubEvent) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	messageID := update.CallbackQuery.Message.MessageID
	switch {
	case strings.Contains(data, "route"):
		fmt.Println(data)
		childName := strings.Split(data, ":")[1]
		user := b.userByID(chatID)
		if childRoute, ok := user.Route.Select(childName); ok {
			msg := tgbotapi.NewEditMessageText(user.ID, messageID, fmt.Sprintf("%s👇🏻", childRoute))
			mkp := GenFor(childRoute)
			msg.ReplyMarkup = &mkp
			// msg := tgbotapi.NewEditMessageReplyMarkup(int64(user), messageID, GenFor(childRoute))
			if _, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
				log.Panic(err)
			}
			if _, err := b.api.Send(msg); err != nil {
				log.Panic(err)
			}
			user.Route = childRoute
		}
	case strings.Contains(data, "sub"):
		scheduleName := strings.Split(data, ":")[1]
		ch, ok := chans[scheduleName]
		if ok {
			fmt.Println(data)
			ch <- root.SubEvent{Action: root.Add, SubID: chatID}
			fbclient.AddSubscriber(chatID, scheduleName)
			// snackMsg := "Our congrats 🥂. We handled your sub!"
			// Поздравляю! Ты подписался на бота! До скорых встреч на паре!
			snackMsg := "Поздравляю! Ты подписался на \"" + cmdMapping[data] + "\". Увидимся на паре 🥂"
			// snackMsg := "Ваша регистрация обработана 🥂 (" + cmdMapping[data] + ")"
			b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
			msg := tgbotapi.NewMessage(chatID, snackMsg)
			if _, err := b.api.Send(msg); err != nil {
				log.Println(err)
			}
		}
	case strings.Contains(data, "back"):
		fmt.Println(data)
		user := b.userByID(chatID)
		if parent := user.Route.Parent; parent != nil {
			msg := tgbotapi.NewEditMessageText(user.ID, messageID, fmt.Sprintf("%s👇🏻", parent))
			mkp := GenFor(parent)
			msg.ReplyMarkup = &mkp
			if _, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
				log.Panic(err)
			}
			if _, err := b.api.Send(msg); err != nil {
				log.Panic(err)
			}
			user.Route = parent
		}
	case strings.Contains(data, "reset"):
		scheduleName := strings.Split(data, ":")[1]

		ch, ok := chans[scheduleName]
		if ok {
			fmt.Println(data)
			ch <- root.SubEvent{Action: root.Del, SubID: chatID}
			fbclient.DeleteSubscriber(chatID, scheduleName)
			// snackMsg := "Un️subscribed ♻️" // ☠️
			snackMsg := "Отписка проведена ♻️ (" + cmdMapping[data] + ")"
			b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
			msg := tgbotapi.NewMessage(chatID, snackMsg)
			if _, err := b.api.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}
