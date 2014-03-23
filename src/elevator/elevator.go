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
	fmt.Println("The elevator program is turned off")
}

func Idle() {
	for {

		<-idleCh

		//fmt.Println("ENTERED IDLE")

		driver.ElevSetDoorOpenLamp(OFF)
		driver.ElevSetSpeed(0) // Maa haandtere braastopp-tingen

		//fmt.Println(order.CheckCurrentFloor())
		//fmt.Println(order.FindDirection())
		for { //Her går den helt til den oppnår betingelsene for en ny state
			if order.CheckCurrentFloor() {
				openCh <- true
				break

			} else if order.FindDirection() == 1 { // || order.FindDirection() == -1 {
				downCh <- true
				break
			} else if order.FindDirection() == 0 {
				upCh <- true
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		//fmt.Println("IDLE END")

	}
}

func Open() {
	for {

		<-openCh
		//fmt.Println("ENTERED OPEN")

		driver.ElevSetSpeed(0) // Maa haandtere braastopp-tingen
		order.ClearOrder()

		doorTimerStartCh <- true
		<-doorTimerDoneCh
		//time.Sleep(3000 * time.Millisecond)
		idleCh <- true

	}
}

func Down() {
	for {

		<-downCh
		//fmt.Println("ENTERED DOWN")

		driver.ElevSetSpeed(-SPEED) // verdi?
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
		//fmt.Println("ENTERED UP")
		elevatorDirection = UP
		driver.ElevSetSpeed(SPEED) // verdi?

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
		time.Sleep(100 * time.Millisecond)
		driver.ElevSetFloorIndicator(driver.ElevGetFloorSensorSignal())
	}
}

// Too make sure the elevator never drives over top floor or under bottom floor. Returns true if elevator reaches top floor or bottom floor without order there.
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
		fmt.Println("Timer started")
		time.Sleep(3000 * time.Millisecond)
		fmt.Println("Timer done")
		driver.ElevSetDoorOpenLamp(OFF)

		doorTimerDoneCh <- true
	}

}
