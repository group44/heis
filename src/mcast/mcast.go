package main

import (
    "fmt"
    "net"
    "time"
)

func main() {
	BROADCAST := "129.241.187.255:12000"
	baddr, _ := net.ResolveUDPAddr("udp4", BROADCAST)
	conn, _ := net.DialUDP("udp4", nil, baddr)
	data := make([]byte, 1)
	
	for {
		conn.Write(data)
		time.Sleep(time.Second)
		fmt.Println("Message sent")
	}
}
