package com

import (
	"../types"
	"net"
	"fmt"
	//"sync"
	"time"
	"encoding/gob"
	"os"
	//"../order"
)


func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// Creates a map with peer IP as key and timer as element
func NewPeerMap() *types.Peermap {
	return &types.Peermap{M: make(map[int]time.Time)}
}

// Checks if peer address is in peer map and time difference is not > 1 sec
/*
func CheckPeerLife(p peermap, addr net.Addr) bool {
	peeraddr, _, err := net.SplitHostPort(addr.String())
	//fmt.Println(peeraddr)
	CheckError(err)
	_, present := p.m[peeraddr]
	if present {
		p.mu.Lock()
		tdiff := time.Since(p.m[peeraddr])
		p.mu.Unlock()
		return tdiff <= timeout
	}
	return false
}
*/

// New version using ID instead of peeraddr
func CheckPeerLife(p types.Peermap, id int) bool {
	_, present := p.M[id]
	if present {
		tdiff := time.Since(p.M[id])
		return tdiff <= types.TIMEOUT
	}
	return false
}


// Updates peermap and sets discovery time from conn input
/*
func UpdatePeermap(p *peermap, conn *net.UDPConn) {
	for {
		var buf [1024]byte
		_, addr, err := conn.ReadFromUDP(buf[:])
		CheckError(err)
		peeraddr, _, err := net.SplitHostPort(addr.String())
		CheckError(err)
		p.mu.Lock()
		p.m[peeraddr] = time.Now()
		p.mu.Unlock()
	}
}
*/

// New version using ID instead of IP
func UpdatePeerMap(p *types.Peermap, id int, peerch chan int) {
	for {
		<- peerch
		p.Mu.Lock()
		p.M[id] = time.Now()
		p.Mu.Unlock()
	}
}


// This is in main() for now
func CreateSocket() {
}

// Receive data from multicast socket. Returns number of bytes read and the return address of the packet. Can be made to timeout and return an error after a fixed time limit; see SetDeadline and SetReadDeadline.

// Listens and receives from connection in seperate go-routine
func ReceiveData(conn *net.UDPConn, peerch chan int, orderch chan []int, tablech chan [][]int, aucch chan int) {
    fmt.Println("test1")
	decoder := gob.NewDecoder(conn)
	for {
		inc := types.Data{}
		err := decoder.Decode(&inc)
		fmt.Println("test2")
		CheckError(err)
		fmt.Println(err)
		// update peermap
		peerch <- inc.ID
		
		if inc.Head == "order" {
			orderch <- inc.Order
		} else if inc.Head == "table" {
			tablech <- inc.Table
		} else if inc.Head == "cost" {
			aucch <- inc.Cost
		}
		
		fmt.Println(inc)
	}
}

func CastData(d types.Data, conn *net.UDPConn) {
	encoder := gob.NewEncoder(conn)
	for {
		err := encoder.Encode(d)
		CheckError(err)
		//fmt.Println(d)
	}
}

func ReceiveTest(c *net.UDPConn, b []byte) {
	c.ReadFromUDP(b)
	fmt.Println(b)
}

func ChannelTester(c1 chan int, c2 chan []int, c3 chan [][]int, c4 chan int) {
	for {
		select {
		case <- c1:
			fmt.Println("c1 read")
		case <- c2:
			fmt.Println("c2 read")
		case <- c3:
			fmt.Println("c3 read")
		case <- c4:
			fmt.Println("c4 read")
		}
	}
}

// Create sockets and start go routines
func Init() {

}


