package order

import (
	"../types"
	"../com"
	"math/rand"
)


func TestDistribute() {

}


// Calculates own cost for an order
/*
func CalculateCost(lt []int, gt [][]int, state elevatorState, order []int) int {
	// Calculate and return an int describing your degree of availability
	// 0 is best
	
	cost := 5
	return cost
}
*/

// Test function, gives random cost, goroutine
func CalculateCost() {
	var cost int
	order := make([]int, 3)

	for {
		order = <- com.OrderCh
		// Temp
		cost = rand.Intn(10) + 1
		// Temp
		com.OutputCh <- types.Data{Head:"cost", Order:order, Cost:cost}
	}

}

// Initiates an "auction" to determine which cart that should dispatch an order.
// Bids in range 0-10, consider changing this - goroutine
func Auction(GlobalOrders types.GlobalTable) {
	var winner int
	var maxCost = 10
	var bid types.Data
	carts := make([]int, types.NUMBER_OF_CARTS)

	for {
		bid = <- com.AuctionCh
		com.PeerMap.Mu.Lock()
		for len(carts) < len (com.PeerMap.M) {
			carts[bid.ID] = bid.Cost
			bid = <-com.AuctionCh
		}
		com.PeerMap.Mu.Unlock()

		// This may cause two or more elevators to claim the same order (if they have equal cost
		for i := 0; i < types.NUMBER_OF_CARTS; i++ {
			if carts[i] < maxCost {
				maxCost = carts[i]
				winner = i + 1
			}
		}
		if winner == types.CART_ID {
			Claim(bid, GlobalOrders)
		}
	}

}


// Claims and order and marks it by setting it's own CART_ID in the ID field of the 
// global table. Should check if another ID is already set, and then not claim it, unless
// the cart who has claimed it is dead.
func Claim(data types.Data, table types.GlobalTable) { // order: [floor, dir, ID]
	floor, dir := data.Order[0], data.Order[1]
	if table[floor][dir + 1] != 0 {
		table[floor][dir + 1] = types.CART_ID
		outData := types.Data{Head:"table", Table:table}
		com.OutputCh <- outData
	}
}

// Removes a successfully dispatched order from the global table
func ClearGlobal(data types.Data) {

}