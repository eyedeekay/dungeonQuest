package main

import (
	"embed"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eyedeekay/onramp"
	"github.com/eyedeekay/unembed"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	gs "github.com/SineYuan/goBrowserQuest/bqs"
)

//go:embed conf/BrowserQuest/*
var defaultContent embed.FS

//go:embed conf/BrowserQuest/client/config/config_local.json
var defaultConfig []byte

var confFilePath = flag.String("config", "./config.json", "configuration file path")
var clientDir = flag.String("client", "./conf/BrowserQuest", "BrowserQuest root directory to serve if provided")
var clientReqPrefix = flag.String("prefix", "/game", "request url prefix when client is provided, cannot be '/' ")
var shortPort = flag.String("port", "7681", "port to present the plugin homepage on, actually a link to the game.")
var useI2P = flag.String("i2p", "127.0.0.1:7656", "The address of the SAMv3 API.")

var e *echo.Echo
var garlic *onramp.Garlic

func main() {
	flag.Parse()
	e = echo.New()
	e.Server.ReadTimeout = time.Hour
	e.Server.WriteTimeout = time.Hour
	e.Server.ReadHeaderTimeout = time.Hour
	e.Use(middleware.Recover())
	e.GET("/index.html", helloFunc)
	var err error
	garlic, err = onramp.NewGarlic("dungeonquest", *useI2P, onramp.OPT_WIDE)
	if err != nil {
		log.Fatal(err)
	}
	garlic.Timeout = time.Hour
	garlic.StreamSession.Timeout = time.Hour
	//garlic.StreamSession.SetReadDeadline(time.Now() + time.Hour)
	e.TLSListener, err = garlic.ListenTLS()
	if err != nil {
		log.Fatal(err)
	}

	if *clientDir != "" {
		if err := fixupDefaultDir(); err != nil {
			log.Fatal(err)
		}
		e.Static(*clientReqPrefix, *clientDir)
		if err := fixupConfigFiles(); err != nil {
			log.Fatal(err)
		}
	}
	config, err := gs.LoadConf(*confFilePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config)

	//go func() {
	bqs := gs.NewBQS(config)
	e.Any("/", bqs.ToEchoHandler())
	addrString := e.TLSListener.Addr().String()
	log.Println("Server is running at https://" + addrString)
	e.Logger.Fatal(http.Serve(e.TLSListener, e))
	//}()

	//server := http.ServeMux{}
	//server.HandleFunc("/index.html", hello)
	//http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", *shortPort), &server)
}

func fixupDefaultDir() error {
	configPath := filepath.Join(*clientDir, "bin")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err := unembed.Unembed(defaultContent, *clientDir, "conf/BrowserQuest")
		if err != nil {
			return err
		}
	}
	return nil
}

func fixupConfigFiles() error {
	configDir := filepath.Join(*clientDir, "/client/config/")
	configPath := filepath.Join(configDir, "config_local.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
		ioutil.WriteFile(configPath, defaultConfig, 0755)
	}
	log.Println("Adjusting config file,", configPath)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	log.Println("Adjusting config file", string(bytes))
	addr, err := garlic.Keys()
	if err != nil {
		return err
	}
	fixed := strings.Replace(string(bytes), "localhost", addr.Address.Base32(), -1)
	log.Println("Adjusted config file", fixed)
	err = ioutil.WriteFile(configPath, []byte(fixed), 0644)
	if err != nil {
		return err
	}
	configBuild := filepath.Join(*clientDir, "/client/config/config_build.json")
	err = ioutil.WriteFile(configBuild, []byte(fixed), 0644)
	if err != nil {
		return err
	}
	return nil
}
