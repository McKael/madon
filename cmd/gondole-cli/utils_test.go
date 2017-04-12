package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFilterURL(t *testing.T) {
    in := "https://example.com"
    out := filterURL(in)
    assert.EqualValues(t, in, out, "equal")
}

func TestBasename(t *testing.T) {
    in := "https://example.com"
    out := basename(in)
    assert.EqualValues(t, "example.com", out, "equal")

    in = "https://example.com:80"
    out = basename(in)
    assert.EqualValues(t, "example.com", out, "equal")

    in = "https://example.com:16443"
    out = basename(in)
    assert.EqualValues(t, "example.com", out, "equal")

    in = "//example.com:443"
    out = basename(in)
    assert.EqualValues(t, "example.com", out, "equal")
}
