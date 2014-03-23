package order

import (
	"../com"
	"../driver"
	"../types"
	"fmt"
	//"math/rand"
	"math"
	"time"
)

func Test() {
	driver.IoReadBit(driver.MOTORDIR)
}

//Calculate the elevators cost for the order
func CalculateCost(order []int) int {
	cost := 0
	var state int
	var motorDir int
	elevatorDir := GetOrderDirection()
	elevatorCurrentFloor := GetCurrentFloor()
	orderFloor := order[0]
	orderDir := order[1]
	wdp := 3 //wrong direction punishment
	rdr := 2 //right direction reward
	wfm := 2 //wrong floor multiplier.

	floorDiff := wfm * (int(math.Abs(float64(elevatorCurrentFloor - orderFloor))))

	if driver.ElevGetFloorSensorSignal() == -1 {
		motorDir = driver.IoReadBit(driver.MOTORDIR)
		if motorDir == 1 {
			state = UP
		} else if motorDir == 0 {
			state = DOWN
		}
	}

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
			if orderDir == DOWN {
				cost = cost + wdp
			}
			cost = floorDiff + cost
			break
		case DOWN:
			if orderDir == UP {
				cost = cost + wdp
			} else if orderDir == DOWN {
				cost = cost - rdr
			}
			cost = floorDiff + cost
			break

		default:
			switch elevatorDir {
			case UP:
				if orderDir == DOWN {
					cost = cost + 5
				}
				cost = floorDiff + cost
				break
			case DOWN:
				if orderDir == UP {
					cost = cost + 5
				}
				cost = floorDiff + cost
				break
			}
			break
		}
	}
	if cost < 0 {
		return 0
	}
	return cost
}

//function that handles cost enquireieieii
func HandleCost() {
	var cost int
	order := make([]int, 2)

	for {
		order = <-com.OrderCh
		cost = CalculateCost(order)
		fmt.Println("Cost calculated and sent")
		com.OutputCh <- types.Data{Head: "cost", Order: order, Cost: cost}
	}

}

// Initiates an "auction" to determine which cart that should dispatch an order.
// Bids in range 0-10, consider changing this - goroutine
func Auction(GlobalOrders types.GlobalTable) {
	var winner int
	var maxCost = 100
	var bid types.Data
	var currentOrder = make([]int, 2)
	currentOrder[0], currentOrder[1] = -1, -1
	carts := make([]int, types.NUMBER_OF_CARTS)

	for {

		//fmt.Println(currentOrder)
		//hva skal komme inn her? går det til cost først eller til denne først?
		bid = <-com.AuctionCh
		com.PeerMap.Mu.Lock()

		//hva gjør egentlig denne??

		fmt.Println("len peermap")
		fmt.Println(len(com.PeerMap.M))
		fmt.Println("len carts")
		fmt.Println(len(carts))

		for !ContainsAll(carts) {
			time.Sleep(50 * time.Millisecond)
			//bid = <-com.AuctionCh

			if currentOrder[0] == -1 {
				copy(currentOrder, bid.Order)
			}
			if bid.Order[0] == currentOrder[0] && bid.Order[1] == currentOrder[1] {
				carts[bid.ID-1] = bid.Cost
			}
			if !ContainsAll(carts) {
				fmt.Println("1111111111111111")
				bid = <-com.AuctionCh
			}
			fmt.Println(carts)
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
			Claim(bid.Order, types.CART_ID)

		}

	}
}

// Claims and order and marks it by setting it's own CART_ID in the ID field of the
// global table. Should check if another ID is already set, and then not claim it, unless
// the cart who has claimed it is dead.
func Claim(order []int, winner int) { // order: [floor, dir, ID]
	floor, dir := order[0], order[1]
	data := types.Data{Head: "addorder"}
	if GlobalOrders[floor][dir] == 0 {
		fmt.Println("order has been claimed")
		data.Order = order
		data.WinnerId = winner
		com.OutputCh <- data
	}
}

func ContainsAll(carts []int) bool {
	for _, t := range carts {
		if t == 0 {
			return false
		}
	}
	return true
}
