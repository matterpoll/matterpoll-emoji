package poll

import (
	"fmt"
	"github.com/mattermost/platform/model"
	"log"
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
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	tests := []Test{
		// All correct
		{"sample_conf.json", "What do you gys wanna grab for lunch?", ":pizza: :sushi:", true},
		// Wrong message format
		{"sample_conf.json", model.NewRandomString(20), "", false},
		// Token missmatch
		{"sample_conf.json", "What do you gys wanna grab for lunch?", ":pizza: :sushi:", true},
	}
	for i, test := range tests {
		err := SetConfig(test.Filename)
		if err != nil {
			t.Fatalf("Test %v: Failed to read config file: %v", i, err)
		}
		var payload string
		switch i {
		// All correct
		case 0:
			payload = fmt.Sprintf("token=%s&user_id=%s&text=\"%s\"%s", Conf.Token, model.NewId(), test.Message, test.Emojis)
		// Wrong message format
		case 1:
			payload = fmt.Sprintf("token=%s&user_id=%s&text=%s", Conf.Token, model.NewId(), test.Message)
		// Token missmatch
		case 2:
			payload = fmt.Sprintf("token=%s&user_id=%s&text=\"%s\"%s", model.NewId(), model.NewId(), test.Message, test.Emojis)
		}
		reader := strings.NewReader(payload)

		r, _ := http.NewRequest(http.MethodPost, "localhost:8505/poll", reader)
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		recorder := httptest.NewRecorder()
		PollCmd(recorder, r)

		response := model.CommandResponseFromJson(recorder.Result().Body)

		if response.Username != RESPONSE_USERNAME {
			t.Errorf("Test %v: Assertion error. Expected: %s, Actual: %s.", i, RESPONSE_USERNAME, response.Username)
		}
		if response.IconURL != RESPONSE_ICON_URL {
			t.Errorf("Test %v: Assertion error. Expected: %s, Actual: %s.", i, RESPONSE_ICON_URL, response.IconURL)
		}
		if test.CorrectPoll {
			if response.ResponseType != model.COMMAND_RESPONSE_TYPE_IN_CHANNEL {
				t.Errorf("Test %v: Assertion error. Expected: %s, Actual: %s.", i, model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, response.Text)
			}
			if response.Text != test.Message+" #poll" {
				t.Errorf("Test %v: Assertion error. Expected: %s, Actual: %s.", i, test.Message, response.Text)
			}
		} else {
			if response.ResponseType != model.COMMAND_RESPONSE_TYPE_EPHEMERAL {
				t.Errorf("Test %v: Assertion error. Expected: %s, Actual: %s.", i, model.COMMAND_RESPONSE_TYPE_EPHEMERAL, response.ResponseType)
			}
			if response.Text != error_wrong_format {
				t.Errorf("Test %v: Assertion error. Expected: %s, Actual: %s.", i, error_wrong_format, response.Text)
			}
		}
	}
}

func TestHeader(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	tests := []Test{
		// All correct
		{"sample_conf.json", "What do you gys wanna grab for lunch?", ":pizza: :sushi:", true},
	}
	for i, test := range tests {
		err := SetConfig(test.Filename)
		if err != nil {
			t.Fatalf("Test %v: Failed to read config file: %v", i, err)
		}
		payload := fmt.Sprintf("token=%s&user_id=%s&text=\"%s\"%s", Conf.Token, model.NewId(), test.Message, test.Emojis)
		reader := strings.NewReader(payload)
		switch i {
		case 0:
			r, _ := http.NewRequest("POST", "localhost:8505/poll", reader)
			recorder := httptest.NewRecorder()
			PollCmd(recorder, r)
			if recorder.Code != http.StatusUnsupportedMediaType {
				t.Errorf("Test %v: Assertion error. Expected: %v, Actual: %v.", i, http.StatusNotAcceptable, recorder.Code)
			}
		}
	}
}

func TestURLFormat(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	i := 0
	err := SetConfig("sample_conf.json")
	if err != nil {
		t.Fatalf("Test %v: Failed to read config file: %v", i, err)
	}
	payload := fmt.Sprintf("%")
	reader := strings.NewReader(payload)

	r, _ := http.NewRequest("POST", "localhost:8505/poll", reader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	PollCmd(recorder, r)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Test %v: Assertion error. Expected: %v, Actual: %v.", i, http.StatusBadRequest, recorder.Code)
	}
}

func SetConfig(path string) (err error) {
	p, err := getTestFilePath(path)
	if err != nil {
		return
	}
	c, err := LoadConf(p)
	if err != nil {
		return
	}
	Conf = c
	return
}
