package poll

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConf(t *testing.T) {
	p, err := getTestFilePath("sample_conf.json")
	if err != nil {
		t.Fatalf("Cannot get test file path %v", err)
	}
	c, err := LoadConf(p)
	if err != nil {
		t.Fatalf("Cannot load conf file: %v", err)
	}

	e := "http://localhost:8065"
	if c.Host != e {
		t.Error("LoadConf Error. Expected: %s, Actual: %s.", e, c.Host)
	}
	e = "9jrxak1ykxrmnaed9cps9i4cim"
	if c.Token != e {
		t.Error("LoadConf Error. Expected: %s, Actual: %s.", e, c.Token)
	}
	e = "bot"
	if c.User.Id != e {
		t.Error("LoadConf Error. Expected: %s, Actual: %s.", e, c.User.Id)
	}
	e = "botbot"
	if c.User.Password != e {
		t.Error("LoadConf Error. Expected: %s, Actual: %s.", e, c.User.Password)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		filename     string
		should_error bool
	}{
		{"sample_conf.json", false},
		{"sample_conf_error.json", true},
		{"sample_conf_error_no_host.json", true},
		{"sample_conf_error_no_token.json", false},
		{"sample_conf_error_wrong_token_length.json", true},
		{"sample_conf_error_no_user.json", true},
		{"sample_conf_error_no_user_id.json", true},
		{"sample_conf_error_no_user_password.json", true},
	}
	for i, test := range tests {
		p, err := getTestFilePath(test.filename)
		if err != nil {
			t.Fatalf("Test %v: Cannot get test file: %v", i, err)
		}
		_, err = LoadConf(p)
		if err != nil && test.should_error == false {
			t.Fatalf("Test %v: Test retured with error %v but there sould be none", i, err)
		}
		if err == nil && test.should_error == true {
			t.Fatalf("Test %v: Test didn't return with an error but it should return with one", i)
		}
	}
}

func TestReadConfNotExistsError(t *testing.T) {
	p, err := getTestFilePath("not_exists.json")
	if err != nil {
		t.Fatalf("Cannot get test file path %v", err)
	}
	_, err = LoadConf(p)
	if err == nil {
		t.Fatalf("Unexpected Error")
	}
}

func getTestFilePath(path string) (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "..", "testdata", path), nil
}
