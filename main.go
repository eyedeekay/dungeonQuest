package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"

	gs "github.com/SineYuan/goBrowserQuest/bqs"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func certgen(address string) *tls.Config {
	// priv, err := rsa.GenerateKey(rand.Reader, *rsaBits)
	var priv crypto.Signer
	prepriv, err := ioutil.ReadFile(address + ".pem")
	if err != nil {
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(address+".pem", pem.EncodeToMemory(pemBlockForKey(priv)), 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		block, _ := pem.Decode(prepriv)
		switch block.Type {
		case "RSA PRIVATE KEY":
			priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				log.Fatal(err)
			}
		case "EC PRIVATE KEY":
			priv, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatalf("Unrecognized key type %q", block.Type)
		}
	}
	derBytes, err := ioutil.ReadFile(address + ".crt")
	if err != nil {
		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization: []string{address},
			},
			Issuer: pkix.Name{
				Organization: []string{address},
			},
			NotBefore: time.Now(),
			NotAfter:  time.Now().Add(time.Hour * 24 * 180),

			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}

		derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
		if err != nil {
			log.Fatalf("Failed to create certificate: %s", err)
		}
		err = ioutil.WriteFile(address+".crt", derBytes, 0644)
		if err != nil {
			log.Fatalf("Failed to write data to file: %s", err)
		}
	}
	config := &tls.Config{
		ServerName: address,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{derBytes},
				PrivateKey:  priv,
			},
		},
	}
	return config
	/*out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	fmt.Println(out.String())
	out.Reset()
	pem.Encode(out, pemBlockForKey(priv))
	fmt.Println(out.String())*/
}

var confFilePath = flag.String("config", "./config.json", "configuration file path")
var clientDir = flag.String("client", "./BrowserQuest", "BrowserQuest root directory to serve if provided")
var clientReqPrefix = flag.String("prefix", "/game", "request url prefix when client is provided, cannot be '/' ")
var useTLS = flag.Bool("tls", true, "use TLS")
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
	addr := fmt.Sprintf("%s%s:8000/game/client/index.html", scheme, string(b32))
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
		e.TLSListener = tls.NewListener(e.Listener, certgen(e.Listener.Addr().(i2pkeys.I2PAddr).Base32()))
		defer e.Listener.Close()
	}

	server := http.ServeMux{}
	server.HandleFunc("/index.html", hello)
	go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", *shortPort), &server)

	if *useTLS {
		addr := fmt.Sprintf("%v", e.TLSListener.Addr().(i2pkeys.I2PAddr).Base32())
		log.Println("Server is running at https://" + addr + ":8000")
		e.Logger.Fatal(http.Serve(e.TLSListener, e))
	} else {
		addr := fmt.Sprintf("%v", e.Listener.Addr().(i2pkeys.I2PAddr).Base32())
		log.Println("Server is running at http://" + addr + ":8000")
		e.Logger.Fatal(http.Serve(e.Listener, e))
	}

}
