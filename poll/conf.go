package poll

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Conf represents the login credentials of a mattermost user
type Conf struct {
	Host   string `json:"host"`
	Listen string `json:"listen"`
	Token  string `json:"token"`
	User   User   `json:"user"`
}

// User represents the login credentials of a mattermost user
type User struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// LoadConf loads a configuration file located at path and parse it to a Conf struct
func LoadConf(path string) (*Conf, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p Conf
	json.Unmarshal(b, &p)
	if err := p.validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *Conf) validate() error {
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
	if len(c.User.ID) == 0 {
		return fmt.Errorf("Config `user.id` is missing")
	}
	if len(c.User.Password) == 0 {
		return fmt.Errorf("Config `user.password` is missing")
	}
	return nil
}
