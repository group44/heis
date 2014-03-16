package order

import (
	"../driver"
	"fmt"
)

const (
	CART_ID int = 1
	N_FLOORS = 4
	N_BUTTONS = 4
	
	UP = 0
	DOWN = 1
	INTERNAL = 2
)

var (
	localTable [N_FLOORS][3]int
	type GlobalTable [N_FLOORS][2]int
)


// Todo, check for errors in order type functions
func CheckError(err) {
}

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne
func UpdateLocalTable() {
	for i := 0; i < N_FLOORS; i++ {
		localTable[i][INTERNAL] = driver.ElevGetButtonSignal(INTERNAL, i)
	}
}

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne
func RemoveOrder() {
	floor := driver.ElevGetFloorSensorSignal() 
	dir := driver.ElevGetDirection()
	if floor != -1 && floor < N_FLOORS {
		localTable[floor][INTERNAL] = 0
		localTable[floor][dir] = 0
	}
}

// Vurder assert, tar ikke hensyn til retning
func CheckCurrentFloor() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	if currentFloor != -1 && currentFloor < N_FLOORS {
		return (localTable[currentFloor][0] == 1 || localTable[currentFloor][1] == 1 || localTable[currentFloor][2] == 1)
	}
	return false
}

// For enkel, returnerer bare den foerste ordren den finner. Kan gjoeres om til aa returnere flere verdier
func CheckAllFloors() int {
	currentFloor := driver.ElevGetFloorSensorSignal()
	for floor := range localTable {
		if floor != currentFloor {
			for i := 0; i < len(localTable[floor]); i++ {
				if localTable[floor][i] == 1 {
					return floor
				}
			}
		}
	}
	return -1
}

func FindDirection() int {
	var diff int
	currentFloor := driver.ElevGetFloorSensorSignal()
	if currentFloor != -1 && currentFloor < N_FLOORS {
		diff = currentFloor - CheckAllFloors()
	}
	if diff > 0 {
		return 1
	} else if diff < 0 {
		return 0
	} else {
		return -1
	}
}

func Init() {
	for floor := range localTable {
		for i := 0; i < len(localTable[floor]); i++ {
			localTable[floor][i] = 0
		}
	}
}

func PrintTable(){
	fmt.Println(localTable)
}
/*
type Order struct {
	Floor, Dir, Cart int
}

func NewOrder(f int, d int, c int) Order {
	return Order{f, d, c}
}

type OrderTable [4][2]Order

func ClaimOrder(o Order, t *OrderTable) {
	t[o.Floor][o.Dir].Cart = CART_ID
}

func InsertOrder(o Order, t *OrderTable) {
	t[o.Floor][o.Dir] = o
}

*/
