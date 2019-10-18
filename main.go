package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/fbclient"
)

// *** ASAP ***
// Fix time on server
// Reliable subscription!
// Migalky (spin[up/down])
//
// *** Proposals ***
// Use hash to determine whether sch update is needed
// Generate buttoms by path like below
// Ukraine?.Mehmat.firstyear.math.group1.subgroup2
// Custom schedules
//

func main() {
	fbclient.CreateSchFromJSON()

	// /*
	fbclient.GenerateTestSchedule()
	fmt.Println("Minuted from week start", cloudfunc.GetMinsOfWeek(time.Now()))
	table := fbclient.FetchTable()
	users := fbclient.FetchSubscribers()

	key := creds.ReadCr()
	b := bot.InitBot(key)

	chans := map[string]chan bot.SubEvent{}

	for _, sch := range table {
		chans[sch.Name] = make(chan bot.SubEvent)
		go bot.ActivateSchedule(sch, users[sch.Name], b, chans[sch.Name])
	}
	bot.Listen(b, chans)
	//*/
	//fbclient.CreateSchedule()
}
