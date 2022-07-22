package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	ht *HTTP
)

// HostSwitch maps host names to http.Handler structs
type HostSwitch map[string]http.Handler

// Implement the ServeHTTP method on HostSwitch
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.
	if handler := hs[r.Host]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		// Handle host names for which no handler is registered
		// Pretend access is denied always
		fmt.Printf("[UNKNOWN] %s %s", r.Host, r.URL)
		http.Error(w, "Forbidden", 403) // Or Redirect?
	}
}

type HTTP struct {
	HostSwitch HostSwitch
	Server    http.Server
	ServerSSL http.Server
}

func NewHTTP() *HTTP {
	hs := make(HostSwitch)

	server := http.Server{
		Addr: ":80",
	}

	/*rootCAs := x509.NewCertPool()
	certCAs, err := ioutil.ReadFile("certs/ca-cert.pem")
	crashError(err)

	if ok := rootCAs.AppendCertsFromPEM(certCAs); !ok {
		panic("invalid cert CA")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionSSL30,
		MaxVersion:         tls.VersionSSL30,
		RootCAs:            rootCAs,
		CipherSuites: []uint16{
			0x0004,
			0x0005,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		},
	}*/

	serverSSL := http.Server{
		Addr: ":443",
		//TLSConfig:    tlsConfig,
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &HTTP{
		HostSwitch: hs,
		Server: server,
		ServerSSL: serverSSL,
	}
}

func (ht *HTTP) ListenAndServe() error {
	ht.Server.Handler = ht.HostSwitch
	return ht.Server.ListenAndServe()
}
func (ht *HTTP) ListenAndServeSSL() error {
	ht.ServerSSL.Handler = ht.HostSwitch
	//return ht.Server.ListenAndServeTLS(certFile, keyFile)
	return nil
}

func (ht *HTTP) GET(host, path string, handler httprouter.Handle) *HTTP {
	router := httprouter.New()
	if _, ok := ht.HostSwitch[host]; ok {
		router = ht.HostSwitch[host].(*httprouter.Router)
	}
	router.GET(path, handler)
	ht.HostSwitch[host] = router
	return ht
}
func (ht *HTTP) POST(host, path string, handler httprouter.Handle) *HTTP {
	router := httprouter.New()
	if _, ok := ht.HostSwitch[host]; ok {
		router = ht.HostSwitch[host].(*httprouter.Router)
	}
	router.POST(path, handler)
	ht.HostSwitch[host] = router
	return ht
}
func (ht *HTTP) RemoveHost(host string) *HTTP {
	delete(ht.HostSwitch, host)
	return ht
}

func (ht *HTTP) Close() error {
	if err := ht.Server.Close(); err != nil {
		return err
	}
	if err := ht.ServerSSL.Close(); err != nil {
		return err
	}
	return nil
}
