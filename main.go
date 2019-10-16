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
//Respond to callbackQuery !
//
// *** Proposals ***
//Generate buttoms by path lik below
//Mehmat.firstyear.math.group1.subgroup2
//Custom schedules
//create kostyl for migalki
//

func main() {
	// fbclient.CreateDemoSched()

	///*
	GenerateTestSchedule()
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

//GenerateTestSchedule is test
func GenerateTestSchedule() {
	mins := cloudfunc.GetMinsOfWeek(time.Now())
	schedule := cloudfunc.Schedule{
		Name:   "test",
		Event:  []string{"We started", "Still alive", "Unbelivable", "5 минут, полёт нормальный", "I`ll send you one more in one min if everything is good"},
		Minute: []int{mins + 6, mins + 7, mins + 8, mins + 9, mins + 10},
	}
	fbclient.SendSchedule(schedule)
	fbclient.SetSubscribers([]int64{259224772, 364448153}, "test")
}
