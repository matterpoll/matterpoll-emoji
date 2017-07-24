package poll

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type PollRequest struct {
	TeamId    string
	ChannelId string
	Message   string
	Emojis    []string
}

func NewPollRequest(s string) (*PollRequest, error) {
	u, err := url.Parse("http://dummy.com/?" + s)
	if err != nil {
		return nil, err
	}

	p := &PollRequest{}
	for key, values := range u.Query() {
		switch key {
		case "team_id":
			if p.TeamId = values[0]; len(p.TeamId) == 0 {
				return nil, fmt.Errorf("Unexpected Error: TeamID in request is empty.")
			}
		case "channel_id":
			if p.ChannelId = values[0]; len(p.ChannelId) == 0 {
				return nil, fmt.Errorf("Unexpected Error: ChannelID in request is empty.")
			}
		case "text":
			p.Message, p.Emojis, err = parseText(values[0])
			if err != nil {
				return nil, err
			}
		default:
		}
	}
	return p, nil
}

func parseText(text string) (string, []string, error) {
	var re *(regexp.Regexp)
	switch text[0] {
	case '`':
		re = regexp.MustCompile("`([^`]+)`(.+)")
	case '\'':
		re = regexp.MustCompile("'([^']+)'(.+)")
	case '"':
		re = regexp.MustCompile("\"([^\"]+)\"(.+)")
	default:
		return "", nil, fmt.Errorf("Command Error: /poll `Here is description` :thumbsup: :thumbsdown:...")
	}
	e := re.FindStringSubmatch(text)
	if len(e) != 3 {
		return "", nil, fmt.Errorf("Command Error: /poll `Here is description` :thumbsup: :thumbsdown:...")
	}
	var emojis []string
	for _, v := range strings.Split(e[2], " ") {
		if len(v) == 0 {
			continue
		}
		if len(v) < 3 || !strings.HasPrefix(v, ":") || !strings.HasSuffix(v, ":") {
			return "", nil, fmt.Errorf("Emoji Error: %s is not emoji format", v)
		}
		emojis = append(emojis, v[1:len(v)-1])
	}
	v := strings.Split(text, " ")
	if len(v) < 2 {
		return "", nil, fmt.Errorf("Error: /poll description emoji1 emoji2...")
	}
	return e[1], emojis, nil
}
