package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/fbclient"
)

//Generate buttoms by path lik below
//Mehmat.firstyear.math.group1.subgroup2
//Custom schedules
//create kostyl for migalki

func main() {
	// fbclient.CreateDemoSched()
	//*
	fmt.Println("Minuted from week start", cloudfunc.GetMinsOfWeek(time.Now()))
	table := fbclient.FetchTable()
	users := fbclient.FetchSubscribers()
	// fmt.Println(table)

	key, err := creds.ReadToken()

	if err != nil {
		panic(err)
	}

	// fmt.Println(key)
	b := bot.InitBot(key)
	chans := map[string]bot.SubChans{}

	for _, sch := range table {
		chans[sch.Name] = bot.SubChans{AddChan: make(chan int64), DelChan: make(chan int64)}
		go bot.ActivateSchedule(sch, users[sch.Name].IDs, b, chans[sch.Name])
	}
	bot.Listen(b, chans)
	//*/
	//fbclient.CreateSchedule()
}
