package main

import (
	"fmt"
	"../comm"
	//"../cost"
	//"../elev"
	"net"
	//"time"
	//"os"
)



func main() {

	laddr, err := net.ResolveUDPAddr("udp4", ":0")
	gaddr, err := net.ResolveUDPAddr("udp4", "224.0.0.1:12000")
	//gaddr, err := net.ResolveUDPAddr("udp4", "129.241.187.255:12000")
	testaddr1, err := net.ResolveUDPAddr("udp4", "129.241.187.148:0")
	testaddr2, err := net.ResolveUDPAddr("udp4", "129.241.187.144:0")
	comm.CheckError(err)

	gconn, err := net.ListenUDP("udp4", gaddr)	
	lconn, err := net.ListenUDP("udp4", laddr)	
	//groupconn, err := net.ListenMulticastUDP("udp", nil, groupaddr)
	comm.CheckError(err)
	
	fmt.Println("Sockets created successfully")
	
	var data = make([]byte, 512)
	
	for {	
		comm.ReceiveTest(gconn, data)
	}


	// Map functions-test
	testmap := comm.NewPeermap()
	
	
	go comm.UpdatePeermap(testmap, gconn)
	
	for {
		fmt.Println(comm.CheckPeerLife(*testmap, testaddr1))
		fmt.Println()
		fmt.Println(comm.CheckPeerLife(*testmap, testaddr2))
	}
	
	
	
	
	//Test of sending/recieving JSON
	/*
	testOrder := elev.NewOrder(2, 1)
	fmt.Println(testOrder)
	
	for {
		go comm.ReceiveData(groupconn)
	}
	comm.CastData(testOrder, groupconn, localconn, groupaddr)
	*/
	

	lconn.Close()
	gconn.Close()

	fmt.Println("End")
}
