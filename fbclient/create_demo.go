package fbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

func invertedDayIndex() map[int]string {
	res := map[int]string{}
	for k, v := range DayIndex {
		res[v] = k
	}
	return res
}

// CreateDemoSched creates demo schedule on firestore
func CreateDemoSched() {
	schedule := createDemoSched()

	strb, err := json.Marshal(&schedule)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(strb))

	fmt.Println("Sending json...")
	resp, err := http.Post(addSchURL, "application/json", bytes.NewBuffer(strb))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(r))
	fmt.Println("done")
}

func createDemoSched() cloudfunc.Schedule {
	schedule := cloudfunc.Schedule{
		Name:   "demo",
		Event:  []string{},
		Minute: []int{},
	}

	idx := invertedDayIndex()
	for wd := 0; wd < 7; wd++ {
		for hour := 0; hour < 24; hour++ {
			for min := 0; min < 60; min += 5 { // every 5 mins
				event := fmt.Sprintf("%s %v:%v", idx[wd], hour, min)
				schedule.Event = append(schedule.Event, event)
				schedule.Minute = append(schedule.Minute, wd*24*60+hour*60+min)
			}
		}
	}
	return schedule
}
