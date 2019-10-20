package fbclient

import (
	"bytes"
	"encoding/json"
	"log"
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
	_, err := http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/SetSubscribers", "application/json", bytes.NewBuffer(strb))
	if err != nil {
		log.Fatalln(err)
	}
}

// AddSubscriber is public
func AddSubscriber(user int64, schName string) {
	subscriberq := cloudfunc.SubscriberQuerie{
		ScheduleName: schName,
		ID:           strconv.FormatInt(user, 10),
	}

	strb, _ := json.Marshal(&subscriberq)
	_, err := http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/AddSubscriber", "application/json", bytes.NewBuffer(strb))
	if err != nil {
		log.Fatalln(err)
	}
}

// DeleteSubscriber is public
func DeleteSubscriber(user int64, schName string) {
	subscriberq := cloudfunc.SubscriberQuerie{
		ScheduleName: schName,
		ID:           strconv.FormatInt(user, 10),
	}

	strb, _ := json.Marshal(&subscriberq)
	_, err := http.Post("https://us-central1-scheduleuabot.cloudfunctions.net/DeleteSubscriber", "application/json", bytes.NewBuffer(strb))
	if err != nil {
		log.Fatalln(err)
	}
}
