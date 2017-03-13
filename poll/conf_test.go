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
	e = "bot"
	if c.User.Id != e {
		t.Error("LoadConf Error. Expected: %s, Actual: %s.", e, c.User.Id)
	}
	e = "botbot"
	if c.User.Password != e {
		t.Error("LoadConf Error. Expected: %s, Actual: %s.", e, c.User.Password)
	}
}

func TestReadConfError(t *testing.T) {
	p, err := getTestFilePath("sample_conf_error.json")
	if err != nil {
		t.Fatalf("Cannot get test file path %v", err)
	}
	_, err = LoadConf(p)
	if err == nil {
		t.Fatalf("PollConf validation error")
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
