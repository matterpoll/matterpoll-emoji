package poll

import (
	"encoding/json"
	"io/ioutil"
)

type PollConf struct {
	Host string
	User PollUser
}

type PollUser struct {
	Id       string
	Password string
}

func LoadConf(path string) (*PollConf, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p PollConf
	json.Unmarshal(b, &p)
	return &p, nil
}
