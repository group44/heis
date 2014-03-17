package types

import (
    "time"
    "sync"
)


const (
	CART_ID int = 1
	N_FLOORS = 4
	N_BUTTONS = 4
	
	// Global timeout const
	TIMEOUT = 1 * time.Second
)


type (
    ElevButtonTypeT int
    OrderTable [][]int
    
    // Map for storing addresses of peers in group
    Peermap struct {
	    Mu sync.Mutex
	    M map[int]time.Time
    }

    // Struct sent over network
    Data struct {
	    Head string
	    Order []int
	    Table [][]int
	    Cost int
	    ID int
	    T time.Time
    }

)








