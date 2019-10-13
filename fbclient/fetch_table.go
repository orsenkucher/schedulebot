package fbclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

// Table is public
type Table []cloudfunc.Schedule

const fetchURL = "https://us-central1-scheduleuabot.cloudfunctions.net/FetchSchedules"

// FetchTable is public
func FetchTable() Table {
	resp, _ := http.Get(fetchURL)
	r, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(r))
	var table Table
	json.Unmarshal(r, &table)
	return table
}
