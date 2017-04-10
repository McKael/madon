// config.go
//
// This file implements the configuration part for when you need the API
// key to modify things in the Mastodon configuration and manage measurements.

package gondole

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

// Config holds our parameters
type Server struct {
	ID          int64
	Name        string
	BearerToken string
}

type Config struct {
	Default string
}

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

	err = toml.Unmarshal(buf, c)
	if err != nil {
		return c, fmt.Errorf("Error parsing toml %s: %v", file, err)
	}
	return
}

func loadInstance(name string) (s *Server, err error) {
	// Load instance-specific file
	file := filepath.Join(baseDir, name+".token")

	log.Printf("instance is %s", file)

	// Check if there is any config file
	if _, err := os.Stat(file); err == nil {
		// Read it
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return s, fmt.Errorf("Can not read %s", file)
		}

		err = toml.Unmarshal(buf, s)
		if err != nil {
			return s, fmt.Errorf("Error parsing toml %s: %v", file, err)
		}
	}
	return
}

func GetInstanceList() (list []string) {
	list, err := filepath.Glob(filepath.Join(baseDir + "*.token"))
	if err != nil {
		log.Printf("warning, no *.token files")
		list = nil
	}
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
