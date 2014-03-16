package main

import (
	"fmt"
	"../com"
	"net"
	//"time"
	//"os"
)



func main() {

	const CART_ID int = 1

	//testaddr1, err := net.ResolveUDPAddr("udp4", "129.241.187.152:0")
	//testaddr2, err := net.ResolveUDPAddr("udp4", "129.241.187.142:0")
	//com.CheckError(err)

	//broadcastAddr := "129.241.187.255:12000" // For sanntidssalen
	listenAddr := ":12000"

	laddr, err := net.ResolveUDPAddr("udp", listenAddr)
	baddr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	com.CheckError(err)

	lconn, err := net.ListenUDP("udp", laddr)
	bconn, err := net.DialUDP("udp", nil, baddr)
	com.CheckError(err)
	
	fmt.Println("Sockets created successfully")
	
	//var data = make([]byte, 1)
	
	/*	
	for {	
		com.ReceiveTest(gconn, data)
	}
	*/

	// Map functions-test
	/*
	testmap := com.NewPeermap()
	
	
	go com.UpdatePeermap(testmap, lconn)
	
	for {
		fmt.Println(com.CheckPeerLife(*testmap, testaddr1))
		fmt.Println(com.CheckPeerLife(*testmap, testaddr2))
		time.Sleep(time.Second)
		fmt.Println()
	}
	*/
	
	
	var testOrder interface{} = data.NewOrder(0, 1, 1)
	//var testTable data.OrderTable
	fmt.Println("testOrder created")
	//fmt.Println(testOrder)
	
	fmt.Println("testOrder sent:")
	com.CastData(testOrder, bconn)
	fmt.Println("testOrder received:")
	com.ReceiveData(lconn)	
	fmt.Println()
		
	//data.ClaimOrder(testOrder, &testTable)
	//fmt.Println(testTable)
	

	lconn.Close()
	bconn.Close()

	fmt.Println("End")
}
