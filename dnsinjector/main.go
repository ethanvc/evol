package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"test.service.": "192.168.0.2",
}

func printRequest(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, req *dns.Msg) {
	printRequest(req)
	resp, err := dns.ExchangeContext(context.Background(), req, "8.8.8.8:53")
	if err != nil {
		fmt.Printf("err:%s\n", err)
		return
	}
	err = w.WriteMsg(resp)
	if err != nil {
		fmt.Printf("err:%s\n", err)
		return
	}
}

// dig @localhost -p 5999 www.baidu.com

func main() {
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := 5999
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
