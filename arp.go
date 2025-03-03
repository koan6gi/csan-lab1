package main

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Hosts struct {
	IP  string
	Mac string
}

func UpdateARPTable(ip net.IPNet) {
	var wg = sync.WaitGroup{}
	const threadCount = 10
	sem := make(chan byte, threadCount)

	firstIP, lastIP := getIPRange(ip)

	for i := firstIP; !isIPGrater(i, lastIP); IncIP(i) {
		sem <- 1
		wg.Add(1)
		go func() {
			Ping(i.String(), ip.IP.To4().String())
			<-sem
			wg.Done()
		}()
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
	close(sem)
}

func ParseARPTable(ip string) ([]Hosts, error) {
	cmd := exec.Command("arp", "-a")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	index := -1
	for i, line := range lines {
		if strings.Contains(line, ip) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil, fmt.Errorf("ARP table does not contain IP")
	}
	index += 2
	result := make([]Hosts, 0)

	for i := index; i < len(lines) && strings.TrimSpace(lines[i]) != ""; i++ {
		data := strings.Split(strings.TrimSpace(lines[i]), " ")
		iIp, iMac := -1, -1
		for j := 0; j < len(data); j++ {
			if data[j] != "" {
				if iIp == -1 {
					iIp = j
				} else {
					iMac = j
					break
				}
			}
		}
		if iIp == -1 || iMac == -1 {
			return nil, fmt.Errorf("unable to parse ARP table")
		}
		result = append(result, Hosts{string(data[iIp]), string(data[iMac])})
	}
	return result, nil
}
