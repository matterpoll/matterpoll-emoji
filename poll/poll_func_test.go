package poll

import (
	"fmt"
	"github.com/mattermost/platform/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Test struct {
	Filename    string
	Message     string
	Emojis      string
	CorrectPoll bool
}

func TestCommand(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []Test{
		// All correct
		{"sample_conf.json", "What do you gys wanna grab for lunch?", ":pizza: :sushi:", true},
		// Wrong message format
		{"sample_conf.json", model.NewRandomString(20), "", false},
		// Token missmatch
		{"sample_conf.json", "What do you gys wanna grab for lunch?", ":pizza: :sushi:", true},
	}
	for i, test := range tests {
		c, err := getConfig(test.Filename)
		require.Nil(err)
		ps := PollServer{c}
		var payload string
		switch i {
		// All correct
		case 0:
			payload = fmt.Sprintf("token=%s&user_id=%s&text=\"%s\"%s", c.Token, model.NewId(), test.Message, test.Emojis)
		// Wrong message format
		case 1:
			payload = fmt.Sprintf("token=%s&user_id=%s&text=%s", c.Token, model.NewId(), test.Message)
		// Token missmatch
		case 2:
			payload = fmt.Sprintf("token=%s&user_id=%s&text=\"%s\"%s", model.NewId(), model.NewId(), test.Message, test.Emojis)
		}
		reader := strings.NewReader(payload)

		r, err := http.NewRequest(http.MethodPost, "localhost:8505/poll", reader)
		require.Nil(err)
		require.NotNil(r)
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		recorder := httptest.NewRecorder()
		ps.PollCmd(recorder, r)

		response := model.CommandResponseFromJson(recorder.Result().Body)
		require.NotNil(response)
		assert.Equal(response.Username, RESPONSE_USERNAME)
		assert.Equal(response.IconURL, RESPONSE_ICON_URL)
		if test.CorrectPoll {
			assert.Equal(response.ResponseType, model.COMMAND_RESPONSE_TYPE_IN_CHANNEL)
			assert.Equal(response.Text, test.Message+" #poll")
		} else {
			assert.Equal(response.ResponseType, model.COMMAND_RESPONSE_TYPE_EPHEMERAL)
			assert.Equal(response.Text, error_wrong_format)
		}
	}
}

func TestHeader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []Test{
		// All correct
		{"sample_conf.json", "What do you gys wanna grab for lunch?", ":pizza: :sushi:", true},
	}
	for i, test := range tests {
		c, err := getConfig(test.Filename)
		require.Nil(err)
		ps := PollServer{c}
		payload := fmt.Sprintf("token=%s&user_id=%s&text=\"%s\"%s", c.Token, model.NewId(), test.Message, test.Emojis)
		reader := strings.NewReader(payload)
		switch i {
		case 0:
			r, err := http.NewRequest("POST", "localhost:8505/poll", reader)
			require.Nil(err)
			require.NotNil(r)

			recorder := httptest.NewRecorder()
			ps.PollCmd(recorder, r)
			assert.Equal(recorder.Code, http.StatusUnsupportedMediaType)
		}
	}
}

func TestURLFormat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c, err := getConfig("sample_conf.json")
	require.Nil(err)
	ps := PollServer{c}

	payload := "%"
	reader := strings.NewReader(payload)

	r, err := http.NewRequest("POST", "localhost:8505/poll", reader)
	require.Nil(err)
	require.NotNil(r)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	ps.PollCmd(recorder, r)
	assert.Equal(recorder.Code, http.StatusBadRequest)
}

func getConfig(path string) (*PollConf, error) {
	p, err := getTestFilePath(path)
	if err != nil {
		return nil, err
	}
	conf, err := LoadConf(p)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
