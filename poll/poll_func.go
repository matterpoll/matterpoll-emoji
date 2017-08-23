package poll

import (
	"fmt"
	"github.com/mattermost/platform/model"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	ResponseUsername = "Matterpoll"
	ResponseIconUrl  = "https://www.mattermost.org/wp-content/uploads/2016/04/icon.png"
)

type PollServer struct {
	Conf *PollConf
}

func (ps PollServer) PollCmd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Check if Content Type is correct
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	poll, err := NewPollRequest(r.Form)
	valid_poll := err == nil

	var response model.CommandResponse
	response.Username = ResponseUsername
	response.IconURL = ResponseIconUrl

	if valid_poll && len(ps.Conf.Token) != 0 && ps.Conf.Token != poll.Token {
		valid_poll = false
		err = fmt.Errorf(ErrorTokenMissmatch)
	}
	if valid_poll {
		response.ResponseType = model.COMMAND_RESPONSE_TYPE_IN_CHANNEL
		response.Text = poll.Message + ` #poll`
	} else {
		response.ResponseType = model.COMMAND_RESPONSE_TYPE_EPHEMERAL
		response.Text = err.Error()
	}
	io.WriteString(w, response.ToJson())
	if valid_poll {
		c := model.NewAPIv4Client(ps.Conf.Host)
		user, err := ps.login(c)
		if err != nil {
			log.Print(err)
			return
		}
		go ps.addReaction(c, user, poll)
	}
}

func (ps PollServer) login(c *model.Client4) (*model.User, error) {
	u, api_response := c.Login(ps.Conf.User.Id, ps.Conf.User.Password)
	if api_response != nil && api_response.StatusCode != 200 {
		return nil, fmt.Errorf("Error: Login failed. API statuscode: %v", api_response.StatusCode)
	}
	return u, nil
}

func (ps PollServer) addReaction(c *model.Client4, user *model.User, poll *PollRequest) {
	for try := 0; try < 5; try++ {
		// Get the last post and compare it to our message text
		result, api_response := c.GetPostsForChannel(poll.ChannelId, 0, 1, "")
		if api_response != nil && api_response.StatusCode != 200 {
			log.Printf("Error: Failed to fetch posts. API statuscode: %v", api_response.StatusCode)
			return
		}
		postId := result.Order[0]
		if result.Posts[postId].Message == poll.Message+" #poll" {
			err := Reaction(c, poll.ChannelId, user.Id, postId, poll.Emojis)
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

func Reaction(c *model.Client4, channelId string, userId string, postId string, emojis []string) error {
	for _, e := range emojis {
		r := model.Reaction{
			UserId:    userId,
			PostId:    postId,
			EmojiName: e,
		}
		_, api_response := c.SaveReaction(&r)
		if api_response != nil && api_response.StatusCode != 200 {
			return fmt.Errorf("Error: Failed to save reaction. API statuscode: %v", api_response.StatusCode)
		}
	}
	return nil
}
