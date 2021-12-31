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
	b32, err := ioutil.ReadFile("dungeonquest.i2p.public.txt")
	if err != nil {
		log.Println(err)
		return
	}
	addr := fmt.Sprintf("%s%s:8000/game/client/index.html", scheme, string(b32))
	//link := fmt.Sprintf("<a href=\"%s\">%s</a>", addr, addr)
	http.Redirect(w, req, addr, http.StatusFound)
	//fmt.Fprintf(w, link)
}
