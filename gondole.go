package gondole

import (

)

const (
    Version = "0.0"

    APIEndpoint = "/api/v1"
)

func NewApp(name, redirectURI string) (g *Gondole, err error) {
    g = &Gondole{
        Name: name,
        RedirectURI: redirectURI,
    }
    return
}
