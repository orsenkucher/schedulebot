package fbclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

const fetchSubscURL = "https://us-central1-scheduleuabot.cloudfunctions.net/FetchSubscribers"

// FetchSubscribers is public
func FetchSubscribers() map[string][]cloudfunc.Subscriber {
	resp, _ := http.Get(fetchSubscURL)
	r, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(r))
	var subs map[string][]cloudfunc.Subscriber
	json.Unmarshal(r, &subs)
	return subs
}
