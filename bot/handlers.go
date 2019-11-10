package bot

import (
	"fmt"
	"log"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/root"
)

func (b *Bot) handleCommand(update tgbotapi.Update) {
	switch update.Message.Command() {
	case "sub", "start", "go":
		b.onSub(update)
	case "reset", "unsub":
		b.onReset(update)
	default:
		return
	}
}

func (b *Bot) handleMessage(update tgbotapi.Update) {
	// // Do nothing
	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// if _, err := bot.Send(msg); err != nil {
	// 	log.Panic(err)
	// }
}

func (b *Bot) onSub(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "Выбери свое расписание👇🏻") // ⬇️ 🎓 👇🏻
	// user.Route = route.Routes
	mkp, ok := GenFor(b.rootnode)
	if !ok {
		log.Panic("Here must be ok!")
	}
	msg.ReplyMarkup = mkp
	if _, err := b.api.Send(msg); err != nil {
		log.Panic(err)
	}
}

func (b *Bot) onReset(update tgbotapi.Update) {
	// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
	// 		"Варианты отписки ("+update.Message.Chat.FirstName+")") // Unsub options ("+update.Message.Chat.FirstName+")"
	// 	msg.ReplyMarkup = inlineResetKeyboard
	// 	fmt.Println("Doing reset for user", update.Message.Chat.ID)
	// 	if _, err := b.api.Send(msg); err != nil {
	// 		log.Panic(err)
	// 	}
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
		nodepath := strings.Split(data, ":")[1]
		if node, ok := b.rootnode.Find(nodepath); ok {
			msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s👇🏻", node))
			if mkp, ok := GenFor(node); ok {
				msg.ReplyMarkup = &mkp
				if _, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
					log.Panic(err)
				}
				if _, err := b.api.Send(msg); err != nil {
					log.Panic(err)
				}
			} else {
				schName := nodepath
				ch, ok := chans[schName]
				if ok {
					fmt.Println(data)
					ch <- root.SubEvent{Action: root.Add, SubID: chatID}
					fbclient.AddSubscriber(chatID, schName)
					// snackMsg := "Our congrats 🥂. We handled your sub!"
					// snackMsg := "Ваша регистрация обработана 🥂 (" + cmdMapping[data] + ")"
					snackMsg := "Поздравляю! Ты подписался на \"" + node.Name + "\". Увидимся на паре 🥂"
					b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
					msg := tgbotapi.NewMessage(chatID, snackMsg)
					if _, err := b.api.Send(msg); err != nil {
						log.Println(err)
					}
				} else {
					snackMsg := "Этот шедуль уже не в базе"
					b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
					log.Printf("Schedule was not found by name: %s\n", schName)
				}
			}
		}
		// case strings.Contains(data, "reset"):
		// 	scheduleName := strings.Split(data, ":")[1]

		// 	ch, ok := chans[scheduleName]
		// 	if ok {
		// 		fmt.Println(data)
		// 		ch <- root.SubEvent{Action: root.Del, SubID: chatID}
		// 		fbclient.DeleteSubscriber(chatID, scheduleName)
		// 		// snackMsg := "Un️subscribed ♻️" // ☠️
		// 		snackMsg := "Отписка проведена ♻️ (" + cmdMapping[data] + ")"
		// 		b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
		// 		msg := tgbotapi.NewMessage(chatID, snackMsg)
		// 		if _, err := b.api.Send(msg); err != nil {
		// 			log.Println(err)
		// 		}
		// 	}
	}
}
