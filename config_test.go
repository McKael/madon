package gondole

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadGlobal(t *testing.T) {
	baseDir = "."

	_, err := loadGlobal(filepath.Join("test", "non.toml"))
	assert.Error(t, err, "does not exist")

	// git does now allow you to checkin 000 files :(
	err = os.Chmod(filepath.Join("test", "perms.toml"), 000)
	_, err = loadGlobal(filepath.Join("test", "perms.toml"))
	assert.Error(t, err, "unreadable")
	err = os.Chmod(filepath.Join("test", "perms.toml"), 600)

	c, err := loadGlobal(filepath.Join("test", "config.toml"))
	assert.NoError(t, err, "should read it fine")
	assert.EqualValues(t, "foobar", c.Default, "equal")
}

func TestLoadInstance(t *testing.T) {
	baseDir = "."

	_, err := loadInstance("nonexistent")
	assert.Error(t, err, "does not exist")

	real := &Server{
		ID:          666,
		Name:        "foo",
		BearerToken: "d3b07384d113edec49eaa6238ad5ff00",
	}
	s, err := loadInstance("test/foo")
	assert.NoError(t, err, "all fine")
	assert.Equal(t, real, s, "equal")
}

func TestGetInstanceList(t *testing.T) {
	baseDir = "test"

	real := []string{"test/foo.token"}
	list := GetInstanceList()
	assert.Equal(t, real, list, "equal")
}
