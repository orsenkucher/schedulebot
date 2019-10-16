package cloudfunc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

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
