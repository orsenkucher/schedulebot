package fbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/hash"
)

const addSchURL = "https://us-central1-scheduleuabot.cloudfunctions.net/AddSchedule"

// SendSchedule sends Schedule to firestore
func SendSchedule(schedule *cloudfunc.Schedule) {
	schedule.Hash64 = hash.EncodeAsBase64(schedule)
	strb, _ := json.Marshal(schedule)
	fmt.Printf("Sending %v bytes of %s schedule\n", len(strb), schedule.Name)
	// prettystrb, _ := json.MarshalIndent(schedule, "", "\t")
	// fmt.Println(string(prettystrb))
	_, err := http.Post(addSchURL, "application/json", bytes.NewBuffer(strb))
	if err != nil {
		panic(err)
	}
}
