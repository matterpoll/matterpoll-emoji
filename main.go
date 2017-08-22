package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/kaakaa/matterpoll-emoji/poll"
)

var (
	port    = flag.Int("p", 8505, "port number")
	address = flag.String("a", "", "optional address to bind and listen on")
	config  = flag.String("c", "config.json", "optional path to the config file")
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	flag.Parse()

	c, err := poll.LoadConf(*config)
	if err != nil {
		log.Fatal(err)
	}
	poll.Conf = c
	http.HandleFunc("/poll", poll.PollCmd)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *address, *port), nil); err != nil {
		log.Fatal(err)
	}
}
