package poll

import (
	"fmt"
	"regexp"
	"strings"
)

// Request wraps up all information needed to answer a poll request
type Request struct {
	ChannelID string
	Token     string
	Message   string
	Emojis    []string
}

const (
	backTick = "`"
	//ErrorTextWrongFormat is an error message and is used, if the the message isn`t formated correct
	ErrorTextWrongFormat = `The message format is wrong. Try this instead: ` + backTick + `/poll \"What do you gys wanna grab for lunch?\" :pizza: :sushi:` + backTick
	// ErrorTokenMissmatch is an error message and is used, if the token comparison fails
	ErrorTokenMissmatch = `An error occurred. Ask your administrator to check the Matterpoll config settings.`
	// ErrorWrongLength is an error message and is used, if the channel id or the token have a wrong length
	ErrorWrongLength = `An error occurred. Try the same command again. If it fails again, contact your administrator.`
)

// NewRequest validates the data in map and wraps it into a Request struct
func NewRequest(u map[string][]string) (*Request, error) {
	p := &Request{}
	for key, values := range u {
		switch key {
		case "channel_id":
			if err := checkIDLength(values[0]); err != nil {
				return nil, err
			}
			p.ChannelID = values[0]
		case "token":
			if err := checkIDLength(values[0]); err != nil {
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

func checkIDLength(id string) error {
	if len(id) != 26 {
		return fmt.Errorf(ErrorWrongLength)
	}
	return nil
}
