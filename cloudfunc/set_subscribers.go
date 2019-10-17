package cloudfunc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

// SubscriberQuerie is public
type SubscriberQuerie struct {
	ID           string `firebase:"id" json:"id"`
	ScheduleName string `firebase:"schedulename" json:"schedulename"`
}

// SetSubscribers is public
func SetSubscribers(w http.ResponseWriter, r *http.Request) {
	str, _ := ioutil.ReadAll(r.Body)
	var subscribers Subscribers
	json.Unmarshal(str, &subscribers)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	subscribersRef := client.Doc("Subscribers/" + subscribers.Name)
	subscribersRef.Set(ctx, subscribers)
}

// AddSubscriber is public
func AddSubscriber(w http.ResponseWriter, r *http.Request) {
	str, _ := ioutil.ReadAll(r.Body)
	var subscriberq SubscriberQuerie
	json.Unmarshal(str, &subscriberq)

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	subscribersRef := client.Doc("Schedules/" + subscriberq.ScheduleName + "/Subscribers/" + subscriberq.ID)
	subscribersRef.Set(ctx, Subscriber{ID: subscriberq.ID})
}
