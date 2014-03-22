package order

import (
	"../com"
	"../types"
	"fmt"
	"math/rand"
	//"math"
	"time"
)

func TestDistribute() {

}

/*
// Calculates own cost for an order får ikke inn state, så den må fjernes hvis vi ikke finner ut av det..
func CalculateCost(lt []int, gt [][]int, state int, order []int) int {
	cost := 0
	elevatorDir := GetOrderDirection()
	elevatorCurrentFloor := GetCurrentFloor()
	orderFloor := order[0]
	orderDir := order[1]
	wdp := 10 //wrong direction punishment
	rdr := 3  //right direction reward
	wfm := 2  //wrong floor multiplier.

	floorDiff := wfm * (int(math.Abs(float64(elevatorCurrentFloor - orderFloor))))

	switch state {
	case UP:
		switch elevatorDir {
		case UP:
			if orderDir == DOWN {
				cost = cost + wdp
			} else if orderDir == UP {
				cost = cost - rdr
			}
			cost = floorDiff + cost
			break
		case DOWN:
			if orderDir == UP {
				cost = cost + wdp
			}
			cost = floorDiff + cost
			break
		}
		break

	case DOWN:
		switch elevatorDir {
		case UP:
				if orderDir == DOWN{
				cost = cost + 3
			}
			cost = floorDiff + cost
			break
		case DOWN:
			if orderDir == UP {
				cost = cost +5
			} else if orderDir == DOWN {
				cost = cost - 1
			}
			cost = floorDiff + cost
			break
		case DOWN:
			if orderDir == UP {
				cost = cost + wdp
			} else if orderDir == DOWN {
				cost = cost - rdr
			}
			break
	default:
		switch elevatorDir {
		case UP:
			if orderDir == DOWN{
				cost = cost + 5
			}
			cost = floorDiff + cost
			break
		case DOWN:
			if orderDir == UP {
				cost = cost +5
			}
			cost = floorDiff + cost
			break
		}
	break
	}
	}
	return cost
}
*/

// Test function, gives random cost, goroutine

func CalculateCost() {
	var cost int
	order := make([]int, 2)

	for {
		order = <-com.OrderCh
		// Temp
		cost = rand.Intn(10) + 1
		// Temp

		fmt.Println("Cost calculated and sent")
		//fmt.Println(cost)
		//fmt.Println("")

		com.OutputCh <- types.Data{Head: "cost", Order: order, Cost: cost}
	}

}

/*
func CalculateCost() {
	var order = make([]int, 2)
	order = <- com.OrderCh
	com.OutputCh <- types.Data{Head:"cost", Order: order, Cost: 1}
}
*/

// Initiates an "auction" to determine which cart that should dispatch an order.
// Bids in range 0-10, consider changing this - goroutine
func Auction(GlobalOrders types.GlobalTable) {
	var winner int
	var maxCost = 10
	var bid types.Data
	var currentOrder = make([]int, 2)
	currentOrder[0], currentOrder[1] = -1, -1
	carts := make([]int, types.NUMBER_OF_CARTS)

	for {

		//fmt.Println(currentOrder)
		time.Sleep(10 * time.Millisecond)
		bid = <-com.AuctionCh

		if currentOrder[0] == -1 {
			copy(currentOrder, bid.Order)
		}

		com.PeerMap.Mu.Lock()
		for len(carts) < len(com.PeerMap.M)+1 {
			time.Sleep(50 * time.Millisecond)
			if bid.Order[0] == currentOrder[0] && bid.Order[1] == currentOrder[1] {
				carts[bid.ID-1] = bid.Cost
			}
			bid = <-com.AuctionCh

		}

		currentOrder[0] = -1
		com.PeerMap.Mu.Unlock()

		// This may cause two or more elevators to claim the same order (if they have equal cost
		for i := 0; i < len(carts); i++ {
			if carts[i] < maxCost {
				maxCost = carts[i]
				winner = i + 1
			}
		}

		fmt.Println("WINNER:")
		fmt.Println(winner)

		if winner == types.CART_ID {
			Claim(bid.Order, GlobalOrders)

		}

	}
}

// Claims and order and marks it by setting it's own CART_ID in the ID field of the
// global table. Should check if another ID is already set, and then not claim it, unless
// the cart who has claimed it is dead.
func Claim(order []int, table types.GlobalTable) { // order: [floor, dir, ID]
	floor, dir := order[0], order[1]
	if table[floor][dir] == 0 {
		table[floor][dir] = types.CART_ID
		outData := types.Data{Head: "table", Table: table}
		fmt.Println("Sending table from Claim")
		com.TableCh <- table
		fmt.Println("Sending Output from Claim")
		com.OutputCh <- outData

		/*
			fmt.Println("Table casted:")
			fmt.Println(outData)
			fmt.Println("")
		*/
	}
}

// Removes a successfully dispatched order from the global table
func ClearGlobal(data types.Data) {

}
