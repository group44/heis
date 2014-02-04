package main

import (
	"fmt"
	"../comm"
	"net"
	"time"
	//"os"
)


// Husk: TimeoutCheck som egen routine


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

	
	// Map-test
	
	timer := time.NewTimer(1 * time.Second)
	peermap := make(map[string]time.Timer)
	peermap["testaddress"] = *timer
	fmt.Println(peermap["testaddress"])
	
	




	fmt.Println("Success")
}
