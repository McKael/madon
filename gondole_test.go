package gondole

import (
	"testing"

	"github.com/sendgrid/rest"
	"github.com/stretchr/testify/assert"
)

func TestPrepareRequest(t *testing.T) {
	g := &Client{
		Name:    "foo",
		ID:      "666",
		Secret:  "biiiip",
		APIBase: "http://example.com",
	}

	req := g.prepareRequest("bar", rest.Get, nil)
	assert.NotNil(t, req.Headers, "not nil")
}
