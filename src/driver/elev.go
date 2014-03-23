package driver

import (
	"../types"
	"math"
	"time"
)

const (
	BUTTON_CALL_UP   types.ElevButtonTypeT = 0
	BUTTON_CALL_DOWN types.ElevButtonTypeT = 1
	BUTTON_COMMAND   types.ElevButtonTypeT = 2
)

// Global variable containing direction bit
var (
	lastSpeed int

	lampChannelMatrix = [types.N_FLOORS][types.N_BUTTONS]int{

		{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
		{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
		{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
		{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	}

	buttonChannelMatrix = [types.N_FLOORS][types.N_BUTTONS]int{

		{FLOOR_UP1, FLOOR_DOWN1, FLOOR_COMMAND1},
		{FLOOR_UP2, FLOOR_DOWN2, FLOOR_COMMAND2},
		{FLOOR_UP3, FLOOR_DOWN3, FLOOR_COMMAND3},
		{FLOOR_UP4, FLOOR_DOWN4, FLOOR_COMMAND4},
	}
)

func ElevInit() int {
	// Init hardware
	if IoInit() == 0 {
		return 0
	}

	// Zero all floor button lamps
	for i := 0; i < types.N_FLOORS; i++ {
		if i != 0 {
			ElevSetButtonLamp(BUTTON_CALL_DOWN, i, 0)
		}

		if i != types.N_FLOORS-1 {
			ElevSetButtonLamp(BUTTON_CALL_UP, i, 0)
		}

		ElevSetButtonLamp(BUTTON_COMMAND, i, 0)
	}

	//Clear stop lamp, door open lamp, and set floor indicator to ground floor
	ElevSetStopLamp(0) // sette inn så alle lys blir nullstilt??
	ElevSetDoorOpenLamp(0)
	ElevSetFloorIndicator(0)
	lastSpeed = 0

	for ElevGetFloorSensorSignal() == -1 {
		ElevSetSpeed(-300) //speed??
	}

	ElevSetSpeed(0)
	ElevInitLights()
	// Return success
	return 1
}

// Implement this
/*
func CheckError(err) {

}
*/

func ElevSetSpeed(speed int) {
	// In order to sharply stop the elevator, the direction bit is toggled
	// before setting speed to zero.

	// If to start (speed > 0)
	// If to stop (speed == 0)
	if speed > 0 {
		IoClearBit(MOTORDIR)
	} else if speed < 0 {
		IoSetBit(MOTORDIR)
	} else if lastSpeed < 0 {
		IoClearBit(MOTORDIR)
	} else if lastSpeed > 0 {
		IoSetBit(MOTORDIR)
	}

	lastSpeed = speed
	//Adjust this to get instant stop, higher value if it does't stop, lower value if
	// it drives a little in the other direction before it stops
	time.Sleep(10 * time.Millisecond)
	// Write new setting to motor
	IoWriteAnalog(MOTOR, 2048+4*int(math.Abs(float64(speed))))
}

func ElevGetFloorSensorSignal() int {
	if IoReadBit(SENSOR1) == 1 {
		return 0
	} else if IoReadBit(SENSOR2) == 1 {
		return 1
	} else if IoReadBit(SENSOR3) == 1 {
		return 2
	} else if IoReadBit(SENSOR4) == 1 {
		return 3
	} else {
		return -1
	}
}

// Fiks assert
func ElevGetButtonSignal(button types.ElevButtonTypeT, floor int) int {
	// assert(floor >= 0);
	//assert(floor < N_FLOORS);
	//assert(!(button == BUTTON_CALL_UP && floor == N_FLOORS-1));
	//assert(!(button == BUTTON_CALL_DOWN && floor == 0));
	//assert( button == BUTTON_CALL_UP ||
	//       button == BUTTON_CALL_DOWN ||
	//       button == BUTTON_COMMAND);

	if IoReadBit(buttonChannelMatrix[floor][button]) != 0 {
		return 1
	} else {
		return 0
	}
}

func ElevGetStopSignal() int {
	return IoReadBit(STOP)
}

func ElevGetObstructionSignal() int {
	return IoReadBit(OBSTRUCTION)
}

func ElevSetFloorIndicator(floor int) {
	// assert(floor >= 0);
	// assert(floor < N_FLOORS);

	// Binary encoding. One light must always be on.
	switch floor {
	case 0:
		IoClearBit(FLOOR_IND1)
		IoClearBit(FLOOR_IND2)
	case 1:
		IoClearBit(FLOOR_IND1)
		IoSetBit(FLOOR_IND2)
	case 2:
		IoSetBit(FLOOR_IND1)
		IoClearBit(FLOOR_IND2)
	case 3:
		IoSetBit(FLOOR_IND1)
		IoSetBit(FLOOR_IND2)
	default:
	}
}

func ElevSetButtonLamp(button types.ElevButtonTypeT, floor, value int) {
	/*
			assert(floor >= 0);
		    assert(floor < N_FLOORS);
		    assert(!(button == BUTTON_CALL_UP && floor == N_FLOORS-1));
		    assert(!(button == BUTTON_CALL_DOWN && floor == 0));
		    assert( button == BUTTON_CALL_UP ||
		            button == BUTTON_CALL_DOWN ||
		            button == BUTTON_COMMAND);
	*/

	if value == 1 {
		IoSetBit(lampChannelMatrix[floor][button])
	} else {
		IoClearBit(lampChannelMatrix[floor][button])
	}
}

func ElevSetStopLamp(value int) {
	if value == 1 {
		IoSetBit(LIGHT_STOP)
	} else {
		IoClearBit(LIGHT_STOP)
	}
}

func ElevSetDoorOpenLamp(value int) {
	if value == 1 {
		IoSetBit(DOOR_OPEN)
	} else {
		IoClearBit(DOOR_OPEN)
	}
}

func ElevGetDirection() int {
	return IoReadBit(MOTORDIR)
}

func ElevInitLights() {
	for i := 0; i < types.N_FLOORS; i++ {
		for j := 0; j < 3; j++ {
			IoClearBit(lampChannelMatrix[i][j])
		}
	}
}

func ElevSetLights(floor int, button int, temp int) {
	if temp == 1 {
		IoSetBit(lampChannelMatrix[floor][button])
	} else if temp == 0 {
		IoClearBit(lampChannelMatrix[floor][button])
	}
}

func ElevCheckLight(floor int, button int) int {
	return IoReadBit(lampChannelMatrix[floor][button])
}
