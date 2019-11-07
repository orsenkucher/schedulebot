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

// Scheduler schedules send message jobs
type Scheduler struct {
	bot         *bot.Bot // replace with chan!
	SubStream   map[string]chan subs.SubEvent
	table       fbclient.Table
	subscribers map[string][]cloudfunc.Subscriber
}

// NewScheduler is used to activate Scheduler
func NewScheduler(bot *bot.Bot, subChan map[string]chan subs.SubEvent) *Scheduler {
	table := fbclient.FetchTable()
	subscribers := fbclient.FetchSubscribers()
	s := Scheduler{
		bot:         bot,
		SubStream:   subChan,
		table:       table,
		subscribers: subscribers}

	s.activateSchedules()
	return &s
}

func (s *Scheduler) activateSchedules() {
	for _, sch := range s.table {
		s.SubStream[sch.Name] = make(chan subs.SubEvent)
		go s.activateSchedule(sch, s.SubStream[sch.Name])
	}
}

// listen subChan here and mutate scheduler users field
// maybe we need scheduler per schedule (not one, but multiple ones)

// func (s *Scheduler) listenSubEvents(ch chan subs.SubEvent) {

// }

func (s *Scheduler) activateSchedule(sch cloudfunc.Schedule, ch chan subs.SubEvent) {
	users := []int64{}
	sub := s.subscribers[sch.Name]
	for i := 0; i < len(sub); i++ {
		n, _ := strconv.ParseInt(sub[i].ID, 10, 64)
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
		s.bot.SpreadMessage(users, sch.Event[ind])
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
