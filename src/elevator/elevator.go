package elevator

import (
	"../driver"
	"../order"
	"../types"
	"fmt"
	"time"
)

const (
	UP   = 0
	DOWN = 1
	IDLE = 2
	OPEN = 3

	ON    = 1
	OFF   = 0
	SPEED = 300
)

var (

	// Global channels sjekke om alle er i bruk
	elevatorDirection int

	// Local channels
	doorTimerStartCh = make(chan bool)
	doorTimerDoneCh  = make(chan bool)
	idleCh           = make(chan bool)
	openCh           = make(chan bool)
	downCh           = make(chan bool)
	upCh             = make(chan bool)
)

// Starts the go routines and initializing of the elevator
func Run() {

	done := make(chan bool)
	// Initialization
	driver.ElevInit()
	go FloorLights()
	go DoorTimer()
	elevatorDirection = DOWN
	for driver.ElevGetFloorSensorSignal() == -1 {
		driver.ElevSetSpeed(-SPEED)
	}
	driver.ElevSetSpeed(0)
	go Idle()
	go Open()
	go Down()
	go Up()
	idleCh <- true
	<-done
}

func Idle() {
	for {
		<-idleCh
		fmt.Println("Idle state entered")
		driver.ElevSetDoorOpenLamp(OFF)
		driver.ElevSetSpeed(0)
		for {
			if order.CheckCurrentFloor() {
				openCh <- true
				break
			} else if order.FindDirection() == 1 {
				downCh <- true
				break
			} else if order.FindDirection() == 0 {
				upCh <- true
				break
			}
		}
	}
}

func Open() {
	for {
		<-openCh
		fmt.Println("Open state entered")
		driver.ElevSetSpeed(0)
		order.ClearOrder()
		doorTimerStartCh <- true
		<-doorTimerDoneCh
		idleCh <- true
	}
}

func Down() {
	for {
		<-downCh
		fmt.Println("Down state entered")
		driver.ElevSetSpeed(-SPEED)
		elevatorDirection = DOWN
		for {
			if order.CheckCurrentFloor() {
				openCh <- true
				break
			} else if Safety() {
				idleCh <- true
				break
			}
		}
	}
}

func Up() {
	for {
		<-upCh
		fmt.Println("Up state entered")
		elevatorDirection = UP
		driver.ElevSetSpeed(SPEED)
		for {
			if order.CheckCurrentFloor() {
				openCh <- true
				break
			} else if Safety() {
				idleCh <- true
				break
			}
		}
	}
}

func FloorLights() {
	for {
		time.Sleep(1 * time.Millisecond)
		driver.ElevSetFloorIndicator(driver.ElevGetFloorSensorSignal())
	}
}

func Safety() bool {
	if driver.ElevGetFloorSensorSignal() == 0 && !order.CheckCurrentFloor() && elevatorDirection == DOWN {
		return true
	} else if driver.ElevGetFloorSensorSignal() == (types.N_FLOORS-1) && !(order.CheckCurrentFloor()) && elevatorDirection == UP {
		return true
	}
	return false
}

func DoorTimer() {
	for {
		<-doorTimerStartCh
		driver.ElevSetDoorOpenLamp(ON)
		time.Sleep(3000 * time.Millisecond)
		driver.ElevSetDoorOpenLamp(OFF)
		doorTimerDoneCh <- true
	}
}
