package elevator

import (
	"../types"
	"../order"
	"../driver"
	"fmt"
	"time"
)

const (
	IDLE = 0
	OPEN = 1
	UP = 2
	DOWN = 3
	
	ON = 1
	OFF = 0
)

var (

	state, nextstate int
	SafetyFloorCh = make(chan bool)
	DoorTimerStartCh = make(chan bool) 
	DoorTimerDoneCh = make(chan bool)
	ElevatorDirectionCh = make(chan string)
	Temp bool

)

func Run(){
	
	// Initialization
	state = IDLE
	nextstate = IDLE
	
    driver.ElevInit()
	driver.ElevInitLights()
    //LightsInit()
    go FloorLights()
    
    for (driver.ElevGetFloorSensorSignal() == -1){
        driver.ElevSetSpeed(-200)
    }
    driver.ElevSetSpeed(0)
    
    //Safety go routine
    //go Safety()
    //Door timer go routine
    //Selveste statemaskinen
    go DoorTimer()
    for {
    	time.Sleep(5*time.Millisecond)
        ControlStateMachine()
        //order.UpdateLocalTable(order.LocalOrders, order.C1)
    }
}

func ControlStateMachine() {
	switch state {
	
	case IDLE:
		if order.CheckCurrentFloor() {
			nextstate = OPEN
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
		if <-DoorTimerDoneCh{
			fmt.Println("Timer done")
			nextstate = IDLE
		}
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
			DoorTimerStartCh <- true
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

/*
func LightsInit(){
    //go order.SetLights(order.LocalOrders, order.C1)
}
*/

func FloorLights() {

    for {
        time.Sleep(10*time.Millisecond)
        driver.ElevSetFloorIndicator(driver.ElevGetFloorSensorSignal())
    }

}


// Kjores i go routine, kan endre channel til string og legge til flere safety ting som nodstopp og obs her lett
func Safety() {

	for{
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

func DoorTimer() {

	for {
		time.Sleep(10*time.Millisecond)
		<- DoorTimerStartCh
		fmt.Println("Timer started")
		time.Sleep(3000*time.Millisecond)

		DoorTimerDoneCh <- true
	}

}


















