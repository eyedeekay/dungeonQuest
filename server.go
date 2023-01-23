package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func hello(w http.ResponseWriter, req *http.Request) {
	scheme := "https://"
	keys, err := garlic.Keys()
	if err == nil {
		b32 := keys.Address.Base32()
		addr := fmt.Sprintf("%s%s:%s/game/client/index.html", scheme, string(b32), "80")
		http.Redirect(w, req, addr, http.StatusFound)
	}
}

// Handler
func helloFunc(c echo.Context) error {
	scheme := "https://"
	keys, err := garlic.Keys()
	if err == nil {
		b32 := keys.Address.Base32()
		addr := fmt.Sprintf("%s%s:%s/game/client/index.html", scheme, string(b32), "80")
		return c.Redirect(http.StatusFound, addr)
	}
	return nil
}
