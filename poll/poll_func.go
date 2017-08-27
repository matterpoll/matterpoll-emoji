package poll

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mattermost/platform/model"
)

const (
	// ResponseUsername is the username which will be used to post the slack command response
	ResponseUsername = "Matterpoll"
	// ResponseIconURL is the profile picture which will be used to post the slack command response
	ResponseIconURL = "https://www.mattermost.org/wp-content/uploads/2016/04/icon.png"
)

// Server handles slash commands from a mattermost instance. One sever may handle multiple requests from one mattermost instance. It uses a provided configuration to handle the requests.
type Server struct {
	Conf *Conf
}

// Cmd handles a slash command request and sends back a response
func (ps Server) Cmd(w http.ResponseWriter, r *http.Request) {
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
	poll, err := NewRequest(r.Form)
	validPoll := err == nil

	var response model.CommandResponse
	response.Username = ResponseUsername
	response.IconURL = ResponseIconURL

	if validPoll && len(ps.Conf.Token) != 0 && ps.Conf.Token != poll.Token {
		validPoll = false
		err = fmt.Errorf(ErrorTokenMissmatch)
	}
	if validPoll {
		response.ResponseType = model.COMMAND_RESPONSE_TYPE_IN_CHANNEL
		response.Text = poll.Message + ` #poll`
	} else {
		response.ResponseType = model.COMMAND_RESPONSE_TYPE_EPHEMERAL
		response.Text = err.Error()
	}
	if _, err := io.WriteString(w, response.ToJson()); err != nil {
		log.Print(err)
		return
	}
	if validPoll {
		c := model.NewAPIv4Client(ps.Conf.Host)
		user, err := ps.login(c)
		if err != nil {
			log.Print(err)
			return
		}
		go ps.addReaction(c, user, poll)
	}
}

func (ps Server) login(c *model.Client4) (*model.User, error) {
	u, apiResponse := c.Login(ps.Conf.User.ID, ps.Conf.User.Password)
	if apiResponse != nil && apiResponse.StatusCode != 200 {
		return nil, fmt.Errorf("Error: Login failed. API statuscode: %v", apiResponse.StatusCode)
	}
	return u, nil
}

func (ps Server) addReaction(c *model.Client4, user *model.User, poll *Request) {
	for try := 0; try < 5; try++ {
		// Get the last post and compare it to our message text
		result, apiResponse := c.GetPostsForChannel(poll.ChannelID, 0, 1, "")
		if apiResponse != nil && apiResponse.StatusCode != 200 {
			log.Printf("Error: Failed to fetch posts. API statuscode: %v", apiResponse.StatusCode)
			return
		}
		postID := result.Order[0]
		if result.Posts[postID].Message == poll.Message+" #poll" {
			err := reaction(c, poll.ChannelID, user.Id, postID, poll.Emojis)
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

func reaction(c *model.Client4, channelID string, userID string, postID string, emojis []string) error {
	for _, e := range emojis {
		r := model.Reaction{
			UserId:    userID,
			PostId:    postID,
			EmojiName: e,
		}
		_, apiResponse := c.SaveReaction(&r)
		if apiResponse != nil && apiResponse.StatusCode != 200 {
			return fmt.Errorf("Error: Failed to save reaction. API statuscode: %v", apiResponse.StatusCode)
		}
	}
	return nil
}
