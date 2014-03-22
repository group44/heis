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

	// Global channels
	SafetyFloorCh       = make(chan bool)
	ElevatorDirectionCh = make(chan string)

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
	go SetDoorTimer()

	for driver.ElevGetFloorSensorSignal() == -1 {
		driver.ElevSetSpeed(-SPEED)
	}
	driver.ElevSetSpeed(0)

	go Idle()
	go Open()
	go Down()
	go Up()

	idleCh <- true

	//Safety go routine
	//go Safety()
	//Door timer go routine
	//Selveste statemaskinen

	//ControlStateMachine()
	//order.UpdateLocalTable(order.LocalOrders, order.C1)

	<-done

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
			time.Sleep(100 * time.Millisecond)
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
		idleCh <- true

	}
}

func Down() {
	for {

		<-downCh
		//fmt.Println("ENTERED DOWN")

		driver.ElevSetSpeed(-SPEED) // verdi?

		for !order.CheckCurrentFloor() {
			time.Sleep(100 * time.Millisecond)
		}

		openCh <- true

	}
}

func Up() {
	for {

		<-upCh
		//fmt.Println("ENTERED UP")

		driver.ElevSetSpeed(SPEED) // verdi?

		for !order.CheckCurrentFloor() {
			time.Sleep(100 * time.Millisecond)
		}

		openCh <- true

	}

}

// Old state machine
/*
func ControlStateMachine() {

	for {

		time.Sleep(100 * time.Millisecond)

		switch state {

		case IDLE:
			if order.CheckCurrentFloor() {
				nextstate = OPEN
				fmt.Println("opentest")
			} else if order.FindDirection() == 1 {
				nextstate = DOWN
			} else if order.FindDirection() == 0 {
				nextstate = UP
			}
			break

		case OPEN:
			// if true
			// timer ferdig
			//in the ghetto timer
			doorTimerCh = <- true
			<-doorTimerCh
			nextstate = IDLE
			break

		case UP:
			if order.CheckCurrentFloor() {
				nextstate = OPEN
			} // else if <-SafetyFloorCh {
			//	nextstate = IDLE
			//}
			break

		case DOWN:
			if order.CheckCurrentFloor() {
				nextstate = OPEN
			} // else if <-SafetyFloorCh {
			//nextstate = IDLE
			//}
			break

		default:
			break

		}

		if state != nextstate {

			fmt.Println(state, nextstate)
			order.PrintTable()
			order.PrintOrderDirection()

			switch nextstate {

			case IDLE:
				driver.ElevSetDoorOpenLamp(OFF)
				driver.ElevSetSpeed(0) // Maa haandtere braastopp-tingen
				break

			case OPEN:
				driver.ElevSetSpeed(0) // Maa haandtere braastopp-tingen

				order.ClearOrder()
				break

			case UP:
				//driver.ElevSetDoorOpenLamp(OFF)
				driver.ElevSetSpeed(300) // verdi?
				break

			case DOWN:
				//driver.ElevSetDoorOpenLamp(OFF)
				driver.ElevSetSpeed(-300) // verdi?
				break

			default:
				break
			}
		}

		state = nextstate
	}
}
*/

func FloorLights() {

	for {
		time.Sleep(100 * time.Millisecond)
		driver.ElevSetFloorIndicator(driver.ElevGetFloorSensorSignal())
	}

}

// Kjores i go routine, kan endre channel til string og legge til flere safety ting som nodstopp og obs her lett
func Safety() {

	for {
		if driver.ElevGetFloorSensorSignal() == 0 && !order.CheckCurrentFloor() {
			//skriv til channel
			SafetyFloorCh <- true
		} else if driver.ElevGetFloorSensorSignal() == (types.N_FLOORS-1) && !(order.CheckCurrentFloor()) {
			//skriv til channel
			SafetyFloorCh <- true
		} else {
			SafetyFloorCh <- false
		}
	}
}

func SetDoorTimer() {

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
