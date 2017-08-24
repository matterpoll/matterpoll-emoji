package poll

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type PollConf struct {
	Host   string   `json:"host"`
	Token  string   `json:"token"`
	User   PollUser `json:"user"`
	Listen string   `json:"listen"`
}

type PollUser struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

func LoadConf(path string) (*PollConf, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p PollConf
	json.Unmarshal(b, &p)
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *PollConf) validate() error {
	if len(c.Host) == 0 {
		return fmt.Errorf("Config `host` is missing")
	}
	if len(c.Listen) == 0 {
		return fmt.Errorf("Config `listen` is missing")
	}
	if len(c.Token) == 0 {
		return fmt.Errorf("Config `token` is missing")
	}
	if len(c.Token) != 26 {
		return fmt.Errorf("Invalid token length. Check you config.json")
	}
	if len(c.User.Id) == 0 {
		return fmt.Errorf("Config `user.id` is missing")
	}
	if len(c.User.Password) == 0 {
		return fmt.Errorf("Config `user.password` is missing")
	}
	return nil
}
