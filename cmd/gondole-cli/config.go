// config.go
//
// This file implements the configuration part for when you need the API
// key to modify things in the Mastodon configuration and manage measurements.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/naoina/toml"
)

/*
Assume the application is registered if $HOME/.config/<gondole>/config.toml already exist
We will store the per-instance token into $HOME/.config/<gondole>/<site>.token
*/

const (
	DefaultName = "config.toml"
)

var (
	baseDir = filepath.Join(os.Getenv("HOME"),
		".config",
		"gondole",
	)
)

func loadGlobal(file string) (c *Config, err error) {
	log.Printf("file=%s", file)
	// Check if there is any config file
	_, err = os.Stat(file)
	if err != nil {
		return
	}

	log.Printf("file=%s, found it", file)
	// Read it
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return c, fmt.Errorf("Can not read %s", file)
	}

	cnf := Config{}
	err = toml.Unmarshal(buf, &cnf)
	if err != nil {
		return c, fmt.Errorf("Error parsing toml %s: %v", file, err)
	}
	c = &cnf
	return
}

func loadInstance(name string) (s *Server, err error) {
	// Load instance-specific file
	file := filepath.Join(baseDir, name+".token")

	log.Printf("instance is %s", file)

	// Check if there is any config file
	if _, err = os.Stat(file); err == nil {
		// Read it
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return s, fmt.Errorf("Can not read %s", file)
		}

		sc := Server{}
		err = toml.Unmarshal(buf, &sc)
		if err != nil {
			return s, fmt.Errorf("Error parsing toml %s: %v", file, err)
		}
		s = &sc
	}
	return
}

func GetInstanceList() (list []string) {
	list, err := filepath.Glob(filepath.Join(baseDir, "*.token"))
	log.Printf("basedir=%s", filepath.Join(baseDir, "*.token"))
	if err != nil {
		log.Printf("warning, no *.token files in %s", baseDir)
		list = nil
	}
	log.Printf("list=%v", list)
	return
}

// LoadConfig reads a file as a TOML document and return the structure
func LoadConfig(name string) (s *Server, err error) {
	// Load global file
	gFile := filepath.Join(baseDir, DefaultName)

	log.Printf("global is %s", gFile)
	c, err := loadGlobal(gFile)
	if err != nil {
		return
	}
	if name == "" {
		s, err = loadInstance(c.Default)
	} else {
		s, err = loadInstance(name)
	}

	return s, err
}

func (c *Config) Write() (err error) {
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		log.Fatalf("error creating configuration directory: %s", baseDir)
	}

	var sc []byte

	if sc, err = toml.Marshal(*c); err != nil {
		log.Fatalf("error saving configuration")
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, DefaultName), sc, 0600)
	return
}

func (s *Server) WriteToken(instance string) (err error) {
	if err = os.MkdirAll(baseDir, 0700); err != nil {
		log.Fatalf("error creating configuration directory: %s", baseDir)
	}

	var sc []byte

	if sc, err = toml.Marshal(s); err != nil {
		log.Fatalf("error saving configuration")
	}

	full := instance + ".token"
	err = ioutil.WriteFile(filepath.Join(baseDir, full), sc, 0600)
	return
}
