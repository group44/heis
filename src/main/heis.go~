package main

import (
	"fmt"
	"comm"
	"net"
	//"os"
)

func main() {
	fmt.Println("CreateSocket begin")

	localaddr, err := net.ResolveUDPAddr("udp4", ":0")
	groupaddr, err := net.ResolveUDPAddr("udp4", "224.0.0.2:12000")
	comm.CheckError(err)

	localconn, err := net.ListenUDP("udp", localaddr)
	groupconn, err := net.ListenMulticastUDP("udp", nil, groupaddr)
	comm.CheckError(err)
	
	fmt.Println("CreateSocket end")	

	var b = make([]byte, 256)
	comm.CastData(b, groupconn, localconn, groupaddr)

	localconn.Close()
	groupconn.Close()

	

}
