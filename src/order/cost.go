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

func Auction() {
	var winner int
	var lowestCost = 100
	var currentOrder = make([]int, 2)
	currentOrder[0], currentOrder[1] = -1, -1
	var auctionMap = make(map[int]int)
	var bidder types.Data

	for {
		for cart := range PeerMap.M {
			auctionMap[cart] = PeerMap.M[cart]
		}
		for {
			time.Sleep(50 * time.Millisecond)
			select {
			case bidder = <-AuctionCh:
				if currentOrder[0] == -1 {
					copy(currentOrder, bidder.Order)
				}
				auctionMap[bidder.ID] = bidder.Cost
			default:
				break
			}

			fmt.Println(carts)
		}

		currentOrder[0] = -1

		for cart := range auctionMap {
			if auctionMap[cart] < lowestCost {
				winner = cart
				lowestCost = auctionMap[cart]
			}
		}

		fmt.Println("Winner:", winner)

		for cart := range auctionMap {
			delete(auctionMap, cart)
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
