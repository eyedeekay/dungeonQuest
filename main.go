package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

var wide = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=4", "outbound.quantity=4"}

func main() {
	flag.Parse()
	config, err := gs.LoadConf(*confFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)
	bqs := gs.NewBQS(config)

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

	if *useTLS {
		e.TLSListener = tls.NewListener(e.Listener, &tls.Config{})
		defer e.Listener.Close()
	}

	if *clientDir != "" {
		bytes, err := ioutil.ReadFile("BrowserQuest/client/config/config_local.json")
		if err != nil {
			log.Fatal(err)
		}
		fixed := strings.Replace(string(bytes), "Set local dev websocket host here", e.Listener.Addr().(i2pkeys.I2PAddr).Base32(), -1)
		err = ioutil.WriteFile("BrowserQuest/client/config/config_local.json", []byte(fixed), 0644)
		if err != nil {
			log.Fatal(err)
		}
		e.Static(*clientReqPrefix, *clientDir)
	}
	e.Any("/", bqs.ToEchoHandler())

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
