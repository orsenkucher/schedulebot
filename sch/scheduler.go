package sch

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/root"
)

// Scheduler schedules send message jobs
type Scheduler struct {
	jobs chan Job
	upds chan root.SubEvent
	subs map[int64]void
	sch  cloudfunc.Schedule
}

// Job is fired out of scheduler for someone to perform
type Job struct {
	Subs  []int64
	Event string
}

const mpw = root.MPW

type void struct{}

var voi void

func newScheduler(jobs chan Job, sch cloudfunc.Schedule, ss map[int64]void) chan root.SubEvent {
	ch := make(chan root.SubEvent)
	s := Scheduler{
		upds: ch,
		jobs: jobs,
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
		delay, idx := calcNextSchedule(s.sch)
		fmt.Println(s.getSubIDs())
		fmt.Println(s.sch.Name, " sleep for:", delay.Minutes())
		// if delay < 3*time.Hour {
		// 	ids := s.getSubIDs()
		// 	s.jobs <- Job{Subs: ids, Event: s.sch.Event[idx]}
		// }
		time.Sleep(delay)

		ids := s.getSubIDs()
		fmt.Println(ids)
		delay2, idx2 := calcNextSchedule(s.sch)
		// s.bot.SpreadMessage(ids, s.sch.Event[ind])
		msg := s.sch.Event[idx]
		if delay2 < 3*time.Hour {
			msg += "\n\n" + s.sch.Event[idx2]
		}
		s.jobs <- Job{Subs: ids, Event: msg}
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
		curmins := (s.Minute[i] - 15 - mins + mpw) % mpw
		//if s.Type[i] != -1 && (s.Type[i] == (thisWeek+1)%2) == (mins < s.Minute[i]-15) {
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
