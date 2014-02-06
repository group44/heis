package main

import (
    //"fmt"
    "net"
    "time"
)

func main() {
	laddr, _ := net.ResolveUDPAddr("udp4", ":0")
	gaddr, _ := net.ResolveUDPAddr("udp4", "224.0.0.1:12000") // .255
	
	lconn, _ := net.ListenUDP("udp", laddr)
	data := make([]byte, 1)
	
	for {
		_, _ = lconn.WriteToUDP(data, gaddr)
		time.Sleep(500 * time.Millisecond)
	}
}
