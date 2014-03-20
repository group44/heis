package order

import (
	"../com"
	"../driver"
	"../types"
	"fmt"
	"os"
	"time"
)

const (
	UP       = 0
	DOWN     = 1
	INTERNAL = 2
)

var (

	// Channel for signaling type of lights to be set, buffer = 2
	UpdateLightCh = make(chan string, 2)
	Direction     int

	GlobalOrders   types.GlobalTable
	InternalOrders types.InternalTable
)

func Run() {

	GlobalOrders = types.NewGlobalTable()
	InternalOrders = types.NewInternalTable()
	Direction = UP

	go UpdateInternalTable()
	go UpdateLights()
	go CheckExternalButtons()
	//go CalculateCost()
	go Auction(GlobalOrders)

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

/*
func GetClosestElevator() {

}
*/

// INTERNAL maa erstattes, vurder assert
// Vurder navn paa denne, denne i en go routine?
func UpdateInternalTable() {
	for {
		time.Sleep(10 * time.Millisecond)
		for i := range InternalOrders {
			if InternalOrders[i] != 1 {
				if driver.ElevGetButtonSignal(INTERNAL, i) == 1 {
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
	dir := GetOrderDirection()

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
		return InternalOrders[currentFloor] == 1 || GlobalOrders[currentFloor][UP] == types.CART_ID || GlobalOrders[currentFloor][DOWN] == types.CART_ID
	}
	return false
}

// MonInternalOrdersors the external buttons on the carts own panel and sends a Data struct wInternalOrdersh order on
// OutputCh if one is found. Runs in separate go routine.
func CheckExternalButtons() {
	data := types.Data{Head: "Order"}
	for {
		time.Sleep(10 * time.Millisecond)
		for i := 0; i < types.N_FLOORS; i++ {
			if driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, i) != 0 {
				data.Order = []int{i, 0}
				com.OutputCh <- data
				for driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, i) == 1 {
					//time.Sleep(50*time.Millisecond)
				}
			} else if driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, i) != 0 {
				data.Order = []int{i, 1}
				com.OutputCh <- data
				for driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, i) == 1 {
					//time.Sleep(50*time.Millisecond)
				}
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

func GetCurrentFloor() int {
	currentFloor := -1
	x, y := driver.IoReadBit(driver.FLOOR_IND1), driver.IoReadBit(driver.FLOOR_IND2)
	switch x {
	case 0:
		switch y {
		case 0:
			currentFloor = 0
		case 1:
			currentFloor = 1
		}
	case 1:
		switch y {
		case 0:
			currentFloor = 2
		case 1:
			currentFloor = 3
		}
	}

	return currentFloor
}

//skal bli forbedret funksjon som sjekker de andre etasjene og returnerer den nermeste ordren
func CheckOtherFloors() int {
	currentFloor := GetCurrentFloor()
	dir := GetOrderDirection()
	switch dir {
	case UP:
		for floor := currentFloor; floor < types.N_FLOORS; floor++ {
			if floor != currentFloor {
				if InternalOrders[floor] == 1 || GlobalOrders[floor][UP] == types.CART_ID {
					return floor
				}
			}
		}
		for floor := currentFloor; floor >= 0; floor-- {
			if floor != currentFloor {
				if GlobalOrders[floor][UP] == types.CART_ID {
					return floor
				} else if InternalOrders[floor] == 1 {
					ChangeOrderDirection(DOWN)
					return floor
				}
			}
		}
		ChangeOrderDirection(DOWN)
		//fmt.Println("Order direction changed too DOWN")

	case DOWN:
		for floor := currentFloor; floor >= 0; floor-- {
			if floor != currentFloor {
				if InternalOrders[floor] == 1 || GlobalOrders[floor][DOWN] == types.CART_ID {
					return floor
				}
			}
		}

		for floor := currentFloor; floor < types.N_FLOORS; floor++ {
			if floor != currentFloor {
				if GlobalOrders[floor][DOWN] == types.CART_ID {
					return floor
				} else if InternalOrders[floor] == 1 {
					ChangeOrderDirection(UP)
					return floor
				}
			}
		}
		ChangeOrderDirection(UP)
		//fmt.Println("Order direction changed too UP")
	}
	return -1
}

func ChangeOrderDirection(dir int) {
	Direction = dir
}

func GetOrderDirection() int {
	return Direction
}

func PrintOrderDirection() {
	dir := GetOrderDirection()
	switch dir {
	case UP:
		fmt.Println("Order direction: UP")
	case DOWN:
		fmt.Println("Order direction: DOWN")
	}
}

func FindDirection() int {
	var diff int
	currentFloor := driver.ElevGetFloorSensorSignal()
	if currentFloor != -1 && currentFloor < types.N_FLOORS && CheckOtherFloors() != -1 {
		diff = currentFloor - CheckOtherFloors()
	}
	if diff > 0 {
		return DOWN
	} else if diff < 0 {
		return UP
	} else {
		return -1
	}
}

func PrintTable() {
	fmt.Println("Internal:")
	fmt.Println(InternalOrders)

	fmt.Println("Global:")
	fmt.Println(GlobalOrders)
}

// Sends out an order from the global table for a new auction. Called if a peer has disconnected and InternalOrders has
// unfinished orders in the global table.
func Redistribute() {

}

// In separate goroutine
func UpdateLights() {
	var msg string

	for {
		time.Sleep(10 * time.Millisecond)
		msg = <-UpdateLightCh

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

func Backup() {
	hore := make([]byte, types.N_FLOORS)
	for {
		time.Sleep(100 * time.Millisecond)

		WriteFile(hore, "backup.txt")
		time.Sleep(100 * time.Millisecond)
		ReadFile("backup.txt")

	}
}

func ReadFile(filename string) {

	//dat, err := ioutil.ReadFile(filename)

	fmt.Println("reading: " + filename)
	f, err := os.Open(filename)
	b1 := make([]byte, 4)
	hore := make([]byte, 4)
	n1, err := f.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, b1)
	for i := 0; i < 4; i++ {
		hore[i] = b1[i]
	}
	fmt.Println(hore)

}
func WriteFile(hore []byte, filename string) {

	//dat, err := ioutil.ReadFile(filename)

	fmt.Println("write: " + filename)
	fmt.Println(hore)
	f, err := os.Open(filename)
	check(err)

	n1, err := f.Write(hore)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, hore)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
