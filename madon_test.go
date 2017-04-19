package madon

import (
	"testing"

	"github.com/sendgrid/rest"
	"github.com/stretchr/testify/assert"
)

func TestPrepareRequest(t *testing.T) {
	mc := &Client{
		Name:    "foo",
		ID:      "666",
		Secret:  "biiiip",
		APIBase: "http://example.com",
	}

	req, err := mc.prepareRequest("bar", rest.Get, nil)
	assert.NoError(t, err, "no error")
	assert.NotNil(t, req.Headers, "not nil")
}
