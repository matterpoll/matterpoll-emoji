package poll

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"time"

	"github.com/mattermost/platform/model"
)

var Conf *PollConf

func PollCmd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("Error", err)
		return
	}

	poll, err := NewPollRequest(string(b))
	if err != nil {
		log.Print("Error", err)
		return
	}

	var responce = `{
		"response_type": "in_channel",
		"text": "` + poll.Message + ` #poll",
		"username": "Poll Bot",
		"icon_url": "https://www.mattermost.org/wp-content/uploads/2016/04/icon.png"
	}`

	c := model.NewClient(Conf.Host)
	c.TeamId = poll.TeamId

	var user *model.User
	user, err = login(c)
	var userId = user.Id
	if err != nil {
		log.Print("Error", err)
		return
	}

	var result *model.Result
	io.WriteString(w, responce)
	/*
		var start = time.Now().Unix()

		time.Sleep(30 * time.Second)

		result, err = c.GetPostsSince(poll.ChannelId, start)
		if err != nil {
			log.Println("Error", err)
			return
		}
	*/
	result, err = c.GetPosts(poll.ChannelId, 0, 10, "")
	if err != nil {
		log.Println("Error", err)
		return
	}

	var postId = result.Data.(*model.PostList).Order[0]
	log.Println(result.Data.(*model.PostList).Order)
	for _, postId := range result.Data.(*model.PostList).Order {
		log.Println("Message is:", result.Data.(*model.PostList).Posts[postId].Message)
	}
	log.Println("PostId is:", result.Data.(*model.PostList).Posts[postId].Id)
	log.Println("Message is:", result.Data.(*model.PostList).Posts[postId].Message)
	log.Println("UserId is:", result.Data.(*model.PostList).Posts[postId].UserId)

	err = reaction(c, poll.ChannelId, userId, postId, poll.Emojis)
	if err != nil {
		log.Print("Error", err)
		return
	}
}

func login(c *model.Client) (*model.User, error) {
	r, err := c.Login(Conf.User.Id, Conf.User.Password)
	if err != nil {
		return nil, err
	}
	return r.Data.(*model.User), nil
}

func post(c *model.Client, poll *PollRequest) (*model.Post, error) {
	p := model.Post{
		ChannelId: poll.ChannelId,
		Message:   poll.Message + " #poll",
	}
	r, err := c.CreatePost(&p)
	if err != nil {
		return nil, err
	}
	return r.Data.(*model.Post), nil
}

func reaction(c *model.Client, channelId string, userId string, postId string, emojis []string) error {
	for _, e := range emojis {
		r := model.Reaction{
			UserId:    userId,
			PostId:    postId,
			EmojiName: e,
		}
		_, err := c.SaveReaction(channelId, &r)
		if err != nil {
			return err
		}
	}
	return nil
}
