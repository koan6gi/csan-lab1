package main

import (
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func isIPGrater(ip1, ip2 net.IP) bool {
	return (ip1[12] > ip2[12]) ||
		(ip1[12] == ip2[12] && ip1[13] > ip2[13]) ||
		(ip1[12] == ip2[12] && ip1[13] == ip2[13] && ip1[14] > ip2[14]) ||
		(ip1[12] == ip2[12] && ip1[13] == ip2[13] && ip1[14] == ip2[14] && ip1[15] > ip2[15])
}

func IncIP(ip net.IP) {
	if ip == nil {
		return
	}

	shift := 1
	for i := 15; i > 11; i-- {
		newByte := int(ip[i])
		newByte += shift
		shift = 0
		if newByte > 255 {
			newByte %= 256
			shift = 1
		}
		ip[i] = byte(newByte)
	}

}

func getIPRange(ip net.IPNet) (net.IP, net.IP) {
	var firstIP, lastIP = make(net.IP, len(ip.IP)), make(net.IP, len(ip.IP))
	copy(firstIP, ip.IP)

	for i := 0; i < 4; i++ {
		firstIP[i+12] &= ip.Mask[i]
	}

	reverseMask := make([]byte, len(ip.Mask))
	copy(reverseMask, ip.Mask)

	for i := 0; i < len(reverseMask); i++ {
		reverseMask[i] ^= 0xFF
	}

	copy(lastIP, firstIP)

	for i := 0; i < 4; i++ {
		lastIP[i+12] |= reverseMask[i]
	}

	return firstIP, lastIP
}

func Ping(dstIP, srcIP string) {
	conn, err := icmp.ListenPacket("ip4:icmp", srcIP)
	if err != nil {
		return
		//log.Fatalf("listen err (%s), %s", dstIP, err)
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("ICMP Echo"),
		},
	}

	byteMsg, err := msg.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := conn.WriteTo(byteMsg, &net.IPAddr{IP: net.ParseIP(dstIP)}); err != nil {
		log.Fatalf("WriteTo err (%s), %s", dstIP, err)
	}

	//conn, err := net.DialTimeout("tsp", dstIP+":80", 1*time.Second)
	//if err != nil {
	//	return
	//}
	//defer conn.Close()
}
