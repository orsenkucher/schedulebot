package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/cloudfunc"

	"github.com/orsenkucher/schedulebot/fbclient"
)

func main() {
	table := fbclient.FetchTable()
	fmt.Println(table)

	for _, sch := range table {
		time.AfterFunc(calc(sch), func() {})
	}

	// 	key, err := creds.ReadToken()

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(key)
	// 	b := bot.InitBot(key)

	// 	bot.Listen(b)
}

func calc(s cloudfunc.Schedule) time.Duration {
	now := time.Now().UTC().Add(3 * time.Hour)
	mins := cloudfunc.GetMinsOfWeek(now)
	nextEvent := 0

	for i := 0; i < len(s.Event); i++ {

	}
	return 0
}
