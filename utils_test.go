package main

import (
	"net"
	"testing"
)

func TestContains(t *testing.T) {
	n := NewNetworks([]string{"::1/128"})
	if !n.Contains(net.ParseIP("::1")) {
		t.Errorf("Expected ::1 to be in ::1/128")
	}
}
