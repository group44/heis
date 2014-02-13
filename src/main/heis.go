package main

import (
	"fmt"
	"../comm"
	//"../cost"
	"../elev"
	"net"
	//"time"
	//"os"
)



func main() {

	//testaddr1, err := net.ResolveUDPAddr("udp4", "129.241.187.152:0")
	//testaddr2, err := net.ResolveUDPAddr("udp4", "129.241.187.142:0")
	//comm.CheckError(err)

	broadcastAddr := "129.241.187.255:12000"
	listenAddr := ":12000"

	laddr, err := net.ResolveUDPAddr("udp", listenAddr)
	baddr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	comm.CheckError(err)

	lconn, err := net.ListenUDP("udp", laddr)
	bconn, err := net.DialUDP("udp", nil, baddr)
	comm.CheckError(err)
	
	fmt.Println("Sockets created successfully")
	
	//var data = make([]byte, 1)
	
	/*	
	for {	
		comm.ReceiveTest(gconn, data)
	}
	*/

	// Map functions-test
	/*
	testmap := comm.NewPeermap()
	
	
	go comm.UpdatePeermap(testmap, lconn)
	
	for {
		fmt.Println(comm.CheckPeerLife(*testmap, testaddr1))
		fmt.Println(comm.CheckPeerLife(*testmap, testaddr2))
		time.Sleep(time.Second)
		fmt.Println()
	}
	*/
	
	
	//Test of sending/recieving JSON
	
	testOrder := elev.NewOrder(2, 1)
	//b := make([]byte, 512)
	//fmt.Println(testOrder)
	
	for {
	
		comm.CastData(testOrder, bconn)	
		//comm.ReceiveTest(lconn, b)
		comm.ReceiveData(lconn)
		fmt.Println()
	}	
	
	

	lconn.Close()
	bconn.Close()

	fmt.Println("End")
}
