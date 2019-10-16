package dev

import (
	"time"

	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/fbclient"
)

//GenerateTestSchedule is test
func GenerateTestSchedule() {
	mins := cloudfunc.GetMinsOfWeek(time.Now())
	schedule := cloudfunc.Schedule{
		Name:   "test",
		Event:  []string{"We started", "Still alive", "Unbelivable", "5 минут, полёт нормальный", "I`ll send you one more in one min if everything is good"},
		Minute: []int{mins + 6, mins + 7, mins + 8, mins + 9, mins + 10},
	}
	fbclient.SendSchedule(schedule)
	fbclient.SetSubscribers([]int64{259224772, 364448153}, "test")
}
