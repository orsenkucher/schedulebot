package main

import (
	"fmt"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/creds"
)

func main() {
	key, err := creds.ReadToken()
	if err != nil {
		panic(err)
	}

	fmt.Println(key)
	bot.InitBot(key)
}
