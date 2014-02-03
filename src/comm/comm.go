package comm

import (
	"net"
	"fmt"
	"code.google.com/p/go.net/ipv4"
)


func Establish() {
	
	// IP til egen maskin her
	testIP := "129.241.187.146:12000"
	udpAddr, err := net.ResolveUDPAddr("udp4", testIP)

	en0, err := net.InterfaceByName("en0")
	//checkError(err)
	group := net.IPv4(224, 0, 0, 250)
	c, err := net.ListenPacket("udp4", "0.0.0.0:12000")
	//checkError(err)
	p := ipv4.NewPacketConn(c)
	p.JoinGroup(en0, udpAddr)
	//checkError(err)
	
	fmt.Println("Test...")

	fmt.Println(err)
	socket.Close()
}
