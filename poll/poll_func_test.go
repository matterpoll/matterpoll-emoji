package poll_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kaakaa/matterpoll-emoji/poll"
	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandCorrect(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := "What do you gys wanna grab for lunch?"
	emojis := ":pizza: :sushi:"
	c, err := getConfig("sample_conf.json")
	require.Nil(err)
	ps := poll.Server{Conf: c}

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", c.Token, model.NewId(), message, emojis)
	response, header := sendHttpRequest(require, &ps, payload)

	assert.Equal("application/json", header.Get("Content-Type"))
	assert.Equal(poll.ResponseUsername, response.Username)
	assert.Equal(poll.ResponseIconURL, response.IconURL)
	assert.Equal(model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, response.ResponseType)
	assert.Equal(message+" #poll", response.Text)
}

func TestCommandWronMessageFormat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := model.NewRandomString(20)
	emojis := ""
	c, err := getConfig("sample_conf.json")
	require.Nil(err)
	ps := poll.Server{Conf: c}

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", model.NewId(), model.NewId(), message, emojis)
	response, header := sendHttpRequest(require, &ps, payload)

	assert.Equal("application/json", header.Get("Content-Type"))
	assert.Equal(poll.ResponseUsername, response.Username)
	assert.Equal(poll.ResponseIconURL, response.IconURL)
	assert.Equal(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, response.ResponseType)
	assert.Equal(poll.ErrorTextWrongFormat, response.Text)
}

func TestCommandTokenMissmatch(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := "What do you gys wanna grab for lunch?"
	emojis := ":pizza: :sushi:"
	c, err := getConfig("sample_conf.json")
	require.Nil(err)
	ps := poll.Server{Conf: c}

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", model.NewId(), model.NewId(), message, emojis)
	response, header := sendHttpRequest(require, &ps, payload)

	assert.Equal("application/json", header.Get("Content-Type"))
	assert.Equal(poll.ResponseUsername, response.Username)
	assert.Equal(poll.ResponseIconURL, response.IconURL)
	assert.Equal(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, response.ResponseType)
	assert.Equal(poll.ErrorTokenMissmatch, response.Text)
}

func TestHeaderMediaTypeWrong(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	message := "What do you gys wanna grab for lunch?"
	emojis := ":pizza: :sushi:"
	c, err := getConfig("sample_conf.json")
	require.Nil(err)
	ps := poll.Server{Conf: c}

	payload := fmt.Sprintf("token=%s&channel_id=%s&text=\"%s\"%s", c.Token, model.NewId(), message, emojis)
	reader := strings.NewReader(payload)
	r, err := http.NewRequest("POST", "localhost:8505/poll", reader)
	require.Nil(err)
	require.NotNil(r)

	recorder := httptest.NewRecorder()
	ps.Cmd(recorder, r)
	assert.Equal(http.StatusUnsupportedMediaType, recorder.Code)
}

func TestURLFormat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c, err := getConfig("sample_conf.json")
	require.Nil(err)
	ps := poll.Server{Conf: c}

	payload := "%"
	reader := strings.NewReader(payload)
	r, err := http.NewRequest("POST", "localhost:8505/poll", reader)
	require.Nil(err)
	require.NotNil(r)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	ps.Cmd(recorder, r)
	assert.Equal(http.StatusBadRequest, recorder.Code)
}

func getConfig(path string) (*poll.Conf, error) {
	p, err := getTestFilePath(path)
	if err != nil {
		return nil, err
	}
	c, err := poll.LoadConf(p)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func sendHttpRequest(require *require.Assertions, ps *poll.Server, payload string) (response *model.CommandResponse, header http.Header) {
	reader := strings.NewReader(payload)

	r, err := http.NewRequest(http.MethodPost, "localhost:8505/poll", reader)
	require.Nil(err)
	require.NotNil(r)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	ps.Cmd(recorder, r)
	header = recorder.Header()
	response = model.CommandResponseFromJson(recorder.Result().Body)
	require.NotNil(response)
	return
}
