package main

import (
	"flag"
	"fmt"
	"github.com/kaakaa/matterpoll-emoji/poll"
	"log"
	"net/http"
)

var port = flag.Int("p", 8505, "port number")

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	c, err := poll.LoadConf("config.json")
	if err != nil {
		log.Fatal(err)
	}
	ps := poll.PollServer{c}
	http.HandleFunc("/poll", ps.PollCmd)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}
