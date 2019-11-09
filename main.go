package main

import (
	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/sch"
)

// *** ASAP ***
// [+] Migalky (spin[up/down])
// [+] Append sch.json with schs for Thu and Fri
//
// *** Current ***
// [.] Generate schedule path from direcory it lies in
// [.] Generate buttons by path like below
//     Ukraine?.Mehmat.firstyear.math.group1.subgroup2
//
// *** Proposals ***
// [.] Use hash to determine whether sch update is needed
// [.] Custom schedules
//
func main() {
	fbclient.GenerateTestSchedule()
	// fbclient.CreateSchFromJSON()

	// /*
	b := bot.NewBot(creds.Cr459)
	updsmap := sch.SpawnSchedulers(b.Jobs)
	b.Listen(updsmap)
	//*/

	//fbclient.CreateSchedule()
}
