package main

import (
	"fmt"
)

func main() {
	iface, err := ChooseNetworkInterface()
	if err != nil {
		_ = fmt.Errorf("%e\n", err)
		return
	}
	ipnet, err := GetIPNet(iface)
	if err != nil {
		_ = fmt.Errorf("%e\n", err)
		return
	}
	UpdateARPTable(*ipnet)
	hosts, err := ParseARPTable(ipnet.IP.To4().String())
	if err != nil {
		_ = fmt.Errorf("%e\n", err)
		return
	}
}
