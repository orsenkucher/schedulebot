package fbclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

const schFile = "fbclient/sch.json"

// Event represents an event
type Event struct {
	Title string `json:"title"`
	Day   string `json:"day"`
	Time  string `json:"time"`
}

// Schedule represents my custom schedule scheme
type Schedule struct {
	Name   string  `json:"name"`
	Events []Event `json:"events"`
}

// CreateSchFromJSON is schedule creator v2, that read data from json
func CreateSchFromJSON() {
	var wg sync.WaitGroup

	bytes, err := readJSON(schFile)
	if err != nil {
		panic(err)
	}

	var schs []Schedule
	json.Unmarshal(bytes, &schs)
	fmt.Println(schs)

	for _, sch := range schs {
		schedule := cloudfunc.Schedule{
			Name:   sch.Name,
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
			schedule.Event = append(schedule.Event, e.Title)
			schedule.Minute = append(schedule.Minute, dayIdx*24*60+hour*60+minute)
		}

		wg.Add(1)
		fmt.Printf("Sending %v\n", schedule)
		doneCh := make(chan struct{})
		go func(doneCh <-chan struct{}, wg *sync.WaitGroup) {
			for range doneCh {
				fmt.Println("Sent " + schedule.Name)
				wg.Done()
			}
		}(doneCh, &wg)
		go SendScheduleGo(&schedule, doneCh)
	}

	fmt.Println("Waiting for cloud functions")
	wg.Wait()
}

func readJSON(path string) ([]byte, error) {
	j, err := ioutil.ReadFile(schFile)
	if err != nil {
		return nil, err
	}
	return j, nil
}
