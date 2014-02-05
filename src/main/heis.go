package main

import (
	"fmt"
	"../comm"
	"net"
	//"time"
	//"os"
)


// Husk: TimeoutCheck som egen routine
// Peer timeout etter 1 sec?


func main() {
	fmt.Println("CreateSocket begin")

	localaddr, err := net.ResolveUDPAddr("udp4", ":0")
	groupaddr, err := net.ResolveUDPAddr("udp4", "224.0.0.2:12000")
	tesetaddr, err .= net.ResolveDUPAddr("udp4", "12.34.56.78:13000")
	comm.CheckError(err)

	localconn, err := net.ListenUDP("udp", localaddr)
	groupconn, err := net.ListenMulticastUDP("udp", nil, groupaddr)
	comm.CheckError(err)
	
	fmt.Println("CreateSocket end")
	
	// Multicast receive test
	/*
	fmt.Println("Multicast recieve test begin")
	
	a := make([]byte, 1)
	for {
		_, casterAddr, _ := groupconn.ReadFromUDP(a)
		fmt.Println(casterAddr)
	}
	
	fmt.Println("Multicast recieve test end")
	*/
	
	// Map functions-test
	testmap := comm.NewPeermap()
	
	comm.UpdatePeermap(testmap, groupconn)
	fmt.Println(testmap)
	
	fmt.Println(comm.CheckPeerLife(*testmap, groupaddr))

	localconn.Close()
	groupconn.Close()

	
	// Map-test
	/*
	timer := time.NewTimer(1 * time.Second)
	peermap := make(map[string]time.Timer)
	peermap["testaddress"] = *timer
	fmt.Println(peermap["testaddress"])
	*/
	
	




	fmt.Println("Success")
}
