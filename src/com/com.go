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
	//broadcastAddr := "78.91.39.255:12000"
	//broadcastAddr := "localhost:12000"
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

	//fmt.Println("Channels created succesfully")

	//go ChannelTester()
	go CastData(bConn)
	//fmt.Println("cast")
	go ReceiveData(lConn)
	//fmt.Println("receive")
	go UpdatePeerMap(PeerMap)
	fmt.Println("UpdatePeerMap")

	//testData := types.Data{"cost", []int{1, 0}, [][]int{}, 2, types.CART_ID, time.Now()}
	//AuctionCh <- testData

	/*
		for {
			OutputCh <- testData
			time.Sleep(500 * time.Millisecond)
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

// Creates a map with peer IP as key and time.Time as element
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

		time.Sleep(100 * time.Millisecond)
		id = <-peerCh

		//p.Mu.Lock()
		p.M[id] = time.Now()
		//p.Mu.Unlock()

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
		fmt.Println("her er det vi har motatt:")
		fmt.Println(inc)

		//fmt.Println("in:")
		//fmt.Println(inc)
		if inc.ID > 0 {
			// update peermap
			peerCh <- inc.ID // c1
		}

		fmt.Println(inc.Head)
		switch inc.Head {

		case "order":
			OrderCh <- inc.Order
			fmt.Println("Order received:")
			//fmt.Println(inc.Order)
			fmt.Println("")

		case "table":
			if inc.ID != types.CART_ID {
				/*
					fmt.Println("Table received")
					//fmt.Println(inc.Table)
					fmt.Println("")
					//TableCh <- inc.Table //
					fmt.Println("kommer vi hertil??????")
				*/
			}

		case "cost":
			fmt.Println(inc)
			AuctionCh <- inc
			//fmt.Println(inc)
			//fmt.Println(AuctionCh)
			fmt.Println("Cost received:")
			fmt.Println(inc.Cost)
			fmt.Println("")

		case "addorder":
			AddOrderCh <- inc
			fmt.Println("Order added:")
			fmt.Println(inc.Order)
			fmt.Println("")

		case "removeorder":
			RemoveOrderCh <- inc.Order
			fmt.Println("Order removed:")
			fmt.Println(inc.Order)
			fmt.Println("")

		default:

			fmt.Println("Default case entered")

		}

		//fmt.Println(inc)
	}

}

func CastData(conn *net.UDPConn) {
	var data types.Data
	var err error

	for {
		data = <-OutputCh
		data.ID = types.CART_ID
		data.T = time.Now()
		fmt.Println("out:")
		fmt.Println("Dette er det vi caster!")
		fmt.Println(data)
		b := make([]byte, 1024)
		b, err = json.Marshal(data)
		CheckError(err)
		_, err = conn.Write(b)
		CheckError(err)

	}

}

/*
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

*/
