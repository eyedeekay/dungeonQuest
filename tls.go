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
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"time"
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
}
