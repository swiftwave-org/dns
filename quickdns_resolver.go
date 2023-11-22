package main

import (
	"fmt"
	"regexp"
	"strings"
)

/**
* Quick DNS Resolver
* This is a simple DNS resolver that resolves a domain name to an IP address without any records
* People need not to configure any DNS records for their domain names
* For example : ip-3-56-23-12.swiftwave.xyz, will resolve to 3.56.23.12
**/

type QuickDNSResolver struct {
	Regex *regexp.Regexp
}

func NewQuickDNSResolver() (*QuickDNSResolver, error) {
	pattern := `^(?:.*\.)?ip-(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)-(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)-(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)-(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.swiftwave\.xyz\.$`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %s", err.Error())
	}
	return &QuickDNSResolver{
		Regex: regex,
	}, nil
}

// ResolveARecord resolves a domain name to an IP address
// Returns true if the domain name is in the format ip-3-56-23-12.swiftwave.xyz
// Also returns the IP address
func (q *QuickDNSResolver) ResolveARecord(domain string) (bool, string) {
	domain = strings.ToLower(domain)
	if matches := q.Regex.FindStringSubmatch(domain); len(matches) == 5 {
		ip := fmt.Sprintf("%s.%s.%s.%s", matches[1], matches[2], matches[3], matches[4])
		return true, ip
	}
	return false, ""
}
