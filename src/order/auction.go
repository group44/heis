package auction

import (
	"../types"
	"fmt"
)


func Auction() {
	var winner int
	var lowestCost = 100
	var currentOrder = make([]int, 2)
	currentOrder[0], currentOrder[1] = -1, -1
	var auctionMap = make(map[int]int)
	var bidder types.Data

	for {

		for cart := range PeerMap {
			auctionMap[cart] = PeerMap.M[cart]
		}

		for {
			time.Sleep(50 * time.Millisecond)

			select {
				case bidder = <- AuctionCh:
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

		fmt.Println("WINNER:")
		fmt.Println(winner)

		for cart := range auctionMap {
			delete(auctionMap, cart)
		}

	}
}


/*
func ContainsAll(carts []int) bool {
	for _, t := range carts {
		if t == 0 {
			return false
		}
	}
	return true
}
*/

func NewPeerMap() *types.PeerMap {
	return &types.PeerMap{M: make(map[int]time.Time)}
}