package poll

import (
	"github.com/mattermost/platform/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPollRequest(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	tests := []struct {
		ChannelId   string
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

			assert.Equal(test.ChannelId, p.ChannelId)
			assert.Equal(test.Token, p.Token)
			assert.Equal(test.Message, p.Message)
			assert.Equal(test.Emojis, p.Emojis)
		}
	}
}
