package main

import (
	"fmt"
	"net"
	"../types"
	"../com"
	//"../order"
	//"time"
)



func main() {

	const CART_ID int = 0

	//broadcastAddr := "129.241.187.255:12000" // For sanntidssalen
	listenAddr := ":12000"

	laddr, err := net.ResolveUDPAddr("udp", listenAddr)
	//baddr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	com.CheckError(err)

	lconn, err := net.ListenUDP("udp", laddr)
	//bconn, err := net.DialUDP("udp", nil, baddr)
	com.CheckError(err)
	fmt.Println("Sockets created successfully")

	
	peerch := make(chan int)
	orderch := make(chan []int)
	tablech := make(chan [][]int)
	aucch := make(chan int)
	fmt.Println("Channels created succesfully")
	
	testMap := com.NewPeerMap()
	//testOrder := order.Data{"order", []int{1, 0, 1}, [][]int{}, 2, time.Now()}

	go com.UpdatePeerMap(testMap, CART_ID, peerch) 
	//go com.CastData(testOrder, bconn)
	go com.ChannelTester(peerch, orderch, tablech, aucch)
	com.ReceiveData(lconn, peerch, orderch, tablech, aucch)

	
		
	//data.ClaimOrder(testOrder, &testTable)
	//fmt.Println(testTable)
	

	//lconn.Close()
	//bconn.Close()

	fmt.Println("End")
}
