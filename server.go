package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"log"
	"crypto/x509"
	"crypto/tls"
)

/* Resources:
	* http://www.levigross.com/2015/11/21/mutual-tls-authentication-in-go/
	* https://github.com/levigross/go-mutual-tls
	* https://ericchiang.github.io/post/go-tls/
 */

// greets the user with the email address of his certificate
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %v!\n", r.TLS.PeerCertificates[0].EmailAddresses[0])
}

func main() {
	certBytes, err := ioutil.ReadFile("./certs/cert.pem")

	if err != nil {
		log.Fatalln("Could not read cert.pem file", err)
	}

	clientCertPool := x509.NewCertPool()

	if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
		log.Fatalln("Could not add certificate to certificate pool")
	}

	tlsConfig := &tls.Config{
		// reject any TLS certificate that cannot be validated
		ClientAuth: tls.RequireAndVerifyClientCert,

		// ensure that only "our" CA is used to validate certificated
		ClientCAs: clientCertPool,

		// PFS becasue we can
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},

		// force it on server side
		PreferServerCipherSuites: true,

		// TLS 1.2 because we can
		MinVersion: tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	// serve index page
	http.HandleFunc("/", HelloHandler)

	httpServer := &http.Server{
		Addr: ":3300",
		TLSConfig: tlsConfig,
	}

	log.Println(httpServer.ListenAndServeTLS("./certs/cert.pem", "./certs/key.pem"))
}