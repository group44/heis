package main

import (
	"fmt"
	"net"
	"../com"
	"../order"
	//"time"
)



func main() {

	const CART_ID int = 1


	broadcastAddr := "129.241.187.255:12000" // For sanntidssalen
	listenAddr := ":12000"
	//sendingAddr := ":12001"

	//saddr, err := net.ResolveUDPAddr("udp", sendingAddr)
	laddr, err := net.ResolveUDPAddr("udp", listenAddr)
	baddr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	com.CheckError(err)

	//sconn, err := net.ListenUDP("udp", saddr)
	lconn, err := net.ListenUDP("udp", laddr)
	bconn, err := net.DialUDP("udp", nil, baddr)
	com.CheckError(err)
	fmt.Println("Sockets created successfully")

	
	peerch := make(chan net.Addr)
	orderch := make(chan []int)
	tablech := make(chan [][]int)
	aucch := make(chan int)
	fmt.Println("Channels created succesfully")
	
	testOrder := order.Data{"order", []int{1, 0, 1}, [][]int{}, 2}

	
	go com.CastData(testOrder, bconn)
	go com.ChannelTester(peerch, orderch, tablech, aucch)
	com.ReceiveData(lconn, peerch, orderch, tablech, aucch)

	
		
	//data.ClaimOrder(testOrder, &testTable)
	//fmt.Println(testTable)
	

	//lconn.Close()
	//bconn.Close()

	fmt.Println("End")
}