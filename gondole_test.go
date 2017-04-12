package gondole

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPrepareRequest(t *testing.T) {
    g := &Client{
        Name: "foo",
        ID: "666",
        Secret: "biiiip",
        APIBase: "http://example.com",
    }

    req := g.prepareRequest("bar")
    assert.NotNil(t, req.Headers, "not nil")
    assert.NotNil(t, req.QueryParams, "not nil")
}

