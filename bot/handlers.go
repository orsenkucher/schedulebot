package bot

import (
	"fmt"
	"log"

	"github.com/orsenkucher/schedulebot/route"

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
	mkp, ok := GenFor(b.root.Rootnode)
	if !ok {
		log.Panic("Here must be ok!")
	}
	msg.ReplyMarkup = mkp
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) onReset(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	urt := b.getResetTree(chatID)
	// dropped := rt.Rootnode.Drop()
	mkp, _ := GenForReset(urt.Rootnode)

	msg := tgbotapi.NewMessage(chatID, urt.Rootnode.Drop().String()) //"Варианты отписки ("+update.Message.Chat.FirstName+")")
	msg.ReplyMarkup = mkp
	fmt.Println("Doing reset for user", chatID)
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) getResetTree(userID int64) *route.TreeRoot {
	_, ok := b.resetTree[userID]
	if !ok {
		// fbclient. Fetch here
		subs := [][]string{}
		b.resetTree[userID] = route.NewTreeRoot(route.GenerateUsersTree(subs))
	}

	return b.resetTree[userID]
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
		nodehash := strings.Split(data, ":")[1]
		if node, ok := b.root.Find(nodehash); ok {
			if mkp, ok := GenFor(node); ok {
				msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s👇🏻", node))
				msg.ReplyMarkup = &mkp
				if _, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
					log.Panic(err)
				}
				if _, err := b.api.Send(msg); err != nil {
					log.Panic(err)
				}
			} else {
				schName := node.MakePath()
				ch, ok := chans[schName]
				if ok {
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
	case strings.Contains(data, "reset"):
		nodehash := strings.Split(data, ":")[1]
		urt := b.getResetTree(chatID)
		if node, ok := urt.Find(nodehash); ok {
			if mkp, ok := GenForReset(node); ok {
				msg := tgbotapi.NewEditMessageText(chatID, messageID, node.Drop().String())
				msg.ReplyMarkup = &mkp
				if _, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "")); err != nil {
					log.Panic(err)
				}
				if _, err := b.api.Send(msg); err != nil {
					log.Panic(err)
				}
			} else {
				scheduleName := node.MakePath()
				ch, ok := chans[scheduleName]
				if ok {
					fmt.Println(data)
					ch <- root.SubEvent{Action: root.Del, SubID: chatID}
					fbclient.DeleteSubscriber(chatID, scheduleName)
					// snackMsg := "Un️subscribed ♻️" // ☠️
					// snackMsg := "Отписка проведена ♻️ (" + cmdMapping[data] + ")"
					snackMsg := "Отписка проведена ♻️" + node.Name
					b.api.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
					msg := tgbotapi.NewMessage(chatID, snackMsg)
					if _, err := b.api.Send(msg); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
