package poll

import (
	"fmt"
	"testing"
)

func TestNewPollRequest(t *testing.T) {
	tables := []struct {
		teamId    string
		channelId string
		text      string
		message   string
		emojis    []string
		err       error
	}{
		{"teamidxxxx", "channelidxxxx", "`description` :emoji1:", "description", []string{"emoji1"}, nil},
		{"teamidxxxx", "channelidxxxx", "'description' :emoji1:", "description", []string{"emoji1"}, nil},
		{"teamidxxxx", "channelidxxxx", "\"description\" :emoji1:", "description", []string{"emoji1"}, nil},

		{"teamidxxxx", "channelidxxxx", "`description including space` :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},
		{"teamidxxxx", "channelidxxxx", "'description including space' :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},
		{"teamidxxxx", "channelidxxxx", "\"description including space\" :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},

		{"teamidxxxx", "channelidxxxx", "`description including space` :emoji1:  :emoji2:   :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},

		{"teamidxxxx", "channelidxxxx", "description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("Start quotation mark missing")},
		{"teamidxxxx", "channelidxxxx", "`description including space :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("End quotation mark missing")},

		{"teamidxxxx", "channelidxxxx", "`description including space' :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("First and second quotes do not match")},
		{"teamidxxxx", "channelidxxxx", "'description including space\" :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("First and second quotes do not match")},
		{"teamidxxxx", "channelidxxxx", "\"description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, fmt.Errorf("First and second quotes do not match")},

		{"", "channelidxxxx", "`description` :emoji1:", "", []string{}, fmt.Errorf("Unexpected Error: TeamID in request is empty.")},
		{"teamidxxxx", "", "`description` :emoji1:", "", []string{}, fmt.Errorf("Unexpected Error: ChannelID in request is empty.")},

		{"teamidxxxx", "channelidxxxx", "`description including space`", "", []string{}, fmt.Errorf("No emoji found")},
		{"teamidxxxx", "channelidxxxx", "`description including space` emoji1", "", []string{}, fmt.Errorf("Emoji format error")},
		{"teamidxxxx", "channelidxxxx", "`` :emoji1:", "", []string{""}, fmt.Errorf("No description found")},
	}

	for _, tt := range tables {
		s := fmt.Sprintf("team_id=%s&channel_id=%s&text=%s", tt.teamId, tt.channelId, tt.text)

		p, err := NewPollRequest(s)

		if err != nil && tt.err == nil {
			t.Fatalf("Test retured with error : %v but there sould be none", err)
		}
		if err == nil && tt.err != nil {
			t.Fatalf("Test didnt not return with an error but it should return with %v", tt.err)
		}
		if err != nil && tt.err != nil {
			continue
		}

		if p.TeamId != tt.teamId {
			t.Errorf("Assertion error `TeamId`. Expected: %s, Actual: %s.", tt.teamId, p.TeamId)
		}
		if p.ChannelId != tt.channelId {
			t.Errorf("Assertion error `ChannelId`. Expected: %s, Actual: %s.", tt.channelId, p.ChannelId)
		}
		if p.Message != tt.message {
			t.Errorf("Assertion error `Message`. Expected: %s, Actual: %s.", tt.message, p.Message)
		}
		if len(p.Emojis) != len(tt.emojis) {
			t.Errorf("Assertion error `Emojis`. Expected: %s, Actual: %s.", tt.emojis, p.Emojis)
		}
		for i, v := range p.Emojis {
			if v != tt.emojis[i] {
				t.Errorf("Assertion error `Emojis`. Expected: %s, Actual: %s.", tt.emojis, p.Emojis)
			}
		}
	}
}
