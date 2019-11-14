package cloudfunc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

// SubscriberQuery is public
type SubscriberQuery struct {
	ID           string `firebase:"ID" json:"ID"`
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
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	str, _ := ioutil.ReadAll(r.Body)
	var subscriberq SubscriberQuery
	json.Unmarshal(str, &subscriberq)

	userRef := client.Collection("Subscribers").Doc(subscriberq.ID)
	userRef.Set(ctx, struct {
		ID string `firebase:"id"`
	}{ID: subscriberq.ID})
	userRef.Collection("Schedules").Doc(subscriberq.ScheduleName).Set(ctx, struct {
		Schedule string `firebase:"schedule"`
	}{Schedule: subscriberq.ScheduleName})
	subscribersRef := client.Collection("Schedules").Doc(subscriberq.ScheduleName).Collection("Subscribers").Doc(subscriberq.ID)
	subscribersRef.Set(ctx, Subscriber{ID: subscriberq.ID})
}

// DeleteSubscriber is public
func DeleteSubscriber(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	str, _ := ioutil.ReadAll(r.Body)
	var subscriberq SubscriberQuery
	json.Unmarshal(str, &subscriberq)
	client.Collection("Schedules").Doc(subscriberq.ScheduleName).Collection("Subscribers").Doc(subscriberq.ID).Delete(ctx)
	client.Collection("Subscribers").Doc(subscriberq.ID).Collection("Schedules").Doc(subscriberq.ScheduleName).Delete(ctx)
}
