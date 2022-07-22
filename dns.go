package main

import (
	"fmt"

	"github.com/miekg/dns"
)

var (
	dnsResolver *DNSResolver
)

func init() {
	dns.HandleFunc(".", dnsHandleRequest)
}

func dnsHandleRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		dnsParseQuery(msg)
	}

	w.WriteMsg(msg)
}

func dnsParseQuery(msg *dns.Msg) {
	for _, q := range msg.Question {
		switch q.Qtype {
		case dns.TypeA:
			ip := dnsResolver.GetIP(q.Name)
			if ip != "" {
				fmt.Printf("[DNS] %s: %s\n", q.Name, ip)
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					msg.Answer = append(msg.Answer, rr)
				}
			}
		}
	}
}

type DNSResolver struct {
	Records map[string]string
	ServerTCP *dns.Server
	ServerUDP *dns.Server
}

func NewDNSResolver() *DNSResolver {
	serverTCP := &dns.Server{
		Addr: ip + ":53",
		Net: "tcp",
	}
	serverUDP := &dns.Server{
		Addr: ip + ":53",
		Net: "udp",
	}

	return &DNSResolver{
		Records: make(map[string]string),
		ServerTCP: serverTCP,
		ServerUDP: serverUDP,
	}
}

func (dns *DNSResolver) ListenAndServeTCP() error {
	defer dns.ServerTCP.Shutdown()
	return dns.ServerTCP.ListenAndServe()
}
func (dns *DNSResolver) ListenAndServeUDP() error {
	defer dns.ServerUDP.Shutdown()
	return dns.ServerUDP.ListenAndServe()
}

func (dns *DNSResolver) Close() error {
	if err := dns.ServerTCP.Shutdown(); err != nil {
		return err
	}
	if err := dns.ServerUDP.Shutdown(); err != nil {
		return err
	}
	return nil
}

func (dns *DNSResolver) GetIP(host string) string {
	if _, ok := dns.Records[host]; ok {
		return dns.Records[host]
	}
	return "0.0.0.0"
}

func (dns *DNSResolver) Add(host, ip string) *DNSResolver {
	dns.Records[host] = ip
	return dns
}

func (dns *DNSResolver) RemoveHost(host string) *DNSResolver {
	delete(dns.Records, host)
	return dns
}

func (dns *DNSResolver) RemoveIP(ip string) *DNSResolver {
	newRecords := make(map[string]string)
	for oldHost, oldIP := range dns.Records {
		if oldIP != ip {
			newRecords[oldHost] = oldIP
		}
	}
	dns.Records = newRecords
	return dns
}