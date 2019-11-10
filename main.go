package main

import (
	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/root"
	"github.com/orsenkucher/schedulebot/route"
	"github.com/orsenkucher/schedulebot/sch"
)

// *** ASAP ***
// [+] Migalky (spin[up/down])
// [+] Append sch.json with schs for Thu and Fri
//
// *** Current ***
// [.] Generate schedule path from direcory it lies in
// [+] Generate buttons by path like below
//     Ukraine?.Mehmat.firstyear.math.group1.subgroup2
//
// *** Proposals ***
// [.] Use hash to determine whether sch update is needed
// [.] Custom schedules
//
func main() {
	// fbclient.GenerateTestSchedule()
	// fbclient.CreateSchFromJSON(root.SchFile)

	// /*
	var lc route.TreeCreator = route.LocalCreator{Root: root.Rootdir}
	t := &route.Tree{Name: "root"}
	fn := func(path, name string) route.MyFn {
		ch := t.MakeChild(name)
		var spawnfn func(tr *route.Tree) route.MyFn
		spawnfn = func(tr *route.Tree) route.MyFn {
			return func(path, name string) route.MyFn {
				ch2 := tr.MakeChild(name)
				return spawnfn(ch2)
			}
		}

		return spawnfn(ch)
	}
	lc.Create(fn)
	t.Print()
	tr := route.NewTreeRoot(t)
	b := bot.NewBot(creds.Cr459, tr)
	updsmap := sch.SpawnSchedulers(b.Jobs)
	b.Listen(updsmap)
	//*/

	//fbclient.CreateSchedule()
}
