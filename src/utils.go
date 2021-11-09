package main

import (
	"net"
)

var ipv4net *net.IPNet
var ipv6net *net.IPNet

// Helper function to check if an ip address is within our subnets
func isInSubnet(ip net.IP) bool {
	if ipv4net == nil {
		_, ipv4net, _ = net.ParseCIDR("128.153.144.0/23")
		_, ipv6net, _ = net.ParseCIDR("2605:6480:c051::1/48")
	}

	return ipv4net.Contains(ip) || ipv6net.Contains(ip)
}
