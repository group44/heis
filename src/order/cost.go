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
	wfm := 3 //wrong floor multiplier.
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
					cost = cost + wdp
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
		}
	}
	if cost < 0 {
		return 0
	}
	return cost
}

func HandleCost() {
	var cost int
	order := make([]int, 2)
	for {
		order = <-com.OrderCh
		cost = CalculateCost(order)
		fmt.Println("Cost calculated and sent. cost:", cost)
		com.OutputCh <- types.Data{Head: "cost", Order: order, Cost: cost}
	}
}

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

		fmt.Println("len peermap:", len(com.PeerMap.M))
		fmt.Println("len carts", len(carts))
		for !ContainsAll(carts) && len(com.PeerMap.M) > 1 {
			time.Sleep(50 * time.Millisecond)
			//bid = <-com.AuctionCh

			if currentOrder[0] == -1 {
				copy(currentOrder, bid.Order)
			}
			if bid.Order[0] == currentOrder[0] && bid.Order[1] == currentOrder[1] {
				carts[bid.ID-1] = bid.Cost
			}
			if !ContainsAll(carts) && len(com.PeerMap.M) > 1 {
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

		fmt.Println("Winner:", winner)

		if winner == types.CART_ID {
			Claim(bid.Order, types.CART_ID)

		}

	}
}

func Claim(order []int, winner int) {
	floor, dir := order[0], order[1]
	data := types.Data{Head: "addorder"}
	if GlobalOrders[floor][dir] == 0 {
		fmt.Println("order has been claimed. order:", order)
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
