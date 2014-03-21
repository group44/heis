package elevator

import (
	"../driver"
	"../order"
	"../types"
	"fmt"
	"time"
)

const (
	IDLE = 2
	OPEN = 3
	UP   = 0
	DOWN = 1

	ON  = 1
	OFF = 0
)

var (

	// Global channels
	state               int
	nextstate           int
	SafetyFloorCh       = make(chan bool)
	ElevatorDirectionCh = make(chan string)

	// Local channels
	doorTimerCh = make(chan bool)
	Temp        bool
)

func Run() {

	// Initialization
	state = IDLE
	nextstate = IDLE

	driver.ElevInit()
	driver.ElevInitLights()
	//LightsInit()
	go FloorLights()

	for driver.ElevGetFloorSensorSignal() == -1 {
		driver.ElevSetSpeed(-200)
	}
	driver.ElevSetSpeed(0)

	//Safety go routine
	//go Safety()
	//Door timer go routine
	//Selveste statemaskinen
	go SetDoorTimer()

	ControlStateMachine()
	//order.UpdateLocalTable(order.LocalOrders, order.C1)

}

func ControlStateMachine() {

	for {

		time.Sleep(5 * time.Millisecond)

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
			<-doorTimerCh
			fmt.Println("Timer done")
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
				driver.ElevSetDoorOpenLamp(ON)
				order.ClearOrder()
				//start timer
				doorTimerCh <- true
				break

			case UP:
				driver.ElevSetDoorOpenLamp(OFF)
				driver.ElevSetSpeed(300) // verdi?
				break

			case DOWN:
				driver.ElevSetDoorOpenLamp(OFF)
				driver.ElevSetSpeed(-300) // verdi?
				break

			default:
				break
			}
		}

		state = nextstate
	}
}

func FloorLights() {

	for {
		time.Sleep(10 * time.Millisecond)
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
		<-doorTimerCh
		time.Sleep(10 * time.Millisecond)

		fmt.Println("Timer started")
		time.Sleep(3000 * time.Millisecond)

		doorTimerCh <- true
	}

}
