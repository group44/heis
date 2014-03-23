package com

import (
	"../types"
	"fmt"
	"net"
	//"sync"
	//"../order"
	//"encoding/gob"
	"encoding/json"
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
	AddOrderCh          = make(chan types.Data, 5)
	RemoveOrderCh       = make(chan []int, 5)
	UpdateGlobalTableCh = make(chan types.GlobalTable)

	// Local channels
	peerCh = make(chan int)
)

// Create sockets and start go routines
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

// Error check
func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// Creates a map with peer IP as key and time.Time as element
func NewPeerMap() *types.PeerMap {
	return &types.PeerMap{M: make(map[int]time.Time)}
}

// Checks if peer address is in peer map and time difference is not > 1 sec
func CheckPeerLife(p types.PeerMap, id int) bool {
	_, present := p.M[id]
	if present {
		tdiff := time.Since(p.M[id])
		return tdiff <= types.TIMEOUT
	}
	return false
}

// Updates peermap and sets discovery time from conn input
func UpdatePeerMap(p *types.PeerMap) {
	var id int
	for {
		time.Sleep(100 * time.Millisecond)
		id = <-peerCh
		p.M[id] = time.Now()
		fmt.Println(PeerMap)
	}

}

// Listens and receives from connection in seperate go-routine
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

// Go routine for sending data over udp
// Sends all messeage 2 times to make sure they get through
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
