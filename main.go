package main

import (
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/sch"

	"github.com/orsenkucher/schedulebot/bot"
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
	fbclient.GenerateTestSchedule()
	// fbclient.CreateSchFromJSON()

	// /*

	// schedStream := map[string]chan subs.SubEvent{}

	b := bot.NewBot(creds.Cr459)
	subStream := sch.SpawnSchedulers(b)
	b.Listen(subStream)

	//*/

	//fbclient.CreateSchedule()
}
