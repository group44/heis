package order

import (
	"../com"
	"../types"
	"fmt"
	"math/rand"
	"time"
)

func TestDistribute() {

}

// Calculates own cost for an order

/*
func CalculateCost(lt []int, gt [][]int, state elevatorState, order []int) int {
	// Calculate and return an int describing your degree of availability
	// 0 is best

	//Tar ikke hensyn til indre ordre. burde kanskje det, men tror det går greit uansett. litt usikker må teste litt.s

	cost := 666
	elevatorDir := GetOrderDirection()
	elevatorCurrentFloor := GetCurrentFloor()
	orderFloor := waddafakka????
	orderDir := fakkawadda????


	floorDiff = 2*(absolute(elevatorCurrentFloor - orderFloor))  // må lages en absoulteverdi func


	switch state {
		case UP:
			switch elevatorDir {
				case UP:
					if orderDir == DOWN{
						cost = cost + 5
					} else if orderDir == UP {
						cost = cost - 1
					}
					cost = floorDiff + cost
					break
				case DOWN:
					if orderDir == UP {
						cost = cost + 3
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



	if cost<0 {
		return 0
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
		com.OutputCh <- types.Data{Head: "cost", Order: order, Cost: cost, Table: GlobalOrders}

		fmt.Println("Cost calculated:")
		fmt.Println(cost)
		fmt.Println("")
	}

}

// Initiates an "auction" to determine which cart that should dispatch an order.
// Bids in range 0-10, consider changing this - goroutine
func Auction(GlobalOrders types.GlobalTable) {
	var winner int
	var maxCost = 10
	var bid types.Data
	carts := make([]int, types.NUMBER_OF_CARTS)
	fmt.Println(carts)

	for {
		time.Sleep(10 * time.Millisecond)
		bid = <-com.AuctionCh
		com.PeerMap.Mu.Lock()

		fmt.Println(len(carts))
		fmt.Println(len(com.PeerMap.M))

		/*
			for len(carts) < len(com.PeerMap.M)+1 {
				carts[bid.ID-1] = bid.Cost
				//For loop for receiving and processing cost from all elevators
				//bid = <-com.AuctionCh
				fmt.Println("test")
			}
		*/

		carts[bid.ID-1] = bid.Cost

		com.PeerMap.Mu.Unlock()

		// This may cause two or more elevators to claim the same order (if they have equal cost
		for i := 0; i < len(carts); i++ {

			if carts[i] < maxCost {
				maxCost = carts[i]
				winner = i + 1
			}
		}

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
		com.TableCh <- table
		com.OutputCh <- outData

		fmt.Println("Table casted:")
		fmt.Println(outData)
		fmt.Println("")
	}
}

// Removes a successfully dispatched order from the global table
func ClearGlobal(data types.Data) {

}
