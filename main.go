package main

import (
	"log"
	"net/http"

	"github.com/kaakaa/matterpoll-emoji/poll"
)

func main() {
	c, err := poll.LoadConf("config.json")
	if err != nil {
		log.Fatal(err)
	}
	poll.Conf = c
	http.HandleFunc("/poll", poll.PollCmd)
	http.ListenAndServe(":8081", nil)
}
