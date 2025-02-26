package main

import (
	"net"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

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
				_ = exec.Command("ping", "-n", "2", "-w", "1000", dstIP+"."+strconv.Itoa(i)).Run()
				time.Sleep(time.Millisecond * 500)
				wg.Done()
			}()
		} else {
			pingIPs(dstIP+"."+strconv.Itoa(i), recDep)
		}
	}
}
