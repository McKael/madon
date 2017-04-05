package gondole

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewApp(t *testing.T) {
	g, err := NewApp("foo", "bar")
	assert.NoError(t, err, "no error")
	assert.Equal(t, reflect.TypeOf(&Gondole{}), reflect.TypeOf(g), "should be Gondole")

	assert.Equal(t, "foo", g.Name, "should be equal")
	assert.Equal(t, "bar", g.RedirectURI, "should be equal")
}
