package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/JoshuaDoes/gotena/services/dns"
	"github.com/JoshuaDoes/gotena/services/web"
	"github.com/gofiber/fiber/v2"
)

var (
	ip          string
	ht          *fiber.App
	dnsResolver *dns.DNSResolver
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

	dnsResolver = dns.NewDNSResolver(ip).
		Add("conntest.nintendowifi.net.", ip).
		Add("nas.nintendowifi.net.", "178.62.43.212").
		Add("nas.wiimmfi.de.", "178.62.43.212").
		Add("ugomemo.hatena.ne.jp.", ip).
		Add("flipnote.hatena.com.", ip)

	ht = web.InitFiber()

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
		httpErr := ht.Listen(":80")
		if httpErr != nil {
			panic(httpErr)
		}
	}()
	dnsResolver.Init()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT)
	<-sc
	fmt.Println(" ")

	fmt.Println("Shutting down...")
	dnsResolver.Close()
	ht.Shutdown()

	fmt.Println("Good-bye!")
}
