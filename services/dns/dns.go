package dns

import (
	"fmt"

	"github.com/miekg/dns"
)

func (d *DNSResolver) Init() {
	dns.HandleFunc(".", d.dnsHandleRequest)
}

func (d *DNSResolver) dnsHandleRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		d.dnsParseQuery(msg)
	}

	w.WriteMsg(msg)
}

func (d *DNSResolver) dnsParseQuery(msg *dns.Msg) {
	for _, q := range msg.Question {
		switch q.Qtype {
		case dns.TypeA:
			ip := d.GetIP(q.Name)
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
