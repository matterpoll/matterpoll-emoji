package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kaakaa/matterpoll-emoji/poll"
	"github.com/mattermost/platform/model"
)

var conf *poll.PollConf

func pollCmd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}
	poll, err := NewPollRequest(string(b))
	if err != nil {
		log.Print(err)
		return
	}
	c := model.NewClient(conf.Host)
	c.TeamId = poll.TeamId

	_, err = login(c)
	if err != nil {
		log.Print(err)
		return
	}
	p, err := post(c, poll)
	if err != nil {
		log.Print(err)
		return
	}
	reaction(c, p, poll)
	fmt.Fprintf(w, "{'text': 'hello'}")
}

func login(c *model.Client) (*model.User, error) {
	r, err := c.Login(conf.User.Id, conf.User.Password)
	if err != nil {
		return nil, err
	}
	return r.Data.(*model.User), nil
}

func post(c *model.Client, poll *PollRequest) (*model.Post, error) {
	p := model.Post{
		UserId:    "kaakaa",
		ChannelId: poll.ChannelId,
		Message:   poll.Message + " #poll",
	}
	r, err := c.CreatePost(&p)
	if err != nil {
		return nil, err
	}
	return r.Data.(*model.Post), nil
}

func reaction(c *model.Client, p *model.Post, poll *PollRequest) {
	for _, e := range poll.Emojis {
		r := model.Reaction{
			UserId:    p.UserId,
			PostId:    p.Id,
			EmojiName: e,
		}
		c.SaveReaction(p.ChannelId, &r)
	}
}

func main() {
	c, err := poll.LoadConf("config.json")
	if err != nil {
		log.Fatal(err)
	}
	conf = c
	http.HandleFunc("/poll", pollCmd)
	http.ListenAndServe(":8080", nil)
}
