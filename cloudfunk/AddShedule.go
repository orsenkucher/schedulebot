package cloudfunc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// Schedule is public
type Schedule struct {
	Name        string   `firebase:"name" json:"name"`
	Event       []string `firebase:"event" json:"event"`
	Minute      []int    `firebase:"minute" json:"minute"`
	Subscribers []string `firebase:"subscribers" json:"subscribers"`
}

func getMinsOfWeek() int {
	t := time.Now()
	mins := t.Minute() + t.Hour()*60 + int(t.Weekday())*24*60
	return mins
}

// AddSchedule is cloud
func AddSchedule(w http.ResponseWriter, r *http.Request) {
	str, _ := ioutil.ReadAll(r.Body)
	var schedule Schedule
	json.Unmarshal(str, &schedule)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	client.Doc("Schedules/"+schedule.Name).Set(ctx, schedule)
}
