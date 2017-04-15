package main

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

	_, err = loadGlobal(filepath.Join("test", "garbage.token"))
	assert.Error(t, err, "just garbage")

	// git does now allow you to checkin 000 files :(
	err = os.Chmod(filepath.Join("test", "perms.toml"), 0000)
	assert.NoError(t, err, "should be fine")
	_, err = loadGlobal(filepath.Join("test", "perms.toml"))
	assert.Error(t, err, "unreadable")
	err = os.Chmod(filepath.Join("test", "perms.toml"), 0600)
	assert.NoError(t, err, "should be fine")

	c, err := loadGlobal(filepath.Join("test", "config.toml"))
	assert.NoError(t, err, "should read it fine")
	assert.EqualValues(t, "foo", c.Default, "equal")
}

func TestLoadInstance(t *testing.T) {
	baseDir = "."

	_, err := loadInstance("nonexistent")
	assert.Error(t, err, "does not exist")

	file := filepath.Join("test", "garbage")
	_, err = loadInstance(file)
	assert.Error(t, err, "just garbage")

	file = filepath.Join("test", "foo.token")
	err = os.Chmod(file, 0000)
	assert.NoError(t, err, "should be fine")

	file = filepath.Join("test", "foo")
	_, err = loadInstance(file)
	assert.Error(t, err, "unreadable")

	file = filepath.Join("test", "foo.token")
	err = os.Chmod(file, 0644)
	assert.NoError(t, err, "should be fine")

	real := &Server{
		ID:          "666abcdef666",
		Name:        "foo",
		BearerToken: "d3b07384d113edec49eaa6238ad5ff00",
		APIBase:     "https://mastodon.social/api/v1",
		InstanceURL: "https://mastodon.social",
	}
	file = filepath.Join("test", "foo")
	s, err := loadInstance(file)
	assert.NoError(t, err, "all fine")
	assert.Equal(t, real, s, "equal")
}

func TestGetInstanceList(t *testing.T) {
	baseDir = "test"

	real := []string{
		filepath.Join("test", "foo.token"),
		filepath.Join("test", "garbage.token"),
	}
	list := GetInstanceList()
	assert.Equal(t, real, list, "equal")

	baseDir = "/tmp"
	real = nil
	list = GetInstanceList()
	assert.Equal(t, real, list, "equal")

	baseDir = "/nonexistent"
	real = nil
	list = GetInstanceList()
	assert.Equal(t, real, list, "equal")
}

func TestLoadConfig(t *testing.T) {
	baseDir = "test"

	_, err := LoadConfig("foo")
	assert.NoError(t, err, "should be fine")

	_, err = LoadConfig("")
	assert.NoError(t, err, "should be fine")

	err = os.Chmod(filepath.Join("test", "config.toml"), 0000)
	assert.NoError(t, err, "should be fine")

	_, err = LoadConfig("")
	assert.Error(t, err, "should be unreadable")

	err = os.Chmod(filepath.Join("test", "config.toml"), 0600)
	assert.NoError(t, err, "should be fine")

}
