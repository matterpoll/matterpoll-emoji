package poll_test

import (
	"testing"

	"github.com/kaakaa/matterpoll-emoji/poll"
	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPollRequest(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	tests := []struct {
		ChannelID   string
		Token       string
		Text        string
		Message     string
		Emojis      []string
		ShouldError bool
	}{
		{model.NewId(), model.NewId(), "`description` :emoji1:", "description", []string{"emoji1"}, false},
		{model.NewId(), model.NewId(), "'description' :emoji1:", "description", []string{"emoji1"}, false},
		{model.NewId(), model.NewId(), "\"description\" :emoji1:", "description", []string{"emoji1"}, false},

		{model.NewId(), model.NewId(), "`description including space` :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},
		{model.NewId(), model.NewId(), "'description including space' :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},
		{model.NewId(), model.NewId(), "\"description including space\" :emoji1: :emoji2: :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},

		{model.NewId(), model.NewId(), "`description including space` :emoji1:  :emoji2:   :emoji3:", "description including space", []string{"emoji1", "emoji2", "emoji3"}, false},

		{model.NewId(), model.NewId(), "description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, true},
		{model.NewId(), model.NewId(), "`description including space :emoji1: :emoji2: :emoji3:", "", []string{""}, true},

		{model.NewId(), model.NewId(), "`description including space' :emoji1: :emoji2: :emoji3:", "", []string{""}, true},
		{model.NewId(), model.NewId(), "'description including space\" :emoji1: :emoji2: :emoji3:", "", []string{""}, true},
		{model.NewId(), model.NewId(), "\"description including space` :emoji1: :emoji2: :emoji3:", "", []string{""}, true},

		{"", model.NewId(), "`description` :emoji1:", "", []string{}, true},
		{model.NewId(), "", "`description` :emoji1:", "", []string{}, true},

		{model.NewId(), model.NewId(), "`description including space`", "", []string{}, true},
		{model.NewId(), model.NewId(), "`description including space` emoji1", "", []string{}, true},
		{model.NewId(), model.NewId(), "`` :emoji1:", "", []string{""}, true},
	}

	for _, test := range tests {
		s := make(map[string][]string)
		s["channel_id"] = []string{test.ChannelID}
		s["token"] = []string{test.Token}
		s["text"] = []string{test.Text}

		p, err := poll.NewRequest(s)
		if test.ShouldError {
			assert.NotNil(err)
			assert.Nil(p)
		} else {
			assert.Nil(err)
			require.NotNil(p)

			assert.Equal(p.ChannelID, test.ChannelID)
			assert.Equal(p.Token, test.Token)
			assert.Equal(p.Message, test.Message)
			assert.Equal(p.Emojis, test.Emojis)
		}
	}
}
