package main

import (
	"fmt"
	"strconv"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/fbclient"
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
	//fmt.Println(cloudfunc.GetMinsOfWeek(time.Now().UTC()))

	//fbclient.GenerateTestSchedule()
	//fbclient.CreateSchFromJSON(root.SchFile)

	// /*
	//fbclient.CreateSchFromOS()

	t := route.BuildOSTree()
	t.Print()
	tr := route.NewTreeRoot(t)
	b := bot.NewBot(creds.Cr459, tr)

	subs := fbclient.FetchSubscribers()
	allsubs := make(map[string]int)

	for _, v := range subs {
		for _, sub := range v {
			allsubs[sub.ID] = 0
		}
	}

	sublist := []int64{}

	for k := range allsubs {
		i, _ := strconv.ParseInt(k, 10, 0)
		//if i == 364448153 || i == 259224772 {
		sublist = append(sublist, i)
		//}
	}

	fmt.Print(sublist)
	fmt.Scanln()

	b.SpreadMessage(sublist, "Дорогой пользователь!\nРад сообщить, что меня обновили и теперь я умею делать ещё и так:\n/today - показать пары на сегодня\n/morrow - показать пары на завтра\n/sub - подписаться\n/unsub - отписаться\n/week - показать пары на неделю")
	//updsmap, table := sch.SpawnSchedulers(b.Jobs)
	//b.Table = table
	//b.Listen(updsmap)
	// */

	//fbclient.CreateSchedule()
}
