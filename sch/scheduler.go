package sch

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/subs"
)

type Scheduler struct {
	bot   *bot.Bot // replace with chan!
	chans map[string]chan subs.SubEvent
}

func NewScheduler(bot *bot.Bot, chans map[string]chan subs.SubEvent) Scheduler {
	return Scheduler{bot: bot, chans: chans}
}

func (s *Scheduler) ActivateSchedules(table fbclient.Table, subscribers map[string][]cloudfunc.Subscriber) {
	for _, sch := range table {
		s.chans[sch.Name] = make(chan subs.SubEvent)
		go s.ActivateSchedule(sch, subscribers[sch.Name], s.bot, s.chans[sch.Name])
	}
}

// ActivateSchedule is public
func (*Scheduler) ActivateSchedule(sch cloudfunc.Schedule, usersstr []cloudfunc.Subscriber, b *bot.Bot, ch chan subs.SubEvent) {
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
				case subs.Add:
					fmt.Println("adding user ", e.ChatID)
					newInf[e.ChatID] = true
				case subs.Del:
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
	minMins := 2*MPW + 1
	_, thisWeek := time.Now().UTC().ISOWeek()
	thisWeek %= 2

	for i := 0; i < len(s.Event); i++ {
		curmins := (s.Minute[i] - 5 - mins + MPW) % MPW
		if s.Type[i] == (thisWeek+1)%2 {
			curmins += MPW
		}
		if minMins > curmins && curmins != 0 {
			nextEvent = i
			minMins = curmins
		}
	}
	return time.Duration(minMins) * time.Minute, nextEvent
}
