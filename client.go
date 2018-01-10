package main

import (
	"crypto/tls"
	"log"
	"io/ioutil"
	"crypto/x509"

	"github.com/levigross/grequests"
	"net/http"
)

func main() {
	cert, err := tls.LoadX509KeyPair("../certs/cert.pem", "../certs/key.pem")

	if err != nil {
		log.Fatalln("Could not load cert", err)
	}

	clientCACert, err := ioutil.ReadFile("../certs/cert.pem")

	if err != nil {
		log.Fatalln("Unable to open cert", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs: clientCertPool,
	}

	tlsConfig.BuildNameToCertificate()

	ro := &grequests.RequestOptions{
		HTTPClient: &http.Client{
			Transport: &http.Transport{TLSClientConfig: tlsConfig},
		},
	}

	resp, err := grequests.Get("https://localhost:3300", ro)

	if err != nil {
		log.Println("Could not communicate with server", err)
	}

	log.Println(resp.String())
}
