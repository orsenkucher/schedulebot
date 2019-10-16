package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/dev"
	"github.com/orsenkucher/schedulebot/fbclient"
)

// *** ASAP ***
// Reliable subscription!
// Migalky
//
// *** Proposals ***
//Generate buttoms by path like below
//Ukraine?.Mehmat.firstyear.math.group1.subgroup2
//Custom schedules
//Create kostyl for migalki
//

func main() {
	// fbclient.CreateDemoSched()

	///*
	dev.GenerateTestSchedule()
	fmt.Println("Minuted from week start", cloudfunc.GetMinsOfWeek(time.Now()))
	table := fbclient.FetchTable()
	users := fbclient.FetchSubscribers()
	// fmt.Println(table)

	// key, err := creds.ReadToken()
	// if err != nil {
	// 	panic(err)
	// }
	key := creds.ReadCr()

	// fmt.Println(key)
	b := bot.InitBot(key)
	chans := map[string]chan bot.SubEvent{}

	for _, sch := range table {
		chans[sch.Name] = make(chan bot.SubEvent)
		go bot.ActivateSchedule(sch, users[sch.Name].IDs, b, chans[sch.Name])
	}
	bot.Listen(b, chans)
	//*/
	//fbclient.CreateSchedule()
}
