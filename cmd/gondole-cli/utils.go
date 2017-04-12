package main

import (
	"github.com/urfave/cli"
	"net/url"
	"strings"
)

// ByAlphabet is for sorting
type ByAlphabet []cli.Command

func (a ByAlphabet) Len() int           { return len(a) }
func (a ByAlphabet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAlphabet) Less(i, j int) bool { return a[i].Name < a[j].Name }

func filterURL(in string) (out string) {
	uri, err := url.Parse(in)
	if err != nil {
		out = ""
	} else {
		uri := url.URL{Scheme: uri.Scheme, Host: uri.Host}
		out = uri.String()
	}
	return
}

func basename(in string) (out string) {
	uri, err := url.Parse(in)
	if err != nil {
		out = ""
	} else {
		// Remove the :NN part of present
		out = strings.Split(uri.Host, ":")[0]
	}
	return
}
