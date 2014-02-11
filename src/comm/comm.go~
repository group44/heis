package comm

import (
	"net"
	"fmt"
	"sync"
	"time"
	"encoding/json"
)

var timeout time.Duration = 1 * time.Second

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}


// leser alle incoming packets med ReadFromUDP og sorterer dem etter innholder/lengde.
// bestemme om en selv er doed ved aa sjekke timeout paa ReadFromMulticast

// Map for storing addresses of peers in group
type peermap struct {
	mu sync.Mutex
	m map[string]time.Time
}

// Creates a map with peer IP as key and timer as element
func NewPeermap() *peermap {
	return &peermap{m: make(map[string]time.Time)}
}

// Checks if peer address is in peer map and time difference is not > 1 sec
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


// Updates peermap and sets discovery time from conn input
func UpdatePeermap(p *peermap, conn *net.UDPConn) {
	for {
		var buf [1024]byte
		_, addr, err := conn.ReadFrom(buf[:])
		CheckError(err)
		peeraddr, _, err := net.SplitHostPort(addr.String())
		CheckError(err)
		fmt.Println(peeraddr)
		p.mu.Lock()
		p.m[peeraddr] = time.Now()
		p.mu.Unlock()
	}
}


// This is in main() for now
func CreateSocket() {
}

// Receive data from multicast socket. Returns number of bytes read and the return address of the packet. Can be made to timeout and return an error after a fixed time limit; see SetDeadline and SetReadDeadline.
func ReceiveData(conn *net.UDPConn) ([]byte) {
	fmt.Println("ReceiveData begin")
	data := make([]byte, 256)
	_, _, err := conn.ReadFromUDP(data)
	CheckError(err)
	fmt.Println("read", string(data))
	fmt.Println("ReceiveData end")
	return data
}

// Testing JSON
func CastData(data struct, conn *net.UDPConn, lconn *net.UDPConn, gaddr *net.UDPAddr) {
	fmt.Println("CastData begin")
	b, _ := json.Marshal(data)
	fmt.Println(b)
	//_, err := lconn.WriteToUDP(data, gaddr)
	//CheckError(err)
	fmt.Println("CastData end")
}



