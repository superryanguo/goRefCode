package main

import (
	"fmt"
	"net"
)

func main() {
	domain := "www.bing.com"
	ips, err := net.LookupIP(domain)
	if err != nil {
		fmt.Printf("Error looking up IP for %s: %v\n", domain, err)
		return
	}

	fmt.Printf("IP addresses for %s:\n", domain)
	for _, ip := range ips {
		fmt.Println(ip)
	}
}
