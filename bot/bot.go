package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// InitBot initializes telegram bot
func InitBot(withKey string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(withKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}

// Listen starts infinite listening
func Listen(bot *tgbotapi.BotAPI, chans map[string]chan int64) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(bot, update, chans)
			continue
		}

		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(bot, update)
			} else if update.Message.Text != "" {
				handleMessage(bot, update)
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			continue
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

var inlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" 🥞 1 група  ", "group1"),
		tgbotapi.NewInlineKeyboardButtonData(" 🍇 2 група  ", "group2"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(" 🤹 demo  ", "demo"),
	),
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
	chans map[string]chan int64) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID
	ch, ok := chans[data]
	if ok {
		fmt.Println(data)
		go addNewUser(ch, chatID)
	}
}

func addNewUser(ch chan int64, user int64) {
	ch <- user
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
