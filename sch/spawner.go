package sch

import (
	"strconv"

	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/root"
)

// SpawnSchedulers spawns and activates all schedulers
func SpawnSchedulers(jobs chan Job) map[string]chan root.SubEvent {
	table := fbclient.FetchSchedules()
	subscribers := fbclient.FetchSubscribers()
	chMap := make(map[string]chan root.SubEvent)
	for _, sch := range table {
		ss := parseSubscribers(subscribers[sch.Name])
		ch := newScheduler(jobs, sch, ss)
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
