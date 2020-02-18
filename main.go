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
	msg := `Привет! 👋
Мы, разработчики *459бота* работаем над новым проектом — текстовой игрой-головоломкой.
	
Вот вопросы, которые мы хотелы бы видеть в нашей игре.
Попробуйте решить 😉

*1.* If X=24 and P=16, what is X-P-A?
*2.* What do you call a FISH with no Eyes?
*3.* A K Q J ?
*4.* What is the end of everything?
*5.* March 14

Нам нужно *очень*, *очень* много вопросов и без вашей помощи мы не справимся. Потому конкурс.
Условия конкурса:
1️⃣_Придумать 5 вопросов с уникальной логикой._
2️⃣_Отправить их_ @sergeycheremshinsky


_Авторы лучших вопросов получат вознаграждение в размере:_
1-е место: *500 грн*
2-е место: *250 грн*
3-е место: *150 грн*
❤️ *Хотя бы 1 хороший вопрос — упоминание в титрах приложения!* ❤️

❗️Итоги *1 марта*

Ответы и разъяснения на примеры вопросов по ссылке: 
@NP2Ans

Доп. информация: @vvvvvvv9sss`
	b.SendToEveryone(msg)
	//updsmap, table := sch.SpawnSchedulers(b.Jobs)

	//b.Table = table
	//b.Listen(updsmap)
	// */

	//fbclient.CreateSchedule()
}
