package fbclient

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

const addSchURL = "https://us-central1-scheduleuabot.cloudfunctions.net/AddSchedule"

// SendSchedule sends Schedule to firestore
func SendSchedule(schedule *cloudfunc.Schedule) {
	schedule.Hash64 = calcSchHashAsBase64(schedule)
	strb, _ := json.Marshal(schedule)
	fmt.Printf("Sending %v bytes of %s schedule\n", len(strb), schedule.Name)
	// prettystrb, _ := json.MarshalIndent(schedule, "", "\t")
	// fmt.Println(string(prettystrb))
	_, err := http.Post(addSchURL, "application/json", bytes.NewBuffer(strb))
	if err != nil {
		panic(err)
	}
}

func calcSchHashAsBase64(sch *cloudfunc.Schedule) string {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(sch); err != nil {
		panic(err)
	}
	bytes := b.Bytes()
	sha := sha1.New()
	sha.Write(bytes)
	bytes = sha.Sum(nil)
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return b64
}
