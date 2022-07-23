package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
		Add("nas.nintendowifi.net.", "178.62.43.212").
		Add("nas.wiimmfi.de.", "178.62.43.212").
		Add("ugomemo.hatena.ne.jp.", ip).
		Add("flipnote.hatena.com.", ip)

	ht = NewHTTP().
		GET("conntest.nintendowifi.net", "/", httpConntest).
		GET("flipnote.hatena.com", "/:namespace", httpFlipnoteRootGET).
		GET("flipnote.hatena.com", "/:namespace/:region", httpFlipnoteRootGET).
		GET("flipnote.hatena.com", "/:namespace/:region/:index", httpFlipnoteRootGET).
		GET("flipnote.hatena.com", "/:namespace/:region/:index/:i1", httpFlipnoteRootGET).
		GET("flipnote.hatena.com", "/:namespace/:region/:index/:i1/:i2", httpFlipnoteRootGET).
		POST("flipnote.hatena.com", "/:namespace", httpFlipnoteRootPOST).
		POST("flipnote.hatena.com", "/:namespace/:region", httpFlipnoteRootPOST).
		POST("flipnote.hatena.com", "/:namespace/:region/:index", httpFlipnoteRootPOST).
		POST("flipnote.hatena.com", "/:namespace/:region/:index/:i1", httpFlipnoteRootPOST).
		POST("flipnote.hatena.com", "/:namespace/:region/:index/:i1/:i2", httpFlipnoteRootPOST)

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

func ReadUTF16String(data []byte) string {
	win16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16le.NewDecoder())
	unicodeReader := transform.NewReader(bytes.NewReader(data), utf16bom)
	decoded, err := ioutil.ReadAll(unicodeReader)
	if err != nil {
		panic(err)
	}

	return string(decoded)
}
func WriteUTF16String(data string) []byte {
	win16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom := unicode.BOMOverride(win16le.NewEncoder())
	unicodeWriter := transform.NewReader(bytes.NewReader([]byte(data)), utf16bom)
	encoded, err := ioutil.ReadAll(unicodeWriter)
	if err != nil {
		panic(err)
	}

	return encoded
}