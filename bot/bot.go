package bot

import (
	"log"
	"time"

	"github.com/orsenkucher/schedulebot/route"

	"github.com/orsenkucher/schedulebot/sch"

	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/root"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Bot is a scheduler bot
type Bot struct {
	credential creds.Credential
	api        *tgbotapi.BotAPI
	Jobs       chan sch.Job
	root       *route.TreeRoot
	sentmap    map[int64]int // map[userID]msgID
}

// NewBot creates new scheduler bot with provided credentials
func NewBot(cr creds.Credential, root *route.TreeRoot) *Bot {
	b := &Bot{
		credential: cr,
		Jobs:       make(chan sch.Job),
		root:       root,
		sentmap:    make(map[int64]int)}
	b.initAPI()
	go b.processJobs()
	return b
}

func (b *Bot) initAPI() {
	var err error
	b.api, err = tgbotapi.NewBotAPI(b.credential.String())
	if err != nil {
		log.Panic(err)
	}

	b.api.Debug = false
	log.Printf("Authorized on account %s", b.api.Self.UserName)

	_, err = b.api.RemoveWebhook()
	if err != nil {
		log.Println("Cant remove webhook")
	}
}

func (b *Bot) processJobs() {
	for {
		select {
		case j := <-b.Jobs:
			b.SpreadMessage(j.Subs, j.Event)
		}
	}
}

// Listen starts infinite listening
func (b *Bot) Listen(updsmap map[string]chan root.SubEvent) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleCallback(update, updsmap)
			continue
		}

		if update.Message != nil {
			if update.Message.IsCommand() {
				b.handleCommand(update)
			} else if update.Message.Text != "" {
				b.handleMessage(update)
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			continue
		}
	}
}

// SpreadMessage is used to spread message to users
func (b *Bot) SpreadMessage(users []int64, msg string) {
	log.Printf("Sending message to %v users\n", len(users))
	for _, u := range users {
		time.Sleep(100 * time.Millisecond)

		log.Printf("Deleting previous msg for %v\n", u)
		if lastsent, ok := b.sentmap[u]; ok {
			delcfg := tgbotapi.NewDeleteMessage(u, lastsent)
			if _, err := b.api.DeleteMessage(delcfg); err != nil {
				log.Println(err)
			}
		}

		log.Printf("Sending to %v\n", u)
		tgmsg := tgbotapi.NewMessage(u, msg)
		sent, err := b.api.Send(tgmsg)
		if err != nil {
			log.Println(err)
		} else {
			b.sentmap[u] = sent.MessageID
		}
	}
}
