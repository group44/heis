package order

import (
    "../types"
	"../driver"
	"fmt"
	"os"
	//"time"
	"../com"
)

const (
    UP = 0
	DOWN = 1
	INTERNAL = 2
)


var (

	// Channel for signaling type of lights to be set, buffer = 2
	UpdateLightCh = make(chan string, 2)

	GlobalOrders types.GlobalTable
	InternalOrders types.InternalTable
	
)


func Run() {
	
	GlobalOrders = types.NewGlobalTable()
	InternalOrders = types.NewInternalTable()

	go UpdateInternalTable()
	go UpdateLights()
	go CalculateCost()

}
   

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

/* Todo
func CheckError(err string) {
	if err != nil {
		fmt.Println("")
	}
}
*/

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne, denne i en go routine?
func UpdateInternalTable() {
	for {
		for i := range InternalOrders {
			if InternalOrders[i] != 1 {
				// mu?
				InternalOrders[i] = driver.ElevGetButtonSignal(types.N_BUTTONS-1, i)
				// mu?
				UpdateLightCh <- "internal"
			}
		}
	}
}

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne
func ClearOrder() {
	floor := driver.ElevGetFloorSensorSignal() 
	dir := driver.ElevGetDirection()
	
	if floor < 1 || floor > types.N_FLOORS {
		// Assert here
		fmt.Println("Invalid floor number")
	}
	if dir != 0 && dir != 1 {
		// Assert here
		fmt.Println("Invalid dir number")
	}

	InternalOrders[floor-1], GlobalOrders[dir][floor-1] = 0, 0
	UpdateLightCh <- "internal"
	UpdateLightCh <- "global"
}


// Vurder assert, tar ikke hensyn til retning. Sjekker kun om det er ordre den skal ta selv
func CheckCurrentFloor() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	fmt.Println(currentFloor)
	if currentFloor < 1 || currentFloor > types.N_FLOORS {
		// Assert here
		fmt.Println("Invalid floor number")
	}
	fmt.Println("test1")
	return InternalOrders[currentFloor-1] == 1 || GlobalOrders[currentFloor-1][0] == types.CART_ID || GlobalOrders[currentFloor-1][1] == types.CART_ID
}

// MonInternalOrdersors the external buttons on the carts own panel and sends a Data struct wInternalOrdersh order on
// OutputCh if one is found. Runs in separate go routine.
func CheckExternalButtons() {
	data := types.Data{ Head:"Order" }
	for {

		for i := 0; i < types.N_FLOORS; i++ {
			if driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, i) != 0 {
				data.Order = []int{i, 0}
				com.OutputCh <- data
			} else if driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, i) != 0 {
				data.Order = []int{i, 1}
				com.OutputCh <- data
			}
		}

	}
}

// For enkel, returnerer bare den foerste ordren den finner. Kan gjoeres om til aa returnere flere verdier
func CheckAllFloors() int {
	//Checks internal orders first
	currentFloor := driver.ElevGetFloorSensorSignal()

	for floor := range InternalOrders {
		if floor != currentFloor {
			if InternalOrders[floor] == 1 || GlobalOrders[floor][0] == 1 || GlobalOrders[floor][1] == 1 {
				return floor
			}
		}
	}

	return -1
}


func FindDirection() int {
	var diff int
	currentFloor := driver.ElevGetFloorSensorSignal()
	if currentFloor != -1 && currentFloor < types.N_FLOORS && CheckAllFloors() != -1 {
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


func PrintTable(){
	fmt.Println("Internal:")
	fmt.Println(InternalOrders)

	fmt.Println("Global:")
	fmt.Println(GlobalOrders)
}

// Sends out an order from the global table for a new auction. Called if a peer has disconnected and InternalOrders has
// unfinished orders in the global table.
func Redistribute() {

}


/*
func SetLights(lt types.InternalTable, gt types.GlobalTable){
    for{
        <- UpdateLightCh
        //time.Sleep(10 * time.Millisecond)
        for floor := range lt {
        	for i := 0; i < len(lt [floor]); i++ {
        	    driver.ElevSetLights(floor, i, lt[floor][i])
        	}
        }
    }
}
*/


// In separate goroutine
func UpdateLights() {
	var msg string
	for {
		//time.Sleep(10 * time.Millisecond)
		msg = <- UpdateLightCh

		switch msg {
		case "internal":
			for i := range InternalOrders {
				driver.ElevSetLights(i, 2, InternalOrders[i])
			}

		case "global":
			for j, k := 0, 0; j < types.N_FLOORS && k < 2; j, k = j+1, k+1 {
				if GlobalOrders[j][k] != 0 {
					driver.ElevSetLights(j, k, 1)
				} else {
					driver.ElevSetLights(j, k, 0)
				}
			}
		}

	}
}