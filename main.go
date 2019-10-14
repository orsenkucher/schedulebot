package main

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/fbclient"
)

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
	chans := map[string]chan int64{}

	for _, sch := range table {
		chans[sch.Name] = make(chan int64)
		go send(sch, users[sch.Name].IDs, b, chans[sch.Name])
	}
	bot.Listen(b, chans)
	//*/
	//fbclient.CreateSchedule()
}

func send(sch cloudfunc.Schedule, usersstr []string, b *tgbotapi.BotAPI, ch chan int64) {
	users := []int64{}
	for i := 0; i < len(usersstr); i++ {
		n, _ := strconv.ParseInt(usersstr[i], 10, 64)
		users = append(users, n)
	}
	for {
		del, ind := calc(sch)
		fmt.Println(users)
		fmt.Println("sleep for:", del.Minutes())
		time.Sleep(del)
		newUsers := []int64{}
	Loop:
		for {
			select {
			case i := <-ch:
				newUsers = append(newUsers, i)
			default:
				break Loop
			}
		}
		users = append(users, newUsers...)

		fmt.Println(users)
		bot.SpreadMessage(b, users, sch.Event[ind])
		fbclient.AddSubscribers(newUsers, sch.Name)
		fmt.Println("Success")
	}
}

func calc(s cloudfunc.Schedule) (time.Duration, int) {
	const mpw = 7 * 60 * 24
	now := time.Now().UTC().Add(3 * time.Hour)
	mins := cloudfunc.GetMinsOfWeek(now)
	nextEvent := 0
	minMins := mpw

	for i := 0; i < len(s.Event); i++ {
		curmins := (s.Minute[i] - 5 - mins + mpw) % mpw
		if minMins > curmins && curmins != 0 {
			nextEvent = i
			minMins = curmins
		}
	}
	return time.Duration(minMins) * time.Minute, nextEvent
}
