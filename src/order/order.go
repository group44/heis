package order

import (
	"../com"
	"../driver"
	"../types"
	"fmt"
	"os"
	"time"
	"strings"
    "strconv"
    "io/ioutil"
)

const (
	UP       = 0
	DOWN     = 1
	INTERNAL = 2
)

var (

	// Channel for signaling type of lights to be set, buffer = 2
	UpdateLightCh = make(chan string)
	//UpdateLightCh = make(chan string, 2)
	Direction     int

	GlobalOrders   types.GlobalTable
	InternalOrders types.InternalTable
)

func Run() {

	done := make(chan bool)

	GlobalOrders = types.NewGlobalTable()
	InternalOrders = types.NewInternalTable()
	ReadFile()
	Direction = UP

	go UpdateInternalTable()
	go UpdateLights()
	go CheckExternalButtons()
	go CalculateCost()
	go Auction(GlobalOrders)
	//go UpdateGlobalTable()
	go PrintTables()

	<-done

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

/*
func UpdateGlobalTable() {

	for {
		time.Sleep(10 * time.Millisecond)
		GlobalOrders = <-com.TableCh
		fmt.Println("GlobalTable is updated")
		fmt.Println(GlobalOrders)
	}

}
*/

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

	InternalOrders[floor], GlobalOrders[floor][dir] = 0, 0
	UpdateLightCh <- "internal"
	UpdateLightCh <- "global"
}

// Vurder assert, tar ikke hensyn til retning. Sjekker kun om det er ordre den skal ta selv
/*func CheckCurrentFloor() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	dir := GetOrderDirection()

	if currentFloor != -1 {
		//fmt.Println(GlobalOrders[currentFloor][dir] == types.CART_ID)
		return InternalOrders[currentFloor] == 1 || GlobalOrders[currentFloor][dir] == types.CART_ID //|| GlobalOrders[currentFloor][DOWN] == types.CART_ID
	}
	return false
}*/

// ny funksjon, kan ikke skjønne at denne ikke skal fungere
func CheckCurrentFloor() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	dir := GetOrderDirection()

	if currentFloor == -1 {
		return false
	}
	if currentFloor != -1 && currentFloor < types.N_FLOORS {
		if InternalOrders[currentFloor] == 1 {
			return true
		} else if GlobalOrders[currentFloor][dir] == types.CART_ID {
			return true
		}
	}
	return false
}

// MonInternalOrdersors the external buttons on the carts own panel and sends a Data struct wInternalOrdersh order on
// OutputCh if one is found. Runs in separate go routine.
func CheckExternalButtons() {
	data := types.Data{Head: "order"}

	for {
		time.Sleep(10 * time.Millisecond)
		for i := 0; i < types.N_FLOORS; i++ {
			time.Sleep(10 * time.Millisecond)
			if driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, i) != 0 {
				data.Order = []int{i, 0}
				com.OutputCh <- data

				fmt.Println("Order created:")
				//fmt.Println(data)
				fmt.Println("")

				for driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, i) == 1 {
					time.Sleep(50 * time.Millisecond)
				}
			} else if driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, i) != 0 {
				data.Order = []int{i, 1}
				com.OutputCh <- data

				fmt.Println("Order created:")
				//fmt.Println(data)
				fmt.Println("")

				for driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, i) == 1 {
					time.Sleep(50 * time.Millisecond)
				}
			}
		}
	}

}

// For enkel, returnerer bare den foerste ordren den finner. Kan gjoeres om til aa returnere flere verdier
/*func CheckAllFloors() int {
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
*/

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

//printer tabellen med 2 sek mellomrom
func PrintTables() {
	for {
		fmt.Println("Internal:")
		fmt.Println(InternalOrders)

		fmt.Println("Global:")
		fmt.Println(GlobalOrders)
		time.Sleep(2000*time.Millisecond)
	}
}

// Sends out an order from the global table for a new auction. Called if a peer has disconnected and InternalOrders has
// unfinished orders in the global table.
func Redistribute() {

}

// In separate goroutine
func UpdateLights() {
	var msg string

	for {
		time.Sleep(100 * time.Millisecond)
		msg = <-UpdateLightCh

		switch msg {
		case "internal":
			for i := range InternalOrders {
				time.Sleep(10 * time.Millisecond)
				driver.ElevSetLights(i, 2, InternalOrders[i])
			}
			fmt.Println("Internal Lights updated")

		case "global":

			for j := 0; j < types.N_FLOORS; j++ {
				for k := 0; k < 2; k++ {
					time.Sleep(10 * time.Millisecond)
					if GlobalOrders[j][k] != 0 {
						driver.ElevSetLights(j, k, 1)
					} else {
						driver.ElevSetLights(j, k, 0)
					}
				}
			}
			fmt.Println("Global Lights updated")
		}

	}

}

func Backup() {
	for {
		time.Sleep(200 * time.Millisecond)
		WriteFile()
	}
}

func ReadFile() {
	b, err := ioutil.ReadFile("backup.txt")
    if err != nil { panic(err) }
    internal := strings.Split(string(b),"")
    for i:=0;i<types.N_FLOORS;i++{
    	InternalOrders[i],_ = strconv.Atoi(internal[i])
    }
    //InternalOrders[1],_ = strconv.Atoi(internal[1])
    //InternalOrders[2],_ = strconv.Atoi(internal[2])
    //InternalOrders[3],_ = strconv.Atoi(internal[3])
}

// i denne må det leges til flere hvis heisen utvides til fler etasjer
func WriteFile() {
	msg := strconv.Itoa(InternalOrders[0])+strconv.Itoa(InternalOrders[1])+strconv.Itoa(InternalOrders[2])+strconv.Itoa(InternalOrders[3])
    buf := []byte(msg)
    _ = ioutil.WriteFile("backup.txt", buf, 0644)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
