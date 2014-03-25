package order

import (
	"../com"
	"../driver"
	"../types"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"os/signal"
)

const (
	UP       = 0
	DOWN     = 1
	INTERNAL = 2
)

var (
	UpdateLightCh  = make(chan string, 5)
	Direction      int
	GlobalOrders   types.GlobalTable
	InternalOrders types.InternalTable
	osChan chan os.Signal
)

func Run() {

	done := make(chan bool)

	GlobalOrders = types.NewGlobalTable()
	InternalOrders = types.NewInternalTable()
	Direction = UP

	go UpdateInternalTable()
	ReadFile()
	go Auction(GlobalOrders)
	go AddOrder()
	go UpdateLights()
	go HandleCost()
	go RemoveOrder()
	go PrintTables()
	go PrintOrderDirection()
	go CheckExternalButtons()
	go Redistribute()

	<-done
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func UpdateInternalTable() {
	for {
		time.Sleep(10 * time.Millisecond)
		for i := range InternalOrders {
			if InternalOrders[i] != 1 {
				if driver.ElevGetButtonSignal(INTERNAL, i) == 1 {
					InternalOrders[i] = driver.ElevGetButtonSignal(INTERNAL, i)
					UpdateLightCh <- "internal"
				}
			}
		}
	}
}

func ClearOrder() {
	floor := GetCurrentFloor()
	dir := GetOrderDirection()
	data := types.Data{Head: "removeorder"}

	if floor < 0 || floor >= types.N_FLOORS {
		fmt.Println("Invalid floor number")
	}
	if dir != 0 && dir != 1 {
		fmt.Println("Invalid dir number")
	}

	InternalOrders[floor] = 0
	data.Order = []int{floor, dir}
	com.OutputCh <- data
	UpdateLightCh <- "internal"
	UpdateLightCh <- "global"
}

// Checks the floor the elevator is on for orders
func CheckCurrentFloor() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	dir := GetOrderDirection()
	var oppDir int
	if dir == UP {
		oppDir = DOWN
	} else if dir == DOWN {
		oppDir = UP
	}

	if currentFloor == 0 {
		ChangeOrderDirection(UP)
	} else if currentFloor == (types.N_FLOORS - 1) {
		ChangeOrderDirection(DOWN)
	}

	if currentFloor == -1 {
		return false
	}

	if currentFloor != -1 && currentFloor < types.N_FLOORS {
		if InternalOrders[currentFloor] == 1 {
			return true
		} else if GlobalOrders[currentFloor][dir] == types.CART_ID {
			return true
		} else if CheckCurrentFloorHelper() == false && GlobalOrders[currentFloor][oppDir] == types.CART_ID {
			ChangeOrderDirection(oppDir)
			return true
		}
	}
	return false
}

func CheckCurrentFloorHelper() bool {
	currentFloor := GetCurrentFloor()
	dir := GetOrderDirection()

	if currentFloor != -1 {
		switch dir {
		case UP:
			for i := currentFloor; i < types.N_FLOORS; i++ {
				if GlobalOrders[i][UP] == types.CART_ID {
					return true
				}
			}
			for j := types.N_FLOORS - 1; j > currentFloor; j-- {
				if GlobalOrders[j][DOWN] == types.CART_ID {
					return true
				}
			}
		case DOWN:
			for k := currentFloor; k >= 0; k-- {
				if GlobalOrders[k][DOWN] == types.CART_ID {
					return true
				}
			}
			for l := 0; l < currentFloor; l++ {
				if GlobalOrders[l][UP] == types.CART_ID {
					return true
				}
			}
		}
	}
	return false
}

func CheckExternalButtons() {
	data := types.Data{Head: "order"}

	for {
		time.Sleep(1 * time.Millisecond)
		for i := 0; i < types.N_FLOORS; i++ {
			time.Sleep(1 * time.Millisecond)
			if driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, i) != 0 {
				data.Order = []int{i, 0}
				com.OutputCh <- data
				fmt.Println("Order created:", data.Order)
				fmt.Println("")
				for driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, i) == 1 {
					time.Sleep(1 * time.Millisecond)
				}
			} else if driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, i) != 0 {
				data.Order = []int{i, 1}
				com.OutputCh <- data
				fmt.Println("Order created:", data.Order)
				fmt.Println("")
				for driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, i) == 1 {
					time.Sleep(1 * time.Millisecond)
				}
			}
		}
	}

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

