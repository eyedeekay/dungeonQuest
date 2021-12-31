package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	scheme := "http://"
	if *useTLS {
		scheme = "https://"
	}
	if *useI2P {
		b32, err := ioutil.ReadFile("dungeonquest.i2p.public.txt")
		if err != nil {
			log.Println(err)
			return
		}
		addr := fmt.Sprintf("%s%s:8000/game/client/index.html", scheme, string(b32))
		http.Redirect(w, req, addr, http.StatusFound)
	} else {
		addr := fmt.Sprintf("%s%s:8000/game/client/index.html", scheme, e.Listener.Addr().String())
		http.Redirect(w, req, addr, http.StatusFound)
	}

}
