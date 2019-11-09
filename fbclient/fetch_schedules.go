package fbclient

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

const fetchSchsURL = "https://us-central1-scheduleuabot.cloudfunctions.net/FetchSchedules"

// FetchSchedules is public
func FetchSchedules() []cloudfunc.Schedule {
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
	var table []cloudfunc.Schedule
	json.Unmarshal(r, &table)
	return table
}
