package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func hello(w http.ResponseWriter, req *http.Request) {
	scheme := "https://"
	b32 := garlic.ServiceKeys.Address.Base32()
	addr := fmt.Sprintf("%s%s:8000/game/client/index.html", scheme, string(b32))
	http.Redirect(w, req, addr, http.StatusFound)
}

// Handler
func helloFunc(c echo.Context) error {
	scheme := "https://"
	b32 := garlic.ServiceKeys.Address.Base32()
	addr := fmt.Sprintf("%s%s:8000/game/client/index.html", scheme, string(b32))
	return c.Redirect(http.StatusFound, addr)
}
