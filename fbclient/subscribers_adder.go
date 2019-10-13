package fbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

// AddSubscribers is public
func AddSubscribers(users []int64, schName string) {
	usersstr := []string{}

	for i := 0; i < len(users); i++ {
		usersstr = append(usersstr, strconv.FormatInt(users[i], 10))
	}

	que := cloudfunc.SubscriberQuerie{
		ScheduleName: schName,
		IDs:          usersstr,
	}

	strb, _ := json.Marshal(&que)
	resp, _ := http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/AddSubscribers", "application/json", bytes.NewBuffer(strb))
	r, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(r))
}