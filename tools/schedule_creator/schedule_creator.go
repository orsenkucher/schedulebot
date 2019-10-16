package main

import (
	"fmt"
	"strconv"

	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/fbclient"
	"github.com/orsenkucher/schedulebot/tools"
)

func main() {
	CreateSchedule()
}

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
		wday, exist := tools.DayIndex[day]
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

	fbclient.SendSchedule(schedule)
}
