package main

import (
    //"fmt"
    "net"
    "time"
)

func main() {
	laddr, _ := net.ResolveUDPAddr("udp4", ":0")
	gaddr, _ := net.ResolveUDPAddr("udp4", "224.0.0.2:12000")
	
	lconn, _ := net.ListenUDP("udp", laddr)
	data := make([]byte, 256)
	
	for {
		_, _ = lconn.WriteToUDP(data, gaddr)
		time.Sleep(500 * time.Millisecond)
	}
}
