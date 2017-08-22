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

func TestCommandCorrect(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := "What do you gys wanna grab for lunch?"
	emojis := ":pizza: :sushi:"
	err := setConfig("sample_conf.json")
	require.Nil(err)

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", Conf.Token, model.NewId(), message, emojis)
	response := sendHttpRequest(require, payload)

	assert.Equal(ResponseUsername, response.Username)
	assert.Equal(ResponseIconUrl, response.IconURL)
	assert.Equal(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, response.ResponseType)
	assert.Equal(message+" #poll", response.Text)
}

func TestCommandWronMessageFormat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := model.NewRandomString(20)
	emojis := ""
	err := setConfig("sample_conf.json")
	require.Nil(err)

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", model.NewId(), model.NewId(), message, emojis)
	response := sendHttpRequest(require, payload)

	assert.Equal(ResponseUsername, response.Username)
	assert.Equal(ResponseIconUrl, response.IconURL)
	assert.Equal(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, response.ResponseType)
	assert.Equal(ErrorTextWrongFormat, response.Text)
}

func TestCommandTokenMissmatch(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := "What do you gys wanna grab for lunch?"
	emojis := ":pizza: :sushi:"
	err := setConfig("sample_conf.json")
	require.Nil(err)

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", model.NewId(), model.NewId(), message, emojis)
	response := sendHttpRequest(require, payload)

	assert.Equal(ResponseUsername, response.Username)
	assert.Equal(ResponseIconUrl, response.IconURL)
	assert.Equal(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, response.ResponseType)
	assert.Equal(ErrorTokenMissmatch, response.Text)
}

func TestHeaderMediaTypeWrong(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := "What do you gys wanna grab for lunch?"
	emojis := ":pizza: :sushi:"
	err := setConfig("sample_conf.json")
	require.Nil(err)

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", Conf.Token, model.NewId(), message, emojis)
	reader := strings.NewReader(payload)
	r, err := http.NewRequest("POST", "localhost:8505/poll", reader)
	require.Nil(err)
	require.NotNil(r)

	recorder := httptest.NewRecorder()
	PollCmd(recorder, r)
	assert.Equal(http.StatusUnsupportedMediaType, recorder.Code)
}

func TestURLFormat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	err := setConfig("sample_conf.json")
	require.Nil(err)
	payload := "%"
	reader := strings.NewReader(payload)

	r, err := http.NewRequest("POST", "localhost:8505/poll", reader)
	require.Nil(err)
	require.NotNil(r)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	PollCmd(recorder, r)
	assert.Equal(http.StatusBadRequest, recorder.Code)
}

func setConfig(path string) (err error) {
	p, err := getTestFilePath(path)
	if err != nil {
		return
	}
	c, err := LoadConf(p)
	if err != nil {
		return
	}
	Conf = c
	return nil
}

func sendHttpRequest(require *require.Assertions, payload string) (response *model.CommandResponse) {
	reader := strings.NewReader(payload)

	r, err := http.NewRequest(http.MethodPost, "localhost:8505/poll", reader)
	require.Nil(err)
	require.NotNil(r)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	PollCmd(recorder, r)
	response = model.CommandResponseFromJson(recorder.Result().Body)
	require.NotNil(response)
	return
}
