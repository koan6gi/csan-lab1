package main

import (
	"log"
)

func main() {
	iface, err := ChooseNetworkInterface()
	if err != nil {
		log.Println(err)
		return
	}
	ipnet, err := GetIPNet(iface)
	if err != nil {
		log.Println(err)
		return
	}
	UpdateARPTable(*ipnet)
	hosts, err := ParseARPTable(ipnet.IP.To4().String())
	if err != nil {
		log.Println(err)
		return
	}
	ShowMyInterfaceInfo(ipnet.IP.To4().String(), iface.HardwareAddr.String())
	ports, _ := GetAvailablePorts(hosts)
	ShowHostsInfo(hosts, ports)
}
