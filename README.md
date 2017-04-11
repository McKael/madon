# gondole

Go version of the Mastodon API

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/keltia/gondole) [![license](https://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/keltia/gondole/master/LICENSE) [![build](https://img.shields.io/travis/keltia/gondole.svg?style=flat)](https://travis-ci.org/keltia/gondole) [![Go Report Card](https://goreportcard.com/badge/github.com/keltia/gondole)](https://goreportcard.com/report/github.com/keltia/gondole)

`gondole` is a [Go](https://golang.org/) library to access the Mastondon [REST API](http://www.rubydoc.info/gems/mastodon-api/Mastodon/REST/API).

**Work in progress, still incomplete**

## Installation

Like many Go-based tools, installation is very easy
  
    go get github.com/keltia/gondole/cmd/...

  or
  
    git clone https://github.com/keltia/gondole
    go install ./cmd/...

The library is fetched, compiled and installed in whichever directory is specified by `$GOPATH`.  The `gondole-cli` binary will also be installed. 

## Configuration

The configuration is stored in the standard place: `$HOME/.config/gondole`.  You can find a global configuration file named `config.toml`and a file for each instance you have connected with, named `<instance>.toml`.

When you connect `gondole` to a given instance you have never connected before, it will register the application, get the token and save it.

## Name

Trying to define a name close to *Mastodon*, one could come up with *godon* and in French, *gondole* (the small boats in Venice) is easy to take.

## References

- [Mastodon API doc](https://github.com/tootsuite/mastodon/blob/master/docs/Using-the-API/API.md)
- [Mastodon Ruby API](http://www.rubydoc.info/gems/mastodon-api/Mastodon/REST/API)
- [Mastodon Python](https://mastodonpy.readthedocs.io/en/latest/)
- [Mastodon repo](https://github.com/tootsuite/mastodon)
