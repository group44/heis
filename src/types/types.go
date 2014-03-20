package types

import (
    "time"
    "sync"
)


const (
	CART_ID = 1
	N_FLOORS = 4
	N_BUTTONS = 4
	NUMBER_OF_CARTS = 2
	
	// Global timeout const
	TIMEOUT = 1 * time.Second
)

var (

)


type (

	GlobalTable [][]int
	InternalTable []int

    ElevButtonTypeT int
    
    // Map for storing addresses of peers in group
    PeerMap struct {
	    Mu sync.Mutex
	    M map[int]time.Time
    }

    // Struct for sending data over network
    Data struct {
	    Head string
	    Order []int
	    Table [][]int
	    Cost int
	    ID int
	    T time.Time
    }

)


func NewGlobalTable() [][]int {
	t := make([][]int, N_FLOORS)
	for i := range t {
		t[i] = make([]int, 4)
		for j := range t[i] {
			t[i][j] = 0
		}
	}
	return t
}

func NewInternalTable() []int {
	t := make([]int, N_FLOORS)
	return t
}


