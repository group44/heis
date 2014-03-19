package order

import (
    "../types"
	"../driver"
	"fmt"
	"os"
	"time"
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
	go Auction()


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
		time.Sleep(10*time.Millisecond)
		for i := range InternalOrders {
			if InternalOrders[i] != 1 {
				if driver.ElevGetButtonSignal(INTERNAL, i) == 1{
					InternalOrders[i] = driver.ElevGetButtonSignal(INTERNAL, i)
					fmt.Println("Internal order table updated")
					UpdateLightCh <- "internal"
				}
			}
		}
	}
}

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne
func ClearOrder() {
	floor := driver.ElevGetFloorSensorSignal() 
	dir := driver.ElevGetDirection()
	
	if floor < 0 || floor >= types.N_FLOORS {
		// Assert here
		fmt.Println("Invalid floor number")
	}
	if dir != 0 && dir != 1 {
		// Assert here
		fmt.Println("Invalid dir number")
	}

	InternalOrders[floor], GlobalOrders[dir][floor] = 0, 0
	UpdateLightCh <- "internal"
	UpdateLightCh <- "global"
}


// Vurder assert, tar ikke hensyn til retning. Sjekker kun om det er ordre den skal ta selv
func CheckCurrentFloor() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	if currentFloor < 0 || currentFloor >= types.N_FLOORS && currentFloor != -1 {
		// Assert here
		//fmt.Println("Invalid floor number troroororoo")
	} 
	if currentFloor != -1 {
		return InternalOrders[currentFloor] == 1/* || GlobalOrders[currentFloor][UP] == types.CART_ID || GlobalOrders[currentFloor][DOWN] == types.CART_ID */
	}
	return false
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

//skal bli forbedret funksjon som sjekker de andre etasjene og returnerer den nermeste ordren
func CheckOtherFloors() int{
	currentFloor := -1
	dir := UP // Maa fikses, maa lagre heisens direction ett sted og hente den her
		
	
	if driver.ElevGetFloorSensorSignal() != -1{
		currentFloor = driver.ElevGetFloorSensorSignal()
	} else {
		x, y := driver.IoReadBit(driver.FLOOR_IND1), driver.IoReadBit(driver.FLOOR_IND2)
		switch x {
			case 0:
				switch y{
					case 0:
						currentFloor = 0
					case 1:
						currentFloor = 1
				}
				
			case 1:
				switch y{
					case 0:
						currentFloor = 2
					case 1:
						currentFloor = 3
				}
		}
	}
	
	
	switch dir {
		case UP:
			for floor:=currentFloor; floor<types.N_FLOORS-1;floor++{
				if floor != currentFloor{
					if InternalOrders[floor] == 1 || GlobalOrders[floor][UP] == 1{
						return floor
					}
				}
			}
			for floor:=currentFloor; floor>=0;floor--{
				if floor != currentFloor{
					if InternalOrders[floor] == 1 || GlobalOrders[floor][UP] == 1{
						return floor
					}
				}
			}
			
		case DOWN:
			for floor:=currentFloor; floor<types.N_FLOORS;floor++{
				if floor != currentFloor{
					if InternalOrders[floor] == 1 || GlobalOrders[floor][DOWN] == 1{
						return floor
					}
				}
			}
			for floor:=currentFloor; floor>0;floor--{
				if floor != currentFloor{
					if InternalOrders[floor] == 1 || GlobalOrders[floor][DOWN] == 1{
						return floor
					}
				}
			}
	}
	
	
	



	for floor := range InternalOrders {
		if floor != currentFloor {
		
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
		time.Sleep(10 * time.Millisecond)
		msg = <- UpdateLightCh

		switch msg {
		case "internal":
			fmt.Println("Internal Lights updated")
			for i := range InternalOrders {
				driver.ElevSetLights(i, 2, InternalOrders[i])
			}

		case "global":
			fmt.Println("Global Lights updated")
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
