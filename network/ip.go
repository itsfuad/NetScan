package network

import (
	"fmt"
	"net"
	"sync"
	"time"
	"context"
)


// DetectSubnet calculates the subnet (in CIDR notation) from an IP address
func DetectSubnet(ip string) string {
	// Parse the provided IP address
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		fmt.Println("Invalid IP address:", ip)
		return ""
	}

	// Determine the subnet mask based on private IP ranges or default
	var subnetMask net.IPMask

	ipv4 := parsedIP.To4()

	if ipv4 != nil {
		firstOctet := ipv4[0]
		fmt.Printf("First octet: %d\n", firstOctet)
		secondOctet := ipv4[1] // For ranges like 172.16.0.0/16
		switch {
		case firstOctet == 10:
			// Private range 10.0.0.0/8
			subnetMask = net.CIDRMask(8, 32)
		case firstOctet == 172 && secondOctet >= 16 && secondOctet <= 31:
			// Private range 172.16.0.0/16
			subnetMask = net.CIDRMask(16, 32)
		case firstOctet == 192 && secondOctet == 168:
			// Private range 192.168.0.0/24
			subnetMask = net.CIDRMask(24, 32)
		default:
			// Default to /24 for other cases
			subnetMask = net.CIDRMask(24, 32)
		}
	} else {
		fmt.Println("Unsupported IP version for:", ip)
		return ""
	}

	// Apply the subnet mask to the IP address to determine the subnet
	ipNet := net.IPNet{
		IP:   parsedIP.Mask(subnetMask),
		Mask: subnetMask,
	}

	return ipNet.String()
}

// DetectLocalIP detects and returns the local machine's non-loopback IPv4 address
func DetectLocalIP() string {
	// Get a list of all network interfaces on the system
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network interfaces:", err)
		return ""
	}

	// Iterate through interfaces to find a non-loopback IPv4 address
	for _, iface := range interfaces {
		// Skip down interfaces and loopback addresses
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		ip := getNonLoopbackIPv4(iface)
		if ip != "" {
			fmt.Printf("Detected Local IP: %s\n", ip)
			return ip
		}
	}

	// If no non-loopback IPv4 address found, return an empty string
	return ""
}

// getNonLoopbackIPv4 returns the first non-loopback IPv4 address for the given interface
func getNonLoopbackIPv4(iface net.Interface) string {
	// Get all addresses associated with the interface
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			// Check if the address is an IPv4 address (ignore IPv6)
			if v.IP.To4() != nil {
				return v.IP.String() // Return the first non-loopback IPv4 address
			}
		}
	}

	return ""
}


// ScanPort checks if a specific port is open on a given IP using context
func ScanPort(ctx context.Context, ip string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", ip, port)

	// Use a context-aware dialer to respect cancellation and timeout
	conn, err := (&net.Dialer{
		Timeout: timeout,
	}).DialContext(ctx, "tcp", address)

	if err != nil {
		return false
	}

	defer conn.Close()
	return true
}

// ScanSubnet scans all IPs in the subnet for open ports in a specified range using context for timeout and cancellation
func ScanSubnet(subnet string, startPort, endPort int, timeout time.Duration) ([]string, error) {
	// Parse the subnet CIDR
	_, network, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet format: %v", err)
	}

	// Prepare context with a 5-minute overall timeout for scanning
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// List of open IPs
	var openIPs []string
	var wg sync.WaitGroup
	sem := make(chan struct{}, 100) // Limit concurrent goroutines

	// Iterate over the network range
	for ip := network.IP.Mask(network.Mask); network.Contains(ip); incrementIP(ip) {
		// Skip the broadcast address
		if ip[3] == 255 {
			continue
		}

		// Logging for debugging
		fmt.Printf("Scanning IP: %s\n", ip.String())
		wg.Add(1)

		go func(ip string) {
			defer wg.Done()

			// Acquire semaphore slot for concurrency control
			sem <- struct{}{}
			defer func() { <-sem }()

			// Check each port for open status
			for port := startPort; port <= endPort; port++ {
				// Respect context cancellation during the scan
				select {
				case <-ctx.Done():
					// If the context is done, exit early
					return
				default:
					// Continue scanning the port
				}

				if ScanPort(ctx, ip, port, timeout) {
					// If any port is open, add to list of open IPs
					openIPs = append(openIPs, ip)
					break // Stop once an open port is found
				}
			}
		}(ip.String())
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return openIPs, nil
}

// incrementIP increments the given IP address (byte-by-byte)
// It will stop once the broadcast address is reached (i.e., 255 in the last octet)
func incrementIP(ip net.IP) {
	// Increment the last byte
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		// Stop incrementing if the IP is beyond the network
		if ip[j] > 0 {
			break
		}
	}
}