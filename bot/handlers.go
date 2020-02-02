package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/route"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/root"
)

func (b *Bot) handleMessage(update tgbotapi.Update) {
	if ok := b.selectCase(update.Message.Text, update); !ok {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, `Я не совсем понял, что ты хочешь, но могу предложить следующие функции:\n
		/today - показать пары на сегодня
		/morrow - показать пары на завтра
		/sub - подписаться
		/unsub - отписаться
		/week - показать пары на неделю`)
		if _, err := b.api.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func (b *Bot) handleCommand(update tgbotapi.Update) {
	b.selectCase(update.Message.Command(), update)
}

func (b *Bot) selectCase(text string, update tgbotapi.Update) bool {
	text = strings.ToLower(text)
	switch {
	case contains(text, "sub", "go", "подпис", "старт"):
		b.onSub(update)
	case contains(text, "reset", "unsub", "отпис"):
		b.onReset(update)
	case contains(text, "week", "недел", "распис"):
		b.onWeek(update)
	case contains(text, "today", "сегодня"):
		b.onToday(update)
	case contains(text, "morrow", "tomorrow", "завтра"):
		b.onMorrow(update)
	default:
		return false
	}
	return true
}

func contains(text string, words ...string) bool {
	for _, w := range words {
		if strings.Contains(text, w) {
			return true
		}
	}
	return false
}

func getSchForDay(sch cloudfunc.Schedule, day int) string {
	str := time.Weekday(day).String() + ":\n"
	_, thisWeek := time.Now().UTC().ISOWeek()
	for i := range sch.Event {
		if sch.Minute[i] >= day*24*60 && sch.Minute[i] < (day+1)*24*60 && (sch.Type[i] == -1 || sch.Type[i] == (thisWeek+1)%2) {
			str += sch.Event[i] + "\n"
		}
	}
	return str
}

func (b *Bot) onWeek(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	subs := fbclient.FetchUsersSubs(chatID)
	evlist := ""
	if len(subs) > 0 {
		schnameb, _ := json.Marshal(subs[0])
		schname := string(schnameb)

		for _, sch := range b.Table {
			if sch.Name == schname {
				for day := 1; day < 6; day++ {
					evlist += getSchForDay(sch, day) + "\n"
				}
				msg = tgbotapi.NewMessage(chatID, evlist)
			}
		}
	} else {
		msg = tgbotapi.NewMessage(chatID, "Для получения рассписания нужно выбрать группу")
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) onToday(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	subs := fbclient.FetchUsersSubs(chatID)
	if len(subs) > 0 {
		schnameb, _ := json.Marshal(subs[0])
		schname := string(schnameb)

		for _, sch := range b.Table {
			if sch.Name == schname {
				msg = tgbotapi.NewMessage(chatID, getSchForDay(sch, int(time.Now().Weekday())))
			}
		}
	} else {
		msg = tgbotapi.NewMessage(chatID, "Для получения рассписания нужно выбрать группу")
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) onMorrow(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")
	subs := fbclient.FetchUsersSubs(chatID)
	if len(subs) > 0 {
		schnameb, _ := json.Marshal(subs[0])
		schname := string(schnameb)

		for _, sch := range b.Table {
			if sch.Name == schname {
				msg = tgbotapi.NewMessage(chatID, getSchForDay(sch, int((time.Now().Weekday()+1)%7)))
			}
		}
	} else {
		msg = tgbotapi.NewMessage(chatID, "Для получения рассписания нужно выбрать группу")
	}
	if _, err := b.api.Send(msg); err != nil {
		log.Println(err)
	}
}

func (b *Bot) onSub(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	dropped := b.root.Rootnode.Drop()
	msg := tgbotapi.NewMessage(chatID, "Выбери свое расписание👇🏻\n"+dropped.String()) // ⬇️ 🎓 👇🏻
	mkp, ok := GenFor(dropped)
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
	if urt, ok := b.getResetTree(chatID, true); ok {
		mkp, _ := GenForReset(urt.Rootnode)
		msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("☠️🔥 %s", urt.Rootnode.Drop()))
		msg.ReplyMarkup = mkp
		fmt.Println("Doing reset for user", chatID)
		if _, err := b.api.Send(msg); err != nil {
			log.Println(err)
		}
	} else {
		b.noSubsMessage(chatID)
	}
}

func (b *Bot) noSubsMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Подписок нет 🙅🏿‍♂️")
	respmsg, err := b.api.Send(msg)
	if err != nil {
		log.Println(err)
	} else {
		b.sentresets[chatID] = respmsg.MessageID
	}
}

func (b *Bot) getResetTree(userID int64, forceUpd bool) (*route.TreeRoot, bool) {
	_, ok := b.resetTree[userID]
	if !ok || forceUpd {
		subs := fbclient.FetchUsersSubs(userID)
		var tr *route.TreeRoot
		if len(subs) > 0 {
			tr = route.NewTreeRoot(route.GenerateUsersTree(subs))
		}
		b.resetTree[userID] = tr
	}
	tr := b.resetTree[userID]
	if tr == nil {
		return nil, false
	}
	return tr, true
}

type idBundle struct {
	data       string
	chatID     int64
	messageID  int
	callbackID string
}

func (b *Bot) handleCallback(
	update tgbotapi.Update,
	chans map[string]chan root.SubEvent) {
	data := update.CallbackQuery.Data
	bundle := idBundle{
		data:       strings.Split(data, ":")[1],
		chatID:     update.CallbackQuery.Message.Chat.ID,
		messageID:  update.CallbackQuery.Message.MessageID,
		callbackID: update.CallbackQuery.ID,
	}
	fmt.Println(data)
	switch {
	case strings.Contains(data, "route"):
		b.onRoute(bundle, chans)
	case strings.Contains(data, "reset"):
		b.onResetCallback(bundle, chans)
	}
}

func (b *Bot) onResetCallback(bundle idBundle, chans map[string]chan root.SubEvent) {
	nodehash := bundle.data
	if urt, ok := b.getResetTree(bundle.chatID, false); ok {
		if node, ok := urt.Find(nodehash); ok {
			if mkp, ok := GenForReset(node); ok {
				msg := tgbotapi.NewEditMessageText(bundle.chatID, bundle.messageID, fmt.Sprintf("☠️🔥 %s", node.Drop()))
				msg.ReplyMarkup = &mkp
				if bundle.callbackID != "" {
					if _, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(bundle.callbackID, "")); err != nil {
						log.Println(err)
					}
				}
				respmsg, err := b.api.Send(msg)
				if err != nil {
					log.Println(err)
				} else {
					b.sentresets[bundle.chatID] = respmsg.MessageID
				}
			} else {
				scheduleName := node.MakePath()
				ch, ok := chans[scheduleName]
				if ok {
					fmt.Println(bundle.data)
					ch <- root.SubEvent{Action: root.Del, SubID: bundle.chatID}
					fbclient.DeleteSubscriber(bundle.chatID, scheduleName)
					// snackMsg := "Un️subscribed ♻️" // ☠️
					// snackMsg := "Отписка проведена ♻️ (" + cmdMapping[data] + ")"
					// snackMsg := "Отписка проведена ♻️" + node.Name
					snackMsg := "📆 " + node.Name + " - отписка проведена♻️"
					b.api.AnswerCallbackQuery(tgbotapi.NewCallback(bundle.callbackID, snackMsg))
					msg := tgbotapi.NewMessage(bundle.chatID, snackMsg)
					if _, err := b.api.Send(msg); err != nil {
						log.Println(err)
					}
					b.getResetTree(bundle.chatID, true)
					b.onResetCallback(idBundle{
						data:      node.Jump().CalcHash64(),
						chatID:    bundle.chatID,
						messageID: bundle.messageID}, chans)
				}
			}
		}
	} else {
		delcfg := tgbotapi.NewDeleteMessage(bundle.chatID, bundle.messageID)
		if _, err := b.api.DeleteMessage(delcfg); err != nil {
			log.Println(err)
		}
		b.noSubsMessage(bundle.chatID)
	}
}

