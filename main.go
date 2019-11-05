package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/creds"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/fbclient"
)

// *** ASAP ***
// [.] Migalky (spin[up/down]) // for now just delete all spin="up" from sch.json 20.10.19 is "down" week
// [+] Append sch.json with schs for Thu and Fri
//
// *** Proposals ***
// [.] Use hash to determine whether sch update is needed
// [.] Generate buttons by path like below
//     Ukraine?.Mehmat.firstyear.math.group1.subgroup2
// [.] Custom schedules
//
// ***
// В основе лежит карта map[int64]User
// тип User глобальный и находится в model.User
// у Юзера есть ID
// у Юзера есть его текущий путь route routes.routeTree
// у Юзера есть канал для общения с шедулером (возможно)
//
// Как только в луп бота приходит новое событие он первым делом находит в карте юзера
//
// У бота есть канал из сообщений для отправки, но его читает 30 раз в сек (возможно)
//
func main() {
	// fbclient.CreateSchFromJSON()

	// /*
	// fbclient.GenerateTestSchedule()
	fmt.Println("Minutes from week start", cloudfunc.GetMinsOfWeek(time.Now()))
	table := fbclient.FetchTable()
	users := fbclient.FetchSubscribers()

	b := bot.NewBot(creds.Cr459)

	chans := map[string]chan bot.SubEvent{}

	for _, sch := range table {
		chans[sch.Name] = make(chan bot.SubEvent)
		go b.ActivateSchedule(sch, users[sch.Name], chans[sch.Name])
	}
	b.Listen(chans)
	//*/
	//fbclient.CreateSchedule()
}
