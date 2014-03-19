package order

import (
	"../types"
	"../com"
)


func TestDistribute() {

}


// Calculates own cost for an order
/*
func Calculate(lt []int, gt [][]int, state elevatorState, order []int) int {
	// Calculate and return an int describing your degree of availability
	// 0 is best
	
	cost := 5
	return cost
}
*/

// Initiates an "auction" to determine which cart that should dispatch an order.
// Bids in range 0-10, consider changing this
func Auction(p *types.PeerMap, aucch chan types.Data, numberOfCarts int) {
	carts := make([]int, numberOfCarts)
	bid := <-aucch
	p.Mu.Lock()
	for len(carts) < len (p.M) {
		carts[bid.ID] = bid.Cost
		bid = <-aucch
	}
	p.Mu.Unlock()

	// This may cause two or more elevators to claim the same order (if they have equal cost)
	winner := 0
	prev := 10 //Arbitrary max value
	for i := 0; i < numberOfCarts; i++ {
		if carts[i] < prev {
			prev = carts[i]
			winner = i
		}
	}
	if winner == types.CART_ID {
		go Claim(GlobalOrders, bid)
	}
}


// Claims and order and marks it by setting it's own CART_ID in the ID field of the 
// global table. Should check if another ID is already set, and then not claim it, unless
// the cart who has claimed it is dead.
func Claim(table types.GlobalTable, data types.Data) { // order: [floor, dir, ID]
	floor, dir := data.Order[0], data.Order[1]
	if table[floor][dir + 1] != 0 {
		table[floor][dir + 1] = types.CART_ID
		outData := types.Data{Head: "table", Table: table}
		com.OutputCh <- outData
	}
}

// Removes a successfully dispatched order from the global table
func ClearGlobal(table types.GlobalTable, data types.Data) {

}


