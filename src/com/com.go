package com

import (
	"../types"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	PeerMap = NewPeerMap()

	// Global channels
	OutputCh            = make(chan types.Data)
	OrderCh             = make(chan []int)
	TableCh             = make(chan types.GlobalTable)
	AuctionCh           = make(chan types.Data)
	AddOrderCh          = make(chan types.Data)
	RemoveOrderCh       = make(chan []int)
	UpdateGlobalTableCh = make(chan types.GlobalTable)

	// Local channels
	peerCh = make(chan int)
)

func Run() {

	done := make(chan bool)

	broadcastAddr := "129.241.187.255:12000" // For sanntidssalen
	//broadcastAddr := "localhost:12000" //localhost
	listenAddr := ":12000"

	lAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)
	bAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err)

	lConn, err := net.ListenUDP("udp", lAddr)
	CheckError(err)
	bConn, err := net.DialUDP("udp", nil, bAddr)
	CheckError(err)

	fmt.Println("Sockets created successfully")

	go CastData(bConn)
	go ReceiveData(lConn)
	go UpdatePeerMap(PeerMap)

	<-done
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func NewPeerMap() *types.PeerMap {
	return &types.PeerMap{M: make(map[int]time.Time)}
}

func CheckPeerLife(p types.PeerMap, id int) bool {
	_, present := p.M[id]
	if present {
		tdiff := time.Since(p.M[id])
		return tdiff <= types.TIMEOUT
	}
	return false
}

func UpdatePeerMap(p *types.PeerMap) {
	var id int
	for {
		time.Sleep(100 * time.Millisecond)
		id = <-peerCh
		p.M[id] = time.Now()
		fmt.Println(PeerMap)
	}

}

func ReceiveData(conn *net.UDPConn) {
	var inc types.Data
	//var err error
	var b = make([]byte, 1024)

	for {
		time.Sleep(100 * time.Millisecond)
		n, _, err := conn.ReadFromUDP(b)
		CheckError(err)
		err = json.Unmarshal(b[:n], &inc)
		CheckError(err)
		fmt.Println("Data Received:", inc)

		if inc.ID > 0 {
			// update peermap
			peerCh <- inc.ID
		}

		fmt.Println("Received case:", inc.Head)
		switch inc.Head {

		case "order":
			OrderCh <- inc.Order
			fmt.Println("Order received:", inc.Order)
			fmt.Println("")

		case "table":
			fmt.Println("Table received", inc.Table)

		case "cost":
			fmt.Println(inc)
			AuctionCh <- inc
			fmt.Println("Cost received:", inc.Cost)
			fmt.Println("")

		case "addorder":
			AddOrderCh <- inc
			fmt.Println("Order added:", inc.Order)
			fmt.Println("")

		case "removeorder":
			RemoveOrderCh <- inc.Order
			fmt.Println("Order removed:", inc.Order)
			fmt.Println("")

		default:

			fmt.Println("Default case entered")
		}
	}
}

func CastData(conn *net.UDPConn) {
	var data types.Data
	var err error

	for {
		data = <-OutputCh
		data.ID = types.CART_ID
		data.T = time.Now()
		for i := 0; i < 2; i++ {
			fmt.Println("Data casted:", data)
			b := make([]byte, 1024)
			b, err = json.Marshal(data)
			CheckError(err)
			_, err = conn.Write(b)
			CheckError(err)
			time.Sleep(1 * time.Millisecond)
		}
	}
}
