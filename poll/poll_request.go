package poll

import (
	"fmt"
	"regexp"
	"strings"
)

type PollRequest struct {
	ChannelId string
	Token     string
	Message   string
	Emojis    []string
}

const (
	backTick             = "`"
	ErrorTextWrongFormat = `The message format is wrong. Try this instead: ` + backTick + `/poll \"What do you gys wanna grab for lunch?\" :pizza: :sushi:` + backTick
	ErrorTokenMissmatch  = `An error occurred. Ask your administrator to check the Matterpoll config settings.`
	ErrorWrongLength     = `An error occurred. Try the same command again. If it fails again, contact your administrator.`
)

func NewPollRequest(u map[string][]string) (*PollRequest, error) {
	p := &PollRequest{}
	for key, values := range u {
		switch key {
		case "channel_id":
			if err := checkIdLength(values[0]); err != nil {
				return nil, err
			}
			p.ChannelId = values[0]
		case "token":
			if err := checkIdLength(values[0]); err != nil {
				return nil, err
			}
			p.Token = values[0]
		case "text":
			message, emojis, err := parseText(values[0])
			if err != nil {
				return nil, err
			}
			p.Message, p.Emojis = message, emojis
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
		return "", nil, fmt.Errorf(ErrorTextWrongFormat)
	}
	e := re.FindStringSubmatch(text)
	if len(e) != 3 {
		return "", nil, fmt.Errorf(ErrorTextWrongFormat)
	}
	var emojis []string
	for _, v := range strings.Split(e[2], " ") {
		if len(v) == 0 {
			continue
		}
		if len(v) < 3 || !strings.HasPrefix(v, ":") || !strings.HasSuffix(v, ":") {
			return "", nil, fmt.Errorf(ErrorTextWrongFormat)
		}
		emojis = append(emojis, v[1:len(v)-1])
	}
	return e[1], emojis, nil
}

func checkIdLength(id string) error {
	if len(id) != 26 {
		return fmt.Errorf(ErrorWrongLength)
	}
	return nil
}
