package fbclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

// Table is public
type Table []cloudfunc.Schedule

const fetchSchsURL = "https://us-central1-scheduleuabot.cloudfunctions.net/FetchSchedules"

// FetchTable is public
func FetchTable() Table {
	resp, err := http.Get(fetchSchsURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(string(r))
	var table Table
	json.Unmarshal(r, &table)
	return table
}
