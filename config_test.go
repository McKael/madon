package gondole

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"path/filepath"
)

func TestLoadGlobal(t *testing.T) {
	baseDir = "."

	_, err := loadGlobal(filepath.Join("test", "non.toml"))
	assert.Error(t, err, "does not exist")

	_, err = loadGlobal(filepath.Join("test", "perms.toml"))
	assert.Error(t, err, "unreadable")

	c, err := loadGlobal(filepath.Join("test", "config.toml"))
	assert.NoError(t, err, "should read it fine")
	assert.EqualValues(t, "foobar", c.Default, "equal")
}

