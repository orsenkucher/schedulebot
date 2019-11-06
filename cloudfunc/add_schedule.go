package cloudfunc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

// Schedule represents firestore Schedule document
// Hash64 is Base64 representation of SHA-1 of this document
type Schedule struct {
	Name   string   `firebase:"name" json:"name"`
	Hash64 string   `firebase:"hash64" json:"hash64"`
	Event  []string `firebase:"event" json:"event"`
	Minute []int    `firebase:"minute" json:"minute"`
	Type   []int    `firebase:"type" json:"type"`
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
