package gondole

import (

)

// NewApp registers a new instance
func NewApp(name, redirectURI string) (g *Gondole, err error) {
	// Load configuration, will register if none is found

	g = &Gondole{
		Name:   name,
		//ID:     config.ID,
		//Secret: config.BearerToken,
	}
	return
}

