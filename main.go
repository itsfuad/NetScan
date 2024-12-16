package main

import (
	"flag"
	"fmt"
	"netscan/network"
	"os"
	"time"
)

func handleFlags() (string, int, int, time.Duration) {
	var ip string
	var startPort int
	var endPort int
	var timeout time.Duration

	flag.StringVar(&ip, "ip", "", "IP address to scan")
	flag.IntVar(&startPort, "start-port", 1, "Start port for scanning")
	flag.IntVar(&endPort, "end-port", 1024, "End port for scanning")
	flag.DurationVar(&timeout, "timeout", 2*time.Second, "Timeout for port scanning")
	flag.Usage = func() {
		fmt.Println("Usage: netscan [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if ip == "" {
		ip = network.DetectLocalIP()
	}

	if startPort < 1 || startPort > 65535 {
		fmt.Printf("Invalid port. Start port %d is not in the range 1-65535\n", startPort)
		os.Exit(-1)
	}

	if startPort > endPort {
		fmt.Printf("Invalid port. Start port %d is greater than end port %d\n", startPort, endPort)
		os.Exit(-1)
	}

	subnet := network.DetectSubnet(ip)

	if subnet == "" {
		fmt.Println("Failed to detect subnet")
		os.Exit(-1)
	}

	return subnet, startPort, endPort, timeout
}

func main() {

	subnet, startPort, endPort, timeout := handleFlags()

	fmt.Println("Scanning the subnet for open ports...")

	openIPs, err := network.ScanSubnet(subnet, startPort, endPort, timeout)
	if err != nil {
		fmt.Println("Error scanning subnet:", err)
		return
	}

	if len(openIPs) == 0 {
		fmt.Println("No open ports found in the subnet.")
	} else {
		fmt.Println("Open IP addresses found:")
		for _, openIP := range openIPs {
			fmt.Println("- http://" + openIP)
		}
	}
}
