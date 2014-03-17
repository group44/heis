package order

import (
	"../driver"
	//"../com"
	//"sort"
)

// Calculates own cost for an order
func Calculate(lt localTable, gt globalTable, state elevatorState, order []int) int {
	// Calculate and return an int describing your degree of availability
	// 0 is best
	
	cost := 5
	return cost
}

// Initiates an "auction" to determine which cart that should dispatch an order.
// Bids in range 0-20, consider changing this
func Auction(p *com.peerMap, aucch chan order.Data, numberOfCarts int){
	var carts [numberOfCarts]int
	bid := <-aucch
	for len(carts) < len (p) {
		carts[bid.ID] = bid.Cost
		bid = <-aucch
	}

	// This may cause two or more elevators to claim the same order (if they have equal cost)
	winner := 0
	prev := 20 //Arbitrary max value
	for i := 0; i < numberOfCarts; i++ {
		if carts[i] < prev {
			prev = carts[i]
			winner = i
		}
	}
	if winner == CART_ID {
		go Claim(bid.Order)
	}
}


// Claims and order and marks it by setting it's own CART_ID in the ID field of the 
// global table. Should check if another ID is already set, and then not claim it, unless
// the cart who has claimed it is dead.
func Claim() {
}

// Removes a successfully dispatched order from the global table
func Clear() {
}
