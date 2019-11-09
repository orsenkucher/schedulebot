package sch

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/root"
)

// Scheduler schedules send message jobs
type Scheduler struct {
	bot  *bot.Bot // replace with chan!
	upds chan root.SubEvent
	subs map[int64]void
	sch  cloudfunc.Schedule
}

const mpw = root.MPW

type void struct{}

var voi void

// SpawnSchedulers spawns and activates all schedulers
func SpawnSchedulers(bot *bot.Bot) map[string]chan root.SubEvent {
	table := fbclient.FetchSchedules()
	subscribers := fbclient.FetchSubscribers()
	chMap := make(map[string]chan root.SubEvent)
	for _, sch := range table {
		ss := parseSubscribers(subscribers[sch.Name])
		ch := newScheduler(bot, sch, ss)
		chMap[sch.Name] = ch
	}
	return chMap
}

func parseSubscribers(subs []cloudfunc.Subscriber) map[int64]void {
	ids := make(map[int64]void)
	for i := 0; i < len(subs); i++ {
		id, _ := strconv.ParseInt(subs[i].ID, 10, 64)
		ids[id] = voi
	}
	return ids
}

func newScheduler(bot *bot.Bot, sch cloudfunc.Schedule, ss map[int64]void) chan root.SubEvent {
	ch := make(chan root.SubEvent)
	s := Scheduler{
		upds: ch,
		bot:  bot,
		sch:  sch,
		subs: ss}
	go s.activateSchedule()
	go s.listenSubEvents()
	return ch
}

func (s *Scheduler) listenSubEvents() {
	for {
		select {
		case e := <-s.upds:
			switch e.Action {
			case root.Add:
				fmt.Println("adding user ", e.SubID)
				s.subs[e.SubID] = voi
			case root.Del:
				fmt.Println("deleting user ", e.SubID)
				delete(s.subs, e.SubID)
			}
		}
	}
}

func (s *Scheduler) getSubIDs() []int64 {
	ids := make([]int64, 0, len(s.subs))
	for k := range s.subs {
		ids = append(ids, k)
	}
	return ids
}

func (s *Scheduler) activateSchedule() {
	for {
		del, ind := calcNextSchedule(s.sch)
		fmt.Println(s.getSubIDs())
		fmt.Println("sleep for:", del.Minutes())
		time.Sleep(del)

		ids := s.getSubIDs()
		fmt.Println(ids)
		s.bot.SpreadMessage(ids, s.sch.Event[ind])
		fmt.Println("Success")
	}
}

func calcNextSchedule(s cloudfunc.Schedule) (time.Duration, int) {
	now := time.Now().UTC()
	mins := cloudfunc.GetMinsOfWeek(now)
	nextEvent := 0
	minMins := 2*mpw + 1
	_, thisWeek := time.Now().UTC().ISOWeek()
	thisWeek %= 2

	for i := 0; i < len(s.Event); i++ {
		curmins := (s.Minute[i] - 5 - mins + mpw) % mpw
		if s.Type[i] == (thisWeek+1)%2 {
			curmins += mpw
		}
		if minMins > curmins && curmins != 0 {
			nextEvent = i
			minMins = curmins
		}
	}
	return time.Duration(minMins) * time.Minute, nextEvent
}
