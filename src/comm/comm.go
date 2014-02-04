package comm

import (
	"net"
	"fmt"
	//"encoding/json"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}


// leser alle incoming packets med ReadFromUDP og sorterer dem etter innholder/lengde.
// bestemme om en selv er doed ved aa sjekke timeout paa ReadFromMulticast

// Struct for storing addresses of peers in group
/*
type peermap struct {
	m map[string]
}

// Creates new list of peers
func NewPeerlist() *peerlist {
	return &peerlist{s: make([]string, 0)}
}

// Todo: Make map entry change after timeout. Are "countdowns" allowed as entries in map?
func UpdatePeerlist(lst *peerlist, a net.Addr) {
	
}

// Adds an address to a peerlist
func AddToPeerlist(lst *peerlist, a net.Addr) {
	lst.s = append(lst.s, a.String())	
}
*/


// Creates multicast-listener for UDP packages in a multicast group. Error check for network operations -- rename to JoinGroup()?

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

func CastData(data []byte, conn, lconn *net.UDPConn, gaddr *net.UDPAddr) {
	fmt.Println("CastData begin")
	_, _ = lconn.WriteToUDP(data, gaddr)
	fmt.Println("CastData end")
}


// Updates an array/slice with the address of the clients sending out multicast packet and whether or not it's alive. If clientlist does not exist -> create. If client already in list -> do nothing. If client reappears after timeout -> do not make a new entry to the list. Instead set it's status as alive.
//func UpdateClientList()

//func SortMulticastPacket(



