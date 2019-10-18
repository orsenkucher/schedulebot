package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/fbclient"
)

// InitBot initializes telegram bot
func InitBot(withKey string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(withKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
	_, err = bot.RemoveWebhook()
	if err != nil {
		log.Println("Cant remove webhook")
	}

	// https://schedulebot-x2gm2h2g4a-uc.a.run.app
	hook := tgbotapi.NewWebhook("https://schedulebot-x2gm2h2g4a-uc.a.run.app:8443/" + bot.Token)
	_, err = bot.SetWebhook(hook)
	if err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	return bot
}

// Listen starts infinite listening
func Listen(bot *tgbotapi.BotAPI, chans map[string]chan SubEvent) {
	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	// if err != nil {
	// 	log.Panic(err)
	// }

	// updates := bot.ListenForWebhook("/" + bot.Token)

	// for update := range updates {
	// 	if update.CallbackQuery != nil {
	// 		handleCallback(bot, update, chans)
	// 		continue
	// 	}

	// 	if update.Message != nil {
	// 		if update.Message.IsCommand() {
	// 			handleCommand(bot, update)
	// 		} else if update.Message.Text != "" {
	// 			handleMessage(bot, update)
	// 		}
	// 		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	// 		continue
	// 	}
	// }
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" 🥞 1 група  ", "sub:group1"),
		tgbotapi.NewInlineKeyboardButtonData(" 🍇 2 група  ", "sub:group2"),
	),
	// tgbotapi.NewInlineKeyboardRow(
	// 	tgbotapi.NewInlineKeyboardButtonData(" 🤹 demo  ", "sub:test"),
	// ),
)

var inlineResetKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" 🖕🏾 1 група  ", "reset:group1"),
		tgbotapi.NewInlineKeyboardButtonData(" 🖕🏾 2 група  ", "reset:group2"),
	),
	// tgbotapi.NewInlineKeyboardRow(
	// 	tgbotapi.NewInlineKeyboardButtonData(" 🖕🏾 demo  ", "reset:test"),
	// ),
)

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "sub", "start", "go":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Go!")
		msg.ReplyMarkup = inlineKeyboard
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	case "reset":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Doing reset for "+update.Message.Chat.FirstName)
		msg.ReplyMarkup = inlineResetKeyboard
		fmt.Println("Doing reset for user", update.Message.Chat.ID)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	default:
		return
	}
}

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

func handleCallback(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
	chans map[string]chan SubEvent) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	ch, ok := chans[strings.Split(data, ":")[1]]
	if ok {
		switch {
		case strings.Contains(data, "sub"):
			fmt.Println(data)
			go sendOnChan(ch, SubEvent{Action: Add, ChatID: chatID})
			snackMsg := "Our congrats 🥂. We handled your sub!"
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))
		case strings.Contains(data, "reset"):
			fmt.Println(data)
			go sendOnChan(ch, SubEvent{Action: Del, ChatID: chatID})
			snackMsg := "Un️subscribed ☠️"
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, snackMsg))

		}
	}
}

func sendOnChan(ch chan SubEvent, e SubEvent) {
	ch <- e
}

// SpreadMessage is public
func SpreadMessage(b *tgbotapi.BotAPI, users []int64, msg string) error {
	for _, u := range users {
		tgmsg := tgbotapi.NewMessage(u, msg)
		if _, err := b.Send(tgmsg); err != nil {
			return err
		}
	}
	return nil
}

// ActivateSchedule is public
func ActivateSchedule(sch cloudfunc.Schedule, usersstr []string, b *tgbotapi.BotAPI, ch chan SubEvent) {
	users := []int64{}
	for i := 0; i < len(usersstr); i++ {
		n, _ := strconv.ParseInt(usersstr[i], 10, 64)
		users = append(users, n)
	}
	for {
		del, ind := calcNextSchedule(sch)
		fmt.Println(users)
		fmt.Println("sleep for:", del.Minutes())
		time.Sleep(del)
		newInf := map[int64]bool{}

		for i := 0; i < len(users); i++ {
			newInf[users[i]] = true
		}

	Loop:
		for {
			select {
			case e := <-ch:
				switch e.Action {
				case Add:
					fmt.Println("adding user ", e.ChatID)
					newInf[e.ChatID] = true
				case Del:
					fmt.Println("deleting user ", e.ChatID)
					newInf[e.ChatID] = false
				}
			default:
				break Loop
			}
		}

		users = make([]int64, 0, len(newInf))

		for k, v := range newInf {
			if v {
				users = append(users, k)
			}
		}

		fmt.Println(users)
		SpreadMessage(b, users, sch.Event[ind])
		fbclient.SetSubscribers(users, sch.Name)
		fmt.Println("Success")
	}
}

func calcNextSchedule(s cloudfunc.Schedule) (time.Duration, int) {
	const mpw = 7 * 60 * 24
	now := time.Now().UTC().Add(3 * time.Hour)
	mins := cloudfunc.GetMinsOfWeek(now)
	nextEvent := 0
	minMins := mpw

	for i := 0; i < len(s.Event); i++ {
		curmins := (s.Minute[i] - 5 - mins + mpw) % mpw
		if minMins > curmins && curmins != 0 {
			nextEvent = i
			minMins = curmins
		}
	}
	return time.Duration(minMins) * time.Minute, nextEvent
}
