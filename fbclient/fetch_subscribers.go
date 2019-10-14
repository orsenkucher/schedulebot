package fbclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

const fetchSubscURL = "https://us-central1-scheduleuabot.cloudfunctions.net/FetchSubscribers"

// FetchSubscribers is public
func FetchSubscribers() map[string]cloudfunc.Subscribers {
	resp, _ := http.Get(fetchSubscURL)
	r, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(r))
	var table map[string]cloudfunc.Subscribers
	json.Unmarshal(r, &table)
	return table
}
