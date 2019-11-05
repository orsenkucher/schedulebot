package bot

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/orsenkucher/schedulebot/creds"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/cloudfunc"
)

// Bot is a scheduler bot
type Bot struct {
	credential creds.Credential
	api        *tgbotapi.BotAPI
}

// NewBot creates new scheduler bot with provided credentials
func NewBot(cr creds.Credential) *Bot {
	b := &Bot{credential: cr}
	b.initAPI()
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

// Listen starts infinite listening
func (b *Bot) Listen(chans map[string]chan SubEvent) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(b.api, update, chans)
			continue
		}

		if update.Message != nil {
			if update.Message.IsCommand() {
				handleCommand(b.api, update)
			} else if update.Message.Text != "" {
				handleMessage(b.api, update)
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			continue
		}
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

func sendOnChan(ch chan SubEvent, e SubEvent) {
	ch <- e
}

// SpreadMessage is public
func (b *Bot) SpreadMessage(users []int64, msg string) {
	log.Printf("Sending message to %v users\n", len(users))
	for _, u := range users {
		time.Sleep(100 * time.Millisecond)
		tgmsg := tgbotapi.NewMessage(u, msg)
		log.Printf("Sending to %v\n", u)
		if _, err := b.api.Send(tgmsg); err != nil {
			log.Println(err)
		}
	}
}

// ActivateSchedule is public
func (b *Bot) ActivateSchedule(sch cloudfunc.Schedule, usersstr []cloudfunc.Subscriber, ch chan SubEvent) {
	users := []int64{}
	for i := 0; i < len(usersstr); i++ {
		n, _ := strconv.ParseInt(usersstr[i].ID, 10, 64)
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
		b.SpreadMessage(users, sch.Event[ind])
		fmt.Println("Success")
	}
}

// MPW is total minutes in week
const MPW = 7 * 60 * 24

func calcNextSchedule(s cloudfunc.Schedule) (time.Duration, int) {
	now := time.Now().UTC()
	mins := cloudfunc.GetMinsOfWeek(now)
	nextEvent := 0
	minMins := MPW

	for i := 0; i < len(s.Event); i++ {
		curmins := (s.Minute[i] - 5 - mins + MPW) % MPW
		if minMins > curmins && curmins != 0 {
			nextEvent = i
			minMins = curmins
		}
	}
	return time.Duration(minMins) * time.Minute, nextEvent
}
