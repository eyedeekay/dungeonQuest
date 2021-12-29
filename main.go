package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"

	gs "github.com/SineYuan/goBrowserQuest/bqs"
)

var confFilePath = flag.String("config", "./config.json", "configuration file path")
var clientDir = flag.String("client", "./BrowserQuest", "BrowserQuest root directory to serve if provided")
var clientReqPrefix = flag.String("prefix", "/game", "request url prefix when client is provided, cannot be '/' ")
var useTLS = flag.Bool("tls", false, "use TLS")
var shortPort = flag.String("port", "7681", "port to present the plugin homepage on, actually a link to the game.")

var wide = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=4", "outbound.quantity=4"}

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
	addr := fmt.Sprintf("%s%s/game/client/index.html", scheme, string(b32))
	link := fmt.Sprintf("<a href=\"%s\">%s</a>", addr, addr)
	fmt.Fprintf(w, link)
}

func main() {
	flag.Parse()
	e := echo.New()
	e.Use(middleware.Recover())
	sam, err := sam3.NewSAM("127.0.0.1:7656")
	if err != nil {
		log.Fatal(err)
	}
	eepkeys, err := sam.EnsureKeyfile("dungeonquest.i2p.private")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("dungeonquest.i2p.public.txt", []byte(eepkeys.Addr().Base32()), 0644)
	if err != nil {
		log.Fatal(err)
	}
	session, err := sam.NewStreamSession("dungeonquest", eepkeys, wide)
	if err != nil {
		log.Fatal(err)
	}
	e.Listener, err = session.Listen()
	if err != nil {
		log.Fatal(err)
	}
	if *clientDir != "" {
		log.Println("Adjusting config file")
		bytes, err := ioutil.ReadFile(filepath.Join(*clientDir, "/client/config/config_local.json"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Adjusting config file", string(bytes))
		fixed := strings.Replace(string(bytes), "localhost", e.Listener.Addr().(i2pkeys.I2PAddr).Base32(), -1)
		log.Println("Adjusted config file", fixed)
		err = ioutil.WriteFile(filepath.Join(*clientDir, "/client/config/config_local.json"), []byte(fixed), 0644)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(filepath.Join(*clientDir, "/client/config/config_build.json"), []byte(fixed), 0644)
		if err != nil {
			log.Fatal(err)
		}
		e.Static(*clientReqPrefix, *clientDir)
	}
	config, err := gs.LoadConf(*confFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)
	bqs := gs.NewBQS(config)

	e.Any("/", bqs.ToEchoHandler())

	if *useTLS {
		e.TLSListener = tls.NewListener(e.Listener, &tls.Config{})
		defer e.Listener.Close()
	}

	server := http.ServeMux{}
	server.HandleFunc("/index.html", hello)
	go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", *shortPort), &server)

	if *useTLS {
		addr := fmt.Sprintf("%v", e.TLSListener.Addr().(i2pkeys.I2PAddr).Base32())
		log.Println("Server is running at https://" + addr)
		e.Logger.Fatal(http.Serve(e.TLSListener, e))
	} else {
		addr := fmt.Sprintf("%v", e.Listener.Addr().(i2pkeys.I2PAddr).Base32())
		log.Println("Server is running at http://" + addr)
		e.Logger.Fatal(http.Serve(e.Listener, e))
	}

}
