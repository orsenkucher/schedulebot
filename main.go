package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	/*key, err := creds.ReadToken()
	if err != nil {
		panic(err)
	}

	fmt.Println(key)
	bot.InitBot(key)//*/

	t := time.Now()
	fmt.Println(t.Weekday())
	fmt.Println(strconv.Itoa(t.Hour()) + ":" + strconv.Itoa(t.Minute()))
}
