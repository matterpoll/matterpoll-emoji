package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kaakaa/matterpoll-emoji/poll"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	c, err := poll.LoadConf("config.json")
	if err != nil {
		log.Fatal(err)
	}
	ps := poll.PollServer{c}
	http.HandleFunc("/poll", ps.Cmd)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *c.Port), nil); err != nil {
		log.Fatal(err)
	}
}
