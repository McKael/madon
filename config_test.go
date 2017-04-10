package gondole

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"os"
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

