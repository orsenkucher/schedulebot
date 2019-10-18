package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/orsenkucher/schedulebot/bot"
	"github.com/orsenkucher/schedulebot/cloudfunc"
	"github.com/orsenkucher/schedulebot/creds"
	"github.com/orsenkucher/schedulebot/dev"
	"github.com/orsenkucher/schedulebot/fbclient"
)

// *** ASAP ***
// Reliable subscription!
// Migalky
//
// *** Proposals ***
//Generate buttoms by path like below
//Ukraine?.Mehmat.firstyear.math.group1.subgroup2
//Custom schedules
//Create kostyl for migalki
//

func main() {
	// fbclient.CreateDemoSched()

	///*
	dev.GenerateTestSchedule()
	fmt.Println("Minuted from week start", cloudfunc.GetMinsOfWeek(time.Now()))
	table := fbclient.FetchTable()
	users := fbclient.FetchSubscribers()
	// fmt.Println(table)

	// key, err := creds.ReadToken()
	// if err != nil {
	// 	panic(err)
	// }
	key := creds.ReadCr()

	// fmt.Println(key)
	b := bot.InitBot(key)
	chans := map[string]chan bot.SubEvent{}

	for _, sch := range table {
		chans[sch.Name] = make(chan bot.SubEvent)
		go bot.ActivateSchedule(sch, users[sch.Name].IDs, b, chans[sch.Name])
	}
	go bot.Listen(b, chans)
	log.Println("bot listen goroutine started")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Start serving")
	go http.ListenAndServe("0.0.0.0:8443", heh2{})
	log.Println("Start serving-2")
	http.ListenAndServe("0.0.0.0:"+port, heh{})
	// http.HandleFunc
	//*/
	//fbclient.CreateSchedule()
}

type heh struct{}

func (h heh) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	log.Println(port + " rec: " + string(bytes))
}

type heh2 struct{}

func (h heh2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	log.Println("8443" + " rec: " + string(bytes))
}
