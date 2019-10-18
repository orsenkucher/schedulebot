package fbclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

//DayIndex is indexes
var DayIndex = map[string]int{
	"Sun":  0,
	"Mon":  1,
	"Tue":  2,
	"Wed":  3,
	"Thu":  4,
	"Fri":  5,
	"Sat":  6,
	"STOP": -1}

// CreateSchedule is public
func CreateSchedule() {
	fmt.Println("Creating Schedule")
	fmt.Print("Name: ")
	var name string
	fmt.Scanln(&name)
	fmt.Println("Now, let's create events")
	schedule := cloudfunc.Schedule{
		Name:   name,
		Event:  []string{},
		Minute: []int{},
	}

	for {
		fmt.Print("WeekDay or 'STOP' if you want to finish: ")
		var day string
		fmt.Scanln(&day)
		wday, exist := DayIndex[day]
		if !exist {
			continue
		}
		if wday == -1 {
			break
		}
		var str string
		fmt.Print("Hour: ")
		fmt.Scanln(&str)
		hour, _ := strconv.Atoi(str)
		fmt.Print("Min: ")
		fmt.Scanln(&str)
		min, _ := strconv.Atoi(str)
		var event string
		fmt.Print("Event Name: ")
		fmt.Scanln(&event)
		fmt.Println("Adding event...")
		fmt.Println(event + "\n at " + strconv.Itoa(hour) + ":" + strconv.Itoa(min) + " " + day)
		fmt.Println("Shure? Write y/n")
		for {
			var ch string
			fmt.Scanln(&ch)
			if ch == "y" || ch == "n" {
				if ch == "y" {
					schedule.Event = append(schedule.Event, event)
					schedule.Minute = append(schedule.Minute, wday*24*60+hour*60+min)
				}
				break
			}
		}
	}

	SendSchedule(&schedule)
}

//GenerateTestSchedule is test
func GenerateTestSchedule() {
	mins := cloudfunc.GetMinsOfWeek(time.Now().UTC())
	schedule := cloudfunc.Schedule{
		Name:   "test",
		Event:  []string{"We started", "Still alive", "Unbelivable", "5 минут, полёт нормальный", "I`ll send you one more in one min if everything is good"},
		Minute: []int{mins + 6, mins + 7, mins + 8, mins + 9, mins + 10},
	}
	SendSchedule(&schedule)
	//fbclient.AddSubscriber(259224772, "test")
	AddSubscriber(364448153, "test")
}
