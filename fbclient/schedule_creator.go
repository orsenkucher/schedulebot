package fbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

// SendSchedule sends Schedule to firestore
func SendSchedule(schedule *cloudfunc.Schedule) {
	strb, _ := json.Marshal(schedule)
	fmt.Printf("Sending %v bytes of %s schedule\n", len(strb), schedule.Name)
	// prettystrb, _ := json.MarshalIndent(schedule, "", "\t")
	// fmt.Println(string(prettystrb))
	_, err := http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/AddSchedule", "application/json", bytes.NewBuffer(strb))
	if err != nil {
		panic(err)
	}
}
