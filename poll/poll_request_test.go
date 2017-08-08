package poll

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPollRequest(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	tests := []struct {
		TeamId      string
		ChannelId   string
		Token       string
		Text        string
		Message     string
		Emojis      []string
		ShouldError bool
	}{
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description` :emoji1:", "description", []string{"emoji1"}, false},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "'description' :emoji1:", "description", []string{"emoji1"}, false},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "\"description\" :emoji1:", "description", []string{"emoji1"}, false},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space` :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "'description including space' :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "\"description including space\" :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space` :emoji1:  :emoji2:   :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, true},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space :emoji1: :emoji2: :emoji3:", "", []string{""}, true},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space' :emoji1: :emoji2: :emoji3:", "", []string{""}, true},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "'description including space\" :emoji1: :emoji2: :emoji3:", "", []string{""}, true},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "\"description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, true},

		{"", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description` :emoji1:", "", []string{}, true},
		{"teamidxxxx", "", "9jrxak1ykxrmnaed9cps9i4cim", "`description` :emoji1:", "", []string{}, true},
		{"teamidxxxx", "channelidxxxx", "", "`description` :emoji1:", "", []string{}, true},

		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space`", "", []string{}, true},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`description including space` emoji1", "", []string{}, true},
		{"teamidxxxx", "channelidxxxx", "9jrxak1ykxrmnaed9cps9i4cim", "`` :emoji1:", "", []string{""}, true},
	}

	for _, test := range tests {
		s := make(map[string][]string)
		s["team_id"] = []string{test.TeamId}
		s["channel_id"] = []string{test.ChannelId}
		s["token"] = []string{test.Token}
		s["text"] = []string{test.Text}

		p, err := NewPollRequest(s)
		if test.ShouldError == true {
			assert.NotNil(err)
			assert.Nil(p)
		} else {
			assert.Nil(err)
			require.NotNil(p)

			assert.Equal(p.TeamId, test.TeamId)
			assert.Equal(p.ChannelId, test.ChannelId)
			assert.Equal(p.Token, test.Token)
			assert.Equal(p.Message, test.Message)
			assert.Equal(p.Emojis, test.Emojis)
		}
	}
}
