package main

import (
	"fmt"
	"net"
)

func GetIPNet() (*net.IPNet, error) {
	iface, err := chooseNetworkInterface()
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet, nil
		}
	}

	return nil, fmt.Errorf("cant find ipnet\n")
}

func chooseNetworkInterface() (*net.Interface, error) {
	ifaces, err := getActiveInterfaces()
	if err != nil {
		return nil, err
	}
	fmt.Println("Choose Network Interface")
	for i, iface := range ifaces {
		fmt.Printf("%d. %s\n", i, iface.Name)
	}
	k := 0
	_, _ = fmt.Scanf("%d", &k)
	if k >= len(ifaces) || k < 0 {
		return nil, fmt.Errorf("Invalid network interface\n")
	}

	return &ifaces[k], nil
}

func getActiveInterfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	activeInterfaces := make([]net.Interface, 0)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			activeInterfaces = append(activeInterfaces, iface)
		}
	}
	return activeInterfaces, nil
}
