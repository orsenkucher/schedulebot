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
	// ch := t
	// fn := func(_, name string) { ch = ch.MakeChild(name) }
	fn := func(path, name string) route.MyFn {
		ch := t.MakeChild(name)
		var fn2 func(path, name string) route.MyFn
		fn2 = func(path, name string) route.MyFn {
			ch = ch.MakeChild(name)
			return fn2
		}
		return fn2
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

func Kek(path, name string) route.MyFn {
	return nil
}
