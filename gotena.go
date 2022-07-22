package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

/* NOTES
- HTTP error codes start with 340, ex 340404 for 404
- On links, when linking to .ugo pages, use .uls anyway, the client will still request .ugo
*/

func main() {
	dnsResolver = NewDNSResolver().
		Add("conntest.nintendowifi.net.", "72.9.147.58").
		Add("dsnattest.available.gs.nintendowifi.net.", "72.9.147.58").
		Add("nas.nintendowifi.net.", "178.62.43.212").
		Add("nas.wiimmfi.de.", "178.62.43.212").
		Add("ugomemo.hatena.ne.jp.", "72.9.147.58").
		Add("flipnote.hatena.com.", "72.9.147.58")

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
