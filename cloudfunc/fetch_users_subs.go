package cloudfunc

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

// FetchUsersSubs give subs of user
func FetchUsersSubs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "scheduleuabot")
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	str, _ := ioutil.ReadAll(r.Body)
	arr := []string{}

	doc := client.Doc("Subscribers/" + string(str))
	schedules, _ := doc.Collection("Schedules").Documents(ctx).GetAll()

	for _, schedule := range schedules {
		arr = append(arr, schedule.Ref.ID)
	}
	log.Print(arr)
	fmt.Fprint(w, arr)
}
