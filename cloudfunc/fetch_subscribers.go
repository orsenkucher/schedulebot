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

//FetchSubscribers returns subscribers
func FetchSubscribers(w http.ResponseWriter, r *http.Request) {
	schs := fetchSubscribers()
	js, err := json.Marshal(&schs)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Fprint(w, string(js))
}

func fetchSubscribers() map[string]Subscribers {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	docsIter := client.Collection("Subscribers").Documents(ctx)
	docs, err := docsIter.GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	schs := map[string]Subscribers{}
	for _, doc := range docs {
		sch := Subscribers{}
		err := doc.DataTo(&sch)
		if err != nil {
			log.Fatalln(err)
		}
		schs[sch.Name] = sch
	}

	return schs
}
