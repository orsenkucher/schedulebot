package fbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

// SendSchedule is public
func SendSchedule(schedule cloudfunc.Schedule) {
	strb, _ := json.Marshal(&schedule)
	fmt.Println("Sending json...")
	prettystrb, _ := json.MarshalIndent(&schedule, "", "\t")
	fmt.Println(string(prettystrb))
	resp, _ := http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/AddSchedule", "application/json", bytes.NewBuffer(strb))
	r, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(r))
}
