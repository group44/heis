package com

import (
	"../types"
	"fmt"
	"net"
	//"sync"
	//"../order"
	"encoding/gob"
	"os"
	"time"
)

var PeerMap = NewPeerMap()

// Global channels
var OutputCh = make(chan types.Data, 5)
var OrderCh = make(chan []int)
var TableCh = make(chan [][]int)
var AuctionCh = make(chan types.Data)

// Local channels
var peerCh = make(chan int)

// Create sockets and start go routines
func Run() {

	done := make(chan bool)

	broadcastAddr := "129.241.187.255:12000" // For sanntidssalen
	//broadcastAddr := "78.91.39.255:12000"
	listenAddr := ":12000"

	lAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	bAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err)

	lConn, err := net.ListenUDP("udp", lAddr)
	bConn, err := net.DialUDP("udp", nil, bAddr)
	CheckError(err)

	fmt.Println("Sockets created successfully")

	//fmt.Println("Channels created succesfully")

	//testData := types.Data{"cost", []int{1, 0, 1}, [][]int{}, 2, types.CART_ID, time.Now()}

	//go ChannelTester()
	go CastData(bConn)
	//fmt.Println("cast")
	go ReceiveData(lConn)
	//fmt.Println("receive")
	go UpdatePeerMap(PeerMap)
	//fmt.Println("UpdatePeerMap")

	/*
		for {

			OutputCh <- testData
			time.Sleep(2 * time.Second)
		}
	*/

	<-done
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// Creates a map with peer IP as key and timer as element
func NewPeerMap() *types.PeerMap {
	return &types.PeerMap{M: make(map[int]time.Time)}
}

/*
func GetPeerMapLength() int {
	return len(PeerMap)
}
*/

// Checks if peer address is in peer map and time difference is not > 1 sec
// New version using ID instead of peeraddr
func CheckPeerLife(p types.PeerMap, id int) bool {
	_, present := p.M[id]
	if present {
		tdiff := time.Since(p.M[id])
		return tdiff <= types.TIMEOUT
	}
	return false
}

// Updates peermap and sets discovery time from conn input
// New version using ID instead of IP
func UpdatePeerMap(p *types.PeerMap) {
	var id int

	for {
		id = <-peerCh
		if id >= 0 && id != types.CART_ID {
			p.Mu.Lock()
			p.M[id] = time.Now()
			p.Mu.Unlock()
			//fmt.Println("peer")
		}

	}

}

// Listens and receives from connection in seperate go-routine
func ReceiveData(conn *net.UDPConn) {
	var inc types.Data
	decoder := gob.NewDecoder(conn)

	for {
		err := decoder.Decode(&inc)
		CheckError(err)

		if inc.ID == types.CART_ID {
			continue
		}
		if inc.ID > 0 {
			// update peermap
			peerCh <- inc.ID // c1
		}

		if inc.Head == "order" {
			OrderCh <- inc.Order // c2
			fmt.Println(inc.Order)
		} else if inc.Head == "table" {

			TableCh <- inc.Table // c3 - is this channel needed?
		} else if inc.Head == "cost" {

			AuctionCh <- inc // c4
		}

		//fmt.Println(inc)
	}

}

func CastData(conn *net.UDPConn) {
	var data types.Data
	encoder := gob.NewEncoder(conn)

	for {
		data = <-OutputCh
		data.ID = types.CART_ID
		data.T = time.Now() // Sets timestamp on outgoing data
		err := encoder.Encode(data)
		CheckError(err)
		fmt.Println(data)
	}

}

func ChannelTester(c1 chan int, c2 chan []int, c3 chan [][]int, c4 chan int) {

	for {
		select {
		case <-c1:
			fmt.Println("c1 read")
		case <-c2:
			fmt.Println("c2 read")
		case <-c3:
			fmt.Println("c3 read")
		case <-c4:
			fmt.Println("c4 read")
		}
	}

}
