package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/cloudfunc"

	"github.com/orsenkucher/schedulebot/fbclient"
)

func main() {
	table := fbclient.FetchTable()
	fmt.Println(table)

	for _, sch := range table {
		dur, event := calc(sch)
		time.AfterFunc(dur, func() { f(event) })
	}
	time.Sleep(30 * time.Second)

	// 	key, err := creds.ReadToken()

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(key)
	// 	b := bot.InitBot(key)

	// 	bot.Listen(b)
}

func f(ind int) {
	fmt.Println(ind)
}

func calc(s cloudfunc.Schedule) (time.Duration, int) {
	const mpw = 7 * 60 * 24
	now := time.Now().UTC().Add(3 * time.Hour)
	mins := cloudfunc.GetMinsOfWeek(now)
	nextEvent := 0
	minMins := (s.Minute[0] - mins + mpw) % mpw

	for i := 1; i < len(s.Event); i++ {
		if minMins > (s.Minute[i]-mins+mpw)%mpw {
			nextEvent = i
			minMins = (s.Minute[i] - mins + mpw) % mpw
		}
	}
	return time.Duration(minMins) * time.Minute, nextEvent
}
