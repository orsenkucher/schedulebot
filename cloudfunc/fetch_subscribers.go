package cloudfunc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

// Subscribers ispublic
type Subscribers struct {
	Name string   `firebase:"name" json:"name"`
	IDs  []string `firebase:"IDs" json:"IDs"`
}

//Subscriber is public
type Subscriber struct {
	ID string `firebase:"ID" json:"ID"`
}

//FetchSubscribers returns subscribers
func FetchSubscribers(w http.ResponseWriter, r *http.Request) {
	schs := fetchSubscribers()
	js, err := json.Marshal(&schs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(w, string(js))
}

func fetchSubscribers() map[string][]Subscriber {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	docsIter := client.Collection("Schedules").DocumentRefs(ctx)
	docs, err := docsIter.GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	subs := map[string][]Subscriber{}
	for _, doc := range docs {
		subsList := []Subscriber{}
		subsIter := doc.Collection("Subscribers").Documents(ctx)
		subsDocs, err := subsIter.GetAll()

		if err == nil {
			for _, subsDoc := range subsDocs {
				var sub Subscriber
				subsDoc.DataTo(&sub)
				subsList = append(subsList, sub)
			}
		}
		subs[doc.ID] = subsList
	}

	return subs
}
