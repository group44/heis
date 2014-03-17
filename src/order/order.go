package order

import (
    "../types"
	"../driver"
	"fmt"
	"os"
	//"time"
)

const (
    UP = 0
	DOWN = 1
	INTERNAL = 2
)


var (
    localTable types.OrderTable
    globalTable types.OrderTable
)
    

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne
func UpdateLocalTable() {
	for i := 0; i < types.N_FLOORS; i++ {
		localTable[i][INTERNAL] = driver.ElevGetButtonSignal(INTERNAL, i)
	}
}

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne
func RemoveOrder() {
	floor := driver.ElevGetFloorSensorSignal() 
	dir := driver.ElevGetDirection()
	if floor != -1 && floor < types.N_FLOORS {
		localTable[floor][INTERNAL] = 0
		localTable[floor][dir] = 0
	}
}

// Vurder assert, tar ikke hensyn til retning
func CheckCurrentFloor() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	if currentFloor != -1 && currentFloor < types.N_FLOORS {
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
	if currentFloor != -1 && currentFloor < types.N_FLOORS {
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


