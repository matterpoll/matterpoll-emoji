package main

import (
	"fmt"
	"log"
	"flag"
	"net/http"

	"github.com/kaakaa/matterpoll-emoji/poll"
)

var config  = flag.String(
	"c", "config.json", "optional path to the config file")

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	c, err := poll.LoadConf(*config)
	if err != nil {
		log.Fatal(err)
	}
	ps := poll.PollServer{Conf: c}
	http.HandleFunc("/poll", ps.PollCmd)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), nil); err != nil {
		log.Fatal(err)
	}
}