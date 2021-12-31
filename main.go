package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
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
var useTLS = flag.Bool("tls", true, "use TLS")
var shortPort = flag.String("port", "7681", "port to present the plugin homepage on, actually a link to the game.")
var useI2P = flag.Bool("i2p", false, "use I2P")

var wide = []string{"inbound.length=1", "outbound.length=1",
	"inbound.lengthVariance=0", "outbound.lengthVariance=0",
	"inbound.backupQuantity=1", "outbound.backupQuantity=1",
	"inbound.quantity=4", "outbound.quantity=4"}

var e *echo.Echo

func main() {
	flag.Parse()
	e = echo.New()
	e.Use(middleware.Recover())
	if *useI2P {
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
	} else {
		var err error
		e.Listener, err = net.Listen("tcp", ":7681")
		if err != nil {
			log.Fatal(err)
		}
	}
	if *clientDir != "" && *useI2P {
		log.Println("Adjusting config file,", *clientDir, "/client/config/config_local.json")
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
	addrString := e.Listener.Addr().String()
	if *useTLS {
		if *useI2P {
			addrString = e.Listener.Addr().(i2pkeys.I2PAddr).Base32()
			e.TLSListener = tls.NewListener(e.Listener, certgen(addrString))
			defer e.Listener.Close()
		} else {
			addrString = e.Listener.Addr().String()
			e.TLSListener = tls.NewListener(e.Listener, certgen(addrString))
			defer e.Listener.Close()
		}
	}

	if *useI2P {
		server := http.ServeMux{}
		server.HandleFunc("/index.html", hello)
		go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", *shortPort), &server)
	}

	if *useTLS {
		addr := fmt.Sprintf("%v", addrString)
		log.Println("Server is running at https://" + addr)
		e.Logger.Fatal(http.Serve(e.TLSListener, e))
	} else {
		addr := fmt.Sprintf("%v", addrString)
		log.Println("Server is running at http://" + addr)
		e.Logger.Fatal(http.Serve(e.Listener, e))
	}

}
