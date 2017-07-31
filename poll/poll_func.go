package poll

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mattermost/platform/model"
)

var Conf *PollConf

func PollCmd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	err := r.ParseForm()
	if err != nil {
		log.Print("Error: ", err)
		return
	}
	poll, err := NewPollRequest(r.Form)
	var response_type string
	var response_text string
	if err == nil {
		response_type = "in_channel"
		response_text = poll.Message + ` #poll`
	} else {
		response_type = "ephemeral"
		response_text = err.Error()
	}

	var response = `{
		"response_type": "` + response_type + `",
		"text": "` + response_text + `",
		"username": "Matterpoll",
		"icon_url": "https://www.mattermost.org/wp-content/uploads/2016/04/icon.png"
	}`
	io.WriteString(w, response)
	if err == nil {
		if len(Conf.Token) != 0 && Conf.Token != poll.Token {
			log.Print("Token missmatch. Check you config.json")
			return
		}

		c := model.NewAPIv4Client(Conf.Host)
		var user *model.User
		user, err = login(c)
		if err != nil {
			log.Print(err)
			return
		}
		go addReaction(c, user, poll)
	}
}

func login(c *model.Client4) (*model.User, error) {
	u, api_response := c.Login(Conf.User.Id, Conf.User.Password)
	if api_response != nil && api_response.StatusCode != 200 {
		return nil, fmt.Errorf("Error: Login failed. API statuscode: ", api_response.StatusCode)
	}
	return u, nil
}

func addReaction(c *model.Client4, user *model.User, poll *PollRequest) {
	for try := 0; try < 5; try++ {
		// Get the last post and compare it to our message text
		result, api_response := c.GetPostsForChannel(poll.ChannelId, 0, 1, "")
		if api_response != nil && api_response.StatusCode != 200 {
			log.Println("Error: Failed to fetch posts. API statuscode: ", api_response.StatusCode)
			return
		}
		var postId = result.Order[0]
		if result.Posts[postId].Message == poll.Message+" #poll" {
			err := reaction(c, poll.ChannelId, user.Id, postId, poll.Emojis)
			if err != nil {
				log.Print(err)
				return
			}
			return
		}
		// Try again later
		time.Sleep(100 * time.Millisecond)
	}
}

func reaction(c *model.Client4, channelId string, userId string, postId string, emojis []string) error {
	for _, e := range emojis {
		r := model.Reaction{
			UserId:    userId,
			PostId:    postId,
			EmojiName: e,
		}
		_, api_response := c.SaveReaction(&r)
		if api_response != nil && api_response.StatusCode != 200 {
			return fmt.Errorf("Error: Failed to save reaction. API statuscode: ", api_response.StatusCode)
		}
	}
	return nil
}
