package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

func parseQuery(m *dns.Msg, qdns *QuickDNSResolver) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			{
				isQDNS, ip := qdns.ResolveARecord(q.Name)
				if isQDNS {
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
					if err == nil {
						m.Answer = append(m.Answer, rr)
					}
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg, qdns *QuickDNSResolver) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, qdns)
	}
	_ = w.WriteMsg(m)
}

func main() {
	// create quick dns resolver
	qdns, err := NewQuickDNSResolver()
	if err != nil {
		panic(err)
	}
	// attach request handler func
	dns.HandleFunc("swiftwave.xyz", func(w dns.ResponseWriter, m *dns.Msg) {
		handleDnsRequest(w, m, qdns)
	})
	// start server
	port := 5555
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
	// defer shutdown
	defer func() {
		err := server.Shutdown()
		if err != nil {
			log.Fatalf("Failed to shutdown server: %s\n ", err.Error())
		}
	}()
}
