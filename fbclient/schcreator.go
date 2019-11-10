package fbclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/root"
)

const mpw = root.MPW
const schFile = "fbclient/sch.json"
const rootdir = "data"

// Event represents an event
type Event struct {
	Title string `json:"title"`
	Day   string `json:"day"`
	Time  string `json:"time"`
	Spin  string `json:"spin"`
}

// Schedule represents my custom schedule scheme
type Schedule struct {
	Name   string  `json:"name"`
	Events []Event `json:"events"`
}

// CreateSchFromJSON is schedule creator v2, that read data from json
func CreateSchFromJSON() {
	bytes, err := readJSON(schFile)
	if err != nil {
		panic(err)
	}

	var schs []Schedule
	json.Unmarshal(bytes, &schs)
	for _, s := range schs {
		for _, e := range s.Events {
			fmt.Println(e)
		}
	}
	// fmt.Scanln()
	fireSchs := makeFirestoreSchedules(schs)
	json.Unmarshal(bytes, &schs)
	for _, s := range fireSchs {
		for i, e := range s.Type {
			fmt.Println(e, s.Event[i])
		}
	}
	// fmt.Scanln()

	count := len(fireSchs)
	var wg sync.WaitGroup
	wg.Add(count)

	fmt.Printf("Sending %v fire schedules\n", count)
	for i := 0; i < count; i++ {
		go func(fsch *cloudfunc.Schedule) {
			SendSchedule(fsch)
			fmt.Println("Sent " + fsch.Name)
			wg.Done()
		}(&fireSchs[i]) // it is safe to path pointer
	}

	fmt.Println("Waiting for cloud functions")
	wg.Wait()
}

func makeFirestoreSchedules(schs []Schedule) []cloudfunc.Schedule {
	fireSchs := make([]cloudfunc.Schedule, 0, len(schs))
	for _, sch := range schs {
		schedule := cloudfunc.Schedule{
			Name:   sch.Name,
			Type:   make([]int, 0, len(sch.Events)),
			Event:  make([]string, 0, len(sch.Events)),
			Minute: make([]int, 0, len(sch.Events))}
		for _, e := range sch.Events {
			dayIdx, ok := DayIndex[e.Day]
			if !ok {
				panic("Invalid Day on " + e.Title)
			}

			timePair := strings.Split(e.Time, ":")
			if len(timePair) != 2 {
				panic("Invalid Time on " + e.Title)
			}
			hour, err := strconv.Atoi(timePair[0])
			if err != nil {
				panic("Invalid Hour on " + e.Title)
			}
			minute, err := strconv.Atoi(timePair[1])
			if err != nil {
				panic("Invalid Minute on " + e.Title)
			}
			spin := -1
			if e.Spin == "up" {
				spin = 0
			}
			if e.Spin == "down" {
				spin = 1
			}
			schedule.Type = append(schedule.Type, spin)
			schedule.Event = append(schedule.Event, e.Title)
			schedule.Minute = append(schedule.Minute, (dayIdx*24*60+(hour-2)*60+minute+mpw)%mpw)
		}
		fireSchs = append(fireSchs, schedule)
	}
	return fireSchs
}

func readJSON(path string) ([]byte, error) {
	j, err := ioutil.ReadFile(schFile)
	if err != nil {
		return nil, err
	}
	return j, nil
}
