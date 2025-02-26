package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Hosts struct {
	IP  string
	Mac string
}

var wg = sync.WaitGroup{}

func UpdateARPTable(ip net.IPNet) {
	pos := 4
	for i := 0; i < len(ip.Mask); i++ {
		if ip.Mask[i] == 0 {
			pos = i
			break
		}
	}
	ifaceIP := ip.IP.To4().String()
	tmp := ifaceIP
	k := pos
	pos = 4 - pos
	for i := 0; i < len(tmp); i++ {
		if tmp[i] == '.' {
			k--
			if k == 0 {
				k = i
				break
			}
		}
	}
	tmp = tmp[:k]
	pingIPs(tmp, pos)
	wg.Wait()
}

func pingIPs(dstIP string, recDep int) {
	recDep--
	if recDep < 0 {
		return
	}
	for i := 0; i < 256; i++ {
		if recDep == 0 {
			wg.Add(1)
			go func() {
				_ = exec.Command("ping", "-n", "1", "-4", dstIP+"."+strconv.Itoa(i)).Run()
				time.Sleep(time.Millisecond * 500)
				wg.Done()
			}()
		} else {
			pingIPs(dstIP+"."+strconv.Itoa(i), recDep)
		}
	}
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
