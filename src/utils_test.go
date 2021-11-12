package main

import (
	"net"
	"testing"
)

func TestIsInSubnet(t *testing.T) {
	var tests = []struct {
		ipstr string
		want  bool
	}{
		{"128.153.145.20", true},
		{"128.153.144.250", true},
		{"128.153.144.140", true},
		{"128.153.1.1", false},
		{"128.153.146.1", false},
		{"127.0.0.1", false},
		{"::1", false},
		{"2605:6480:c051:100::1", true},
		{"2605:6480::1", false},
	}

	for _, tt := range tests {
		t.Run(tt.ipstr, func(t *testing.T) {
			ip := net.ParseIP(tt.ipstr)
			ans := isInSubnet(ip)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}
