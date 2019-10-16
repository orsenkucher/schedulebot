package fbclient

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/orsenkucher/schedulebot/cloudfunc"
)

// SetSubscribers is public
func SetSubscribers(users []int64, schName string) {
	usersstr := []string{}

	for i := 0; i < len(users); i++ {
		usersstr = append(usersstr, strconv.FormatInt(users[i], 10))
	}

	subscribers := cloudfunc.Subscribers{
		Name: schName,
		IDs:  usersstr,
	}

	strb, _ := json.Marshal(&subscribers)
	http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/SetSubscribers", "application/json", bytes.NewBuffer(strb))
}
