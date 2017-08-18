package poll_test

import (
	"github.com/kaakaa/matterpoll-emoji/poll"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestReadConf(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	p, err := getTestFilePath("sample_conf.json")
	assert.Nil(err)
	require.NotNil(p)

	c, err := poll.LoadConf(p)
	assert.Nil(err)
	require.NotNil(c)

	assert.Equal(c.Host, "http://localhost:8065")
	assert.Equal(c.Token, "9jrxak1ykxrmnaed9cps9i4cim")
	assert.Equal(c.User.ID, "bot")
	assert.Equal(c.User.Password, "botbot")
}

func TestValidate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	tests := []struct {
		Filename    string
		ShouldError bool
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
	for _, test := range tests {
		p, err := getTestFilePath(test.Filename)
		assert.Nil(err)
		require.NotNil(p)

		c, err := poll.LoadConf(p)
		if test.ShouldError == true {
			assert.NotNil(err)
			assert.Nil(c)
		} else {
			assert.Nil(err)
			assert.NotNil(c)
		}
	}
}

func TestReadConfNotExistsError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	p, err := getTestFilePath("not_exists.json")
	assert.Nil(err)
	require.NotNil(p)

	c, err := poll.LoadConf(p)
	assert.NotNil(err)
	assert.Nil(c)
}

func getTestFilePath(path string) (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "..", "testdata", path), nil
}