func CheckOtherFloors() int {
	currentFloor := GetCurrentFloor()
	dir := GetOrderDirection()
	if currentFloor == 0 {
		ChangeOrderDirection(UP)
	} else if currentFloor == (types.N_FLOORS - 1) {
		ChangeOrderDirection(DOWN)
	}

	switch dir {
	case UP:
		for floor := currentFloor; floor < types.N_FLOORS; floor++ {
			if floor != currentFloor {
				if InternalOrders[floor] == 1 || GlobalOrders[floor][UP] == types.CART_ID {
					return floor
				}
			}
		}
		for floor := currentFloor; floor < types.N_FLOORS; floor++ {
			if floor != currentFloor {
				if GlobalOrders[floor][DOWN] == types.CART_ID {
					return floor
				}
			}
		}
		for floor := currentFloor; floor >= 0; floor-- {
			if floor != currentFloor {
				if InternalOrders[floor] == 1 {
					ChangeOrderDirection(DOWN)
					return floor
				}
			}
		}
		ChangeOrderDirection(DOWN)

	case DOWN:
		for floor := currentFloor; floor >= 0; floor-- {
			if floor != currentFloor {
				if InternalOrders[floor] == 1 || GlobalOrders[floor][DOWN] == types.CART_ID {
					return floor
				}
			}
		}
		for floor := currentFloor; floor >= 0; floor-- {
			if floor != currentFloor {
				if GlobalOrders[floor][UP] == types.CART_ID {
					return floor
				}
			}
		}
		for floor := currentFloor; floor < types.N_FLOORS; floor++ {
			if floor != currentFloor {
				if InternalOrders[floor] == 1 {
					ChangeOrderDirection(UP)
					return floor
				}
			}
		}
		ChangeOrderDirection(UP)
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
	for {
		time.Sleep(3000 * time.Millisecond)
		dir := GetOrderDirection()
		switch dir {
		case UP:
			fmt.Println("Order direction: UP")
		case DOWN:
			fmt.Println("Order direction: DOWN")
		}
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

func PrintTables() {
	for {
		time.Sleep(3000 * time.Millisecond)
		fmt.Println("Internal:", InternalOrders)
		fmt.Println("Global:", GlobalOrders)
	}
}

func Redistribute() {
	var data types.Data
	var present bool
	var temp int
	for {
		time.Sleep(1000 * time.Millisecond)
		for floor := range GlobalOrders {
			for i := 0; i < 2; i++ {
				temp = GlobalOrders[floor][i]
				if temp != 0 {

					_, present = com.PeerMap.M[temp]
					if GlobalOrders[floor][i] != 0 && !present {
						fmt.Println("HELVETES HORE FAEN")
						data = types.Data{Head: "removeorder"}
						data.Order = []int{floor, i}
						com.OutputCh <- data
						time.Sleep(20 * time.Millisecond)
						data = types.Data{Head: "order"}
						data.Order = []int{floor, i}
						com.OutputCh <- data
					}
				}
			}
		}
	}
}

func UpdateLights() {
	var msg string
	for {
		time.Sleep(1 * time.Millisecond)
		msg = <-UpdateLightCh

		switch msg {
		case "internal":
			for i := range InternalOrders {
				time.Sleep(1 * time.Millisecond)
				driver.ElevSetLights(i, 2, InternalOrders[i])
			}
		case "global":

			for j := 0; j < types.N_FLOORS; j++ {
				for k := 0; k < 2; k++ {
					time.Sleep(1 * time.Millisecond)
					if GlobalOrders[j][k] != 0 {
						driver.ElevSetLights(j, k, 1)
					} else {
						driver.ElevSetLights(j, k, 0)
					}
				}
			}
		}
	}
}

func  OsTest() {  
        osChan := make(chan os.Signal, 1)                                                      
    signal.Notify(osChan, os.Interrupt)
    <- osChan    
    WriteFile()
    fmt.Println("Programmet er blitt avsluttet")
    time.Sleep(100*time.Millisecond)
    //stop elevator her...
    os.Exit(1)
}

func ReadFile() {
	b, err := ioutil.ReadFile("backup.txt")
	if err != nil {
		panic(err)
	}
	internal := strings.Split(string(b), "")
	for i := 0; i < types.N_FLOORS; i++ {
		InternalOrders[i], _ = strconv.Atoi(internal[i])
	}
	UpdateLightCh <- "internal"
}

func WriteFile() {
	fmt.Println("Backup skrevet til fil")
	msg := strconv.Itoa(InternalOrders[0]) + strconv.Itoa(InternalOrders[1]) + strconv.Itoa(InternalOrders[2]) + strconv.Itoa(InternalOrders[3])
	buf := []byte(msg)
	_ = ioutil.WriteFile("backup.txt", buf, 0644)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func AddOrder() {
	var inc types.Data
	var winnerId int
	order := make([]int, 2)
	for {
		inc = <-com.AddOrderCh
		order = inc.Order
		winnerId = inc.WinnerId
		GlobalOrders[order[0]][order[1]] = winnerId
		com.OutputCh <- types.Data{Head: "table", Order: order, Table: GlobalOrders}
		fmt.Println("new global table:")
		fmt.Println(GlobalOrders)
		UpdateLightCh <- "global"
		time.Sleep(25 * time.Millisecond)
	}
}

func RemoveOrder() {
	order := make([]int, 2)
	for {
		order = <-com.RemoveOrderCh
		GlobalOrders[order[0]][order[1]] = 0
		com.OutputCh <- types.Data{Head: "table", Order: order, Table: GlobalOrders}
		fmt.Println("new global table removed:")
		fmt.Println(GlobalOrders)
		UpdateLightCh <- "global"
		time.Sleep(25 * time.Millisecond)
	}
}
