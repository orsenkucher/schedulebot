package main

import "github.com/orsenkucher/schedulebot/fbclient"

func main() {
	/*
		key, err := creds.ReadToken()

		if err != nil {
			panic(err)
		}

		fmt.Println(key)
		b := bot.InitBot(key) //*/
	/*
		j := func() { job.Job(b) }
		doneCh := scheduler.ScheduleJob(j, 10*time.Second)
		defer func() { doneCh <- struct{}{} }()

		bot.Listen(b)//*/
	fbclient.CreateSchedule()
}
