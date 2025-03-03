package main

import (
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	MaxPortsCount = 2000
)

var waitGroupPorts sync.WaitGroup

func GetAvailablePorts(hosts []Hosts) ([][]int, error) {
	ports := make([][]int, len(hosts))

	for k, host := range hosts {

		ports[k] = make([]int, 0)

		for i := 1; i < MaxPortsCount; i++ {

			waitGroupPorts.Add(1)

			go func() {
				defer waitGroupPorts.Done()

				conn, err := net.DialTimeout("tcp", host.IP+":"+strconv.Itoa(i), time.Millisecond*2500)
				if err != nil {
					return
				}

				_ = conn.Close()
				ports[k] = append(ports[k], i)
			}()
		}
	}

	waitGroupPorts.Wait()

	return ports, nil
}
