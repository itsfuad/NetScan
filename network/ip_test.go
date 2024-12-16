package network

import (
	"fmt"
	"testing"
)


func TestDetectSubnet(t *testing.T) {
	tests := []struct {
		ip       string
		expected string
	}{
		{"10.0.0.1", "10.0.0.0/8"},
		{"172.16.0.1", "172.16.0.0/16"},
		{"192.168.1.1", "192.168.1.0/24"},
		{"255.255.255.255", "255.255.255.0/24"},
		{"invalid_ip", ""},
	}

	for _, test := range tests {
		t.Run(test.ip, func(t *testing.T) {
			result := DetectSubnet(test.ip)
			fmt.Printf("DetectSubnet(%s) returned: %s\n", test.ip, result)
			if result != test.expected {
				t.Errorf("DetectSubnet(%s) = %s; expected %s", test.ip, result, test.expected)
			}
		})
	}
}
func TestDetectLocalIP(t *testing.T) {
	ip := DetectLocalIP()
	if ip == "" {
		t.Error("DetectLocalIP() returned an empty string; expected a non-loopback IPv4 address")
	} else {
		fmt.Printf("DetectLocalIP() returned: %s\n", ip)
	}
}