func (b *Bot) onRoute(bundle idBundle, chans map[string]chan root.SubEvent) {
	nodehash := bundle.data
	if node, ok := b.root.Find(nodehash); ok {
		node = node.Drop()
		if mkp, ok := GenFor(node); ok {
			msg := tgbotapi.NewEditMessageText(bundle.chatID, bundle.messageID, fmt.Sprintf("%s👇🏻", node))
			msg.ReplyMarkup = &mkp
			if _, err := b.api.AnswerCallbackQuery(tgbotapi.NewCallback(bundle.callbackID, "")); err != nil {
				log.Println(err)
			}
			if _, err := b.api.Send(msg); err != nil {
				log.Println(err)
			}
		} else {
			schName := node.MakePath()
			ch, ok := chans[schName]
			if ok {
				subs := fbclient.FetchUsersSubs(bundle.chatID)

				for _, path := range subs {
					schnameb, _ := json.Marshal(path)
					schname := string(schnameb)
					b.updsmap[schname] <- root.SubEvent{Action: root.Del, SubID: bundle.chatID}
					fbclient.DeleteSubscriber(bundle.chatID, schname)
				}

				ch <- root.SubEvent{Action: root.Add, SubID: bundle.chatID}
				fbclient.AddSubscriber(bundle.chatID, schName)
				// snackMsg := "Our congrats 🥂. We handled your sub!"
				// snackMsg := "Ваша регистрация обработана 🥂 (" + cmdMapping[data] + ")"
				// snackMsg := "Поздравляю🥂. Подписочка \"" + node.Name + "\" подписана."
				// snackMsg := "Еба захендил🥂 " + node.Name + " +1"
				// snackMsg := "📆 " + node.Name + " - подписка подписана🥂"
				snackMsg := "📆 " + node.Name + " - поздравляем, вы подписались🥂"
				b.api.AnswerCallbackQuery(tgbotapi.NewCallback(bundle.callbackID, snackMsg))
				msg := tgbotapi.NewMessage(bundle.chatID, snackMsg)
				if _, err := b.api.Send(msg); err != nil {
					log.Println(err)
				}
				if resetID, ok := b.sentresets[bundle.chatID]; ok {
					rt, _ := b.getResetTree(bundle.chatID, true)
					b.onResetCallback(idBundle{
						data:      rt.Rootnode.CalcHash64(),
						chatID:    bundle.chatID,
						messageID: resetID}, chans)
				}
			} else {
				snackMsg := "Этот шедуль уже не в базе"
				b.api.AnswerCallbackQuery(tgbotapi.NewCallback(bundle.callbackID, snackMsg))
				log.Printf("Schedule was not found by name: %s\n", schName)
			}
		}
	}
}
