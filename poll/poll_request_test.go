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
		{"teamidxxxx", "channelidxxxx", "`description including space` :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, nil},
		{"", "channelidxxxx", "`description` :emoji1:", "", []string{}, fmt.Errorf("Unexpected Error: TeamID in request is empty.")},
		{"teamidxxxx", "", "`description` :emoji1:", "", []string{}, fmt.Errorf("Unexpected Error: ChannelID in request is empty.")},
		{"teamidxxxx", "channelidxxxx", "`` :emoji1:", "", []string{""}, fmt.Errorf("Command Error: /poll `Here is description` :thumbsup: :thumbsdown:...")},
		{"teamidxxxx", "channelidxxxx", "`description`", "", []string{""}, fmt.Errorf("Command Error: /poll `Here is description` :thumbsup: :thumbsdown:...")},
	}

	for _, tt := range tables {
		s := fmt.Sprintf("team_id=%s&channel_id=%s&text=%s", tt.teamId, tt.channelId, tt.text)

		p, err := NewPollRequest(s)
		if err != tt.err {
			if err != nil && tt.err != nil && err.Error() == tt.err.Error() {
				continue
			} else {
				t.Fatalf("Unexpected error. Expected: %v, Actual: %v", tt.err, err)
			}
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
