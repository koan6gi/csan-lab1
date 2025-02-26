package main

import (
	"fmt"
)

func main() {
	ipnet, err := GetIPNet()
	if err != nil {
		_ = fmt.Errorf("%e\n", err)
		return
	}
	UpdateARPTable(*ipnet)
	fmt.Println("Check arp")
}
