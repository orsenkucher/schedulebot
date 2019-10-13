package main

import (
	"fmt"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/job"
	"github.com/orsenkucher/schedulebot/scheduler"
)

func main() {

	// /*
	key, err := creds.ReadToken()
	if err != nil {
		panic(err)
	}

	fmt.Println(key)
	b := bot.InitBot(key) //*/

	j := func() { job.Job(b) }
	doneCh := scheduler.ScheduleJob(j, 10*time.Second)
	defer func() { doneCh <- struct{}{} }()

	bot.Listen(b)
}
