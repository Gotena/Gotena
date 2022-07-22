package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	ip string
)

/* NOTES
- HTTP error codes start with 340, ex 340404 for 404
- On links, when linking to .ugo pages, use .uls anyway, the client will still request .ugo
*/

func main() {
	if len(os.Args) == 1 {
		panic("please specify the IP to host on!")
	}
	ip = os.Args[1]

	dnsResolver = NewDNSResolver().
		Add("conntest.nintendowifi.net.", ip).
		Add("dsnattest.available.gs.nintendowifi.net.", ip).
		Add("nas.nintendowifi.net.", "178.62.43.212").
		Add("nas.wiimmfi.de.", "178.62.43.212").
		Add("ugomemo.hatena.ne.jp.", ip).
		Add("flipnote.hatena.com.", ip)

	ht = NewHTTP().
		GET("conntest.nintendowifi.net", "/", httpConntest).
		GET("dsnattest.available.gs.nintendowifi.net", "/", httpConntest).
		GET("flipnote.hatena.com", "/:index", httpFlipnoteRootGET).
		GET("flipnote.hatena.com", "/:index/:region", httpFlipnoteRootGET).
		GET("flipnote.hatena.com", "/:index/:region/:language", httpFlipnoteRootGET).
		GET("flipnote.hatena.com", "/:index/:region/:language/:page", httpFlipnoteRootGET).
		POST("flipnote.hatena.com", "/:index", httpFlipnoteRootPOST).
		POST("flipnote.hatena.com", "/:index/:region", httpFlipnoteRootPOST).
		POST("flipnote.hatena.com", "/:index/:region/:language", httpFlipnoteRootPOST).
		POST("flipnote.hatena.com", "/:index/:region/:language/:page", httpFlipnoteRootPOST)

	go func() {
		dnsErr := dnsResolver.ListenAndServeTCP()
		if dnsErr != nil {
			panic(dnsErr)
		}
	}()
	go func() {
		dnsErr := dnsResolver.ListenAndServeUDP()
		if dnsErr != nil {
			panic(dnsErr)
		}
	}()
	go func() {
		httpErr := ht.ListenAndServe()
		if httpErr != nil {
			panic(httpErr)
		}
	}()

	fmt.Println("Idling...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT)
	<-sc
	fmt.Println(" ")

	fmt.Println("Shutting down...")
	dnsResolver.Close()
	ht.Close()

	fmt.Println("Good-bye!")
}
