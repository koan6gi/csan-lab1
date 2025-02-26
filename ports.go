package main

import (
	"net"
	"strconv"
	"time"
)

const (
	MAX_PORTS_COUNT = 2000
)

func GetAvailablePorts(hosts []Hosts) ([][]int, error) {
	ports := make([][]int, len(hosts))
	for k, host := range hosts {
		ports[k] = make([]int, 0)
		for i := 1; i < MAX_PORTS_COUNT; i++ {
			go func() {
				conn, err := net.DialTimeout("tcp", host.IP+":"+strconv.Itoa(i), time.Millisecond*2500)
				if err != nil {
					return
				}
				_ = conn.Close()
				ports[k] = append(ports[k], i)
			}()
		}
	}
	return ports, nil
}
