package utils

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// GetMacAddr get local pc Mac Address list
func GetMacAddr() ([]string, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

// GetPCName get pc Name
func GetPCName() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	fmt.Println("hostname:", name)
	return name, nil
}

// GetMapMacIP get local pc mac map to ip
func GetMapMacIP() (map[string][]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	macips := make(map[string][]string)

	for _, iface := range ifaces {
		a := iface.HardwareAddr.String()
		if a == "" {
			continue
		}
		var ips []string
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			//return ip.String(), nil
			ips = append(ips, ip.String())
		}
		macips[a] = ips
	}
	return macips, nil
}

// GetIPAddrs get local ip list
func GetIPAddrs() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			//return ip.String(), nil
			ips = append(ips, ip.String())
		}
	}
	if len(ips) > 0 {
		return ips, nil
	}
	return nil, errors.New("are you connected to the network?")

}
