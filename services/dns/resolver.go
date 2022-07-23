package dns

import "github.com/miekg/dns"

type DNSResolver struct {
	Records   map[string]string
	Ip        string
	ServerTCP *dns.Server
	ServerUDP *dns.Server
}

func NewDNSResolver(ip string) *DNSResolver {
	serverTCP := &dns.Server{
		Addr: ip + ":53",
		Net:  "tcp",
	}
	serverUDP := &dns.Server{
		Addr: ip + ":53",
		Net:  "udp",
	}

	return &DNSResolver{
		Records:   make(map[string]string),
		ServerTCP: serverTCP,
		ServerUDP: serverUDP,
		Ip:        ip,
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
	return dns.Ip
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
