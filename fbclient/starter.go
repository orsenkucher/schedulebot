package fbclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

type Table []cloudfunc.Schedule

const fetchURL = "https://us-central1-scheduleuabot.cloudfunctions.net/FetchSchedules"

func FetchTable() Table {
	resp, _ := http.Get(fetchURL)
	r, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(r))
	var table Table
	json.Unmarshal(r, &table)
	return table
}
