package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

func forwardToExternalDNS(q dns.Question) ([]dns.RR, error) {
	// Query the external DNS server (e.g., 1.1.1.1)
	client := new(dns.Client)
	message := new(dns.Msg)
	message.SetQuestion(q.Name, q.Qtype)

	// Using 1.1.1.1 as an external resolver
	resp, _, err := client.Exchange(message, "1.1.1.1:53")
	if err != nil {
		return nil, err
	}
	return resp.Answer, nil
}

func parseQuery(m *dns.Msg, qdns *QuickDNSResolver) {
	for _, q := range m.Question {
		println("Query: ", q.Name, q.Qtype, q.Qclass)

		queryName := strings.ToLower(q.Name)

		if !strings.HasSuffix(queryName, ".swiftwave.xyz.") {
			// Query 1.1.1.1 or another DNS server and add the answer
			answers, err := forwardToExternalDNS(q)
			if err != nil {
				log.Printf("Failed to forward query: %v", err)
				return
			}
			m.Answer = append(m.Answer, answers...)
			return
		}

		switch q.Qtype {
		case dns.TypeNone:
			fallthrough
		case dns.TypeANY:
			fallthrough
		case dns.TypeNS:
			{
				// ns1.swiftwave.xyz and ns2.swiftwave.xyz
				rr, err := dns.NewRR(fmt.Sprintf("%s NS ns1.swiftwave.xyz", q.Name))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
				rr, err = dns.NewRR(fmt.Sprintf("%s NS ns2.swiftwave.xyz", q.Name))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		case dns.TypeA:
			{
				isQDNS, ip := qdns.ResolveARecord(queryName)
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
	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		handleDnsRequest(w, m, qdns)
	})
	// start server
	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "0.0.0.0"
		log.Printf("Defaulting to address %s", address)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "53"
		log.Printf("Defaulting to port %s", port)
	}
	server := &dns.Server{Addr: ":" + port, Net: "udp"}
	log.Printf("Starting at %s\n", port)
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
