package main

import (
	"testing"
)

func TestResolveARecord(t *testing.T) {
	resolver, err := NewQuickDNSResolver()
	if err != nil {
		t.Fatalf("Failed to create resolver: %v", err)
	}

	tests := []struct {
		domain string
		valid  bool
		ip     string
	}{
		{"ip-192-168-1-1.swiftwave.xyz.", true, "192.168.1.1"},
		{"ip-10-0-0-1.swiftwave.xyz.", true, "10.0.0.1"},
		{"ip-256-0-0-1.swiftwave.xyz.", false, ""},
		{"ip-192-168-1.swiftwave.xyz.", false, ""},
	}

	for _, test := range tests {
		valid, ip := resolver.ResolveARecord(test.domain)
		if valid != test.valid || ip != test.ip {
			t.Errorf("For domain %s, expected valid=%v and ip=%s, but got valid=%v and ip=%s",
				test.domain, test.valid, test.ip, valid, ip)
		}
	}
}
