package poll

import (
	"fmt"
	"testing"
)

func TestNewPollRequest(t *testing.T) {
	tables := []struct {
		teamId    string
		channelId string
		token     string
		text      string
		message   string
		emojis    []string
		err       error
	}{
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description` :emoji1:", "description", []string{"emoji1"}, nil},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "'description' :emoji1:", "description", []string{"emoji1"}, nil},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "\"description\" :emoji1:", "description", []string{"emoji1"}, nil},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space` :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "'description including space' :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "\"description including space\" :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space` :emoji1:  :emoji2:   :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("Start quotation mark missing")},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("End quotation mark missing")},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space' :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("First and second quotes do not match")},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "'description including space\" :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("First and second quotes do not match")},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "\"description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("First and second quotes do not match")},

		{"", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description` :emoji1:", "", []string{}, fmt.Errorf("Unexpected Error: TeamID in request is empty.")},
		{"teamidxxxx", "", "9jrxak1ykxrmnaed9cps9i4cim", "`description` :emoji1:", "", []string{}, fmt.Errorf("Unexpected Error: ChannelID in request is empty.")},
    {"teamidxxxx", "channelidxxxx", "", "`description` :emoji1:", "", []string{}, fmt.Errorf("Unexpected Error: Token in request is empty.")},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space`", "", []string{}, fmt.Errorf("No emoji found")},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space` emoji1", "", []string{}, fmt.Errorf("Emoji format error")},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`` :emoji1:", "", []string{""}, fmt.Errorf("No description found")},
	}

	for i, tt := range tables {
		s := fmt.Sprintf("team_id=%s&channel_id=%s&token=%s&text=%s", tt.teamId, tt.channelId, tt.token, tt.text)
		p, err := NewPollRequest(s)

		if err != nil && tt.err == nil {
			t.Fatalf("Test %v retured with error : %v but there sould be none", i, err)
		}
		if err == nil && tt.err != nil {
			t.Fatalf("Test %v didnt not return with an error but it should return with %v", i, tt.err)
		}
		if err != nil && tt.err != nil {
			continue
		}

		if p.TeamId != tt.teamId {
			t.Errorf("Test %v: Assertion error `TeamId` in Expected: %s, Actual: %s.", i, tt.teamId, p.TeamId)
		}
		if p.ChannelId != tt.channelId {
			t.Errorf("Test %v: Assertion error `ChannelId`. Expected: %s, Actual: %s.", i, tt.channelId, p.ChannelId)
		}
    if p.Token != tt.token {
			t.Errorf("Test %v: Assertion error `Token`. Expected: %s, Actual: %s.", i, tt.token, p.Token)
		}
		if p.Message != tt.message {
			t.Errorf("Test %v: Assertion error `Message`. Expected: %s, Actual: %s.", i, tt.message, p.Message)
		}
		if len(p.Emojis) != len(tt.emojis) {
			t.Errorf("Test %v: Assertion error `Emojis`. Expected: %s, Actual: %s.", i, tt.emojis, p.Emojis)
		}
		for i, v := range p.Emojis {
			if v != tt.emojis[i] {
				t.Errorf("Test %v: Assertion error `Emojis`. Expected: %s, Actual: %s.", i, tt.emojis, p.Emojis)
			}
		}
	}
}
