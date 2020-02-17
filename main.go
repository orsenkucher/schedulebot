package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/route"
)

// *** ASAP ***
// [+] Migalky (spin[up/down])
// [+] Append sch.json with schs for Thu and Fri
//
// *** Current ***
// [.] Generate schedule path from direcory it lies in
// [+] Generate buttons by path like below
//     Ukraine?.Mehmat.firstyear.math.group1.subgroup2
// [.] Finish schmanager cmd
//
// *** Proposals ***
// [.] Use hash to determine whether sch update is needed
// [.] Custom schedules
//
// *** v2.0 ***
// [.] /start message
// [.] day/week events with "remove" button
// [.] on sub edit start msg
// [.] on sub unsub from previous sub
// [.] cmds: [sub unsub week day]
// [.] возможность для правки страростами /edit cmd
func main() {
	fmt.Println(cloudfunc.GetMinsOfWeek(time.Now().UTC()))

	//fbclient.GenerateTestSchedule()
	// fbclient.CreateSchFromJSON(root.SchFile)

	// /*
	//fbclient.CreateSchFromOS()

	t := route.BuildOSTree()
	t.Print()
	tr := route.NewTreeRoot(t)
	b := bot.NewBot(creds.Cr459, tr)
	b.SendToEveryone("q")
	//updsmap, table := sch.SpawnSchedulers(b.Jobs)

	//b.Table = table
	//b.Listen(updsmap)
	// */

	//fbclient.CreateSchedule()
}
