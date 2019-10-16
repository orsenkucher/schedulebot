package cloudfunc

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
)

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
