package gondole

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPrepareRequest(t *testing.T) {
    g, err := NewApp("foo", nil, NoRedirect)
    assert.NoError(t, err, "no error")

    req := g.prepareRequest("bar")
    assert.NotNil(t, req.Headers, "not nil")
    assert.NotNil(t, req.QueryParams, "not nil")
}

