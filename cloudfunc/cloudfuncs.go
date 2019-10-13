package cloudfunc

import (
	"context"
	"encoding/json"
	"fmt"
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

// SubscriberQuerie is public
type SubscriberQuerie struct {
	ScheduleName string   `json:"schName"`
	IDs          []string `json:"IDs"`
}

// GetMinsOfWeek return mins passed from week start
func GetMinsOfWeek(t time.Time) int {
	mins := t.Minute() + t.Hour()*60 + int(t.Weekday())*24*60
	return mins
}

// AddSchedule appends Schedule to firebase collection
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

// FetchSchedules returns Schedules json
func FetchSchedules(w http.ResponseWriter, r *http.Request) {
	schs := fetchSchedules()
	js, err := json.Marshal(&schs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(w, string(js))
}

func fetchSchedules() []Schedule {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	docsIter := client.Collection("Schedules").Documents(ctx)
	docs, err := docsIter.GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	schs := make([]Schedule, 0, len(docs))
	for _, doc := range docs {
		sch := Schedule{}
		err := doc.DataTo(&sch)
		if err != nil {
			log.Fatalln(err)
		}
		schs = append(schs, sch)
	}

	return schs
}

// AddSubscribers is public
func AddSubscribers(w http.ResponseWriter, r *http.Request) {
	str, _ := ioutil.ReadAll(r.Body)
	var querie SubscriberQuerie
	json.Unmarshal(str, &querie)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	scheduleRef := client.Doc("Schedules/" + querie.ScheduleName)
	var schedule Schedule
	docsnap, _ := scheduleRef.Get(ctx)
	docsnap.DataTo(&schedule)

	schedule.Subscribers = append(schedule.Subscribers, querie.IDs...)
	scheduleRef.Set(ctx, schedule)
}
