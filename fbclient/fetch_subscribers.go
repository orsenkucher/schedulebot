package fbclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

const fetchSubsURL = "https://us-central1-scheduleuabot.cloudfunctions.net/FetchSubscribers"

// FetchSubscribers is public
func FetchSubscribers() map[string][]cloudfunc.Subscriber {
	resp, err := http.Get(fetchSubsURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(r))
	var subs map[string][]cloudfunc.Subscriber
	json.Unmarshal(r, &subs)
	return subs
}
