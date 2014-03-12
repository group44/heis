
#include "elevator_control.h"


#include <stdio.h>
#include "timer_control.h"
#include "order_control.h"
#include "motor_control.h"
#include "elev.h"


/*typedef_enum {
	IDLE 		= 0,
	OPEN		= 1,
	CLOSE 		= 2,
	DRIVE		= 3,
	EMERGENCY 	= 4,
} states_enum;*/

states_enum state, nextstate;
int currentFloorTest;


void elevator_control_state_machine(void)
{
	switch (state) {

		case IDLE:
			if		(elev_get_obstruction_signal())
						nextstate = IDLE;
			else if		(order_control_order_on_floor())
						nextstate = OPEN;
			else if		(order_control_is_order(currentFloorTest) != -1 && timer_control_stop(3))
						nextstate = DRIVE;
			break;

		case OPEN:
			if		(timer_control_stop(3))
						nextstate = CLOSE;			
			break;

		case CLOSE:
			if		(elev_get_obstruction_signal() || order_control_order_in_direction(currentFloorTest))
						nextstate = OPEN;
			else if		(order_control_is_order(currentFloorTest) != -1)
						nextstate = DRIVE;
			else
						nextstate = IDLE;
			break;

		case DRIVE:
			if		(order_control_order_in_direction(currentFloorTest) == 1)
						nextstate = OPEN;
			else if		(elevator_control_safety_check())
						nextstate = CLOSE;
			break;

		case EMERGENCY:
			if		(order_control_emergency_floor_check(currentFloorTest))
						nextstate = OPEN;
			else if		(order_control_emergency_start() && timer_control_stop(3)) {
						motor_control_close_door();
						nextstate = DRIVE;
						order_control_emergency_direction(currentFloorTest);
			}
			else
						order_control_init_order();
			break;
		default:
			break;
	}


	if (elev_get_stop_signal())
		nextstate = EMERGENCY;
	//if (state == IDLE)
		//printf("\n\n\n\n\n\n %d", state);
	//printf("\n\n\n\n\n\n\n\n\nlaststate: %d, currentFloorTest: %d, state: %d, direction: %d", state, currentFloorTest, nextstate, order_control_get_direction());

	//printf("order control is %d \n\n current floor: %d \n\n \n\n", order_control_is_order(currentFloorTest), currentFloorTest);
	if (state != nextstate) {
		printf("laststate: %d, currentFloorTest: %d, state: %d, direction: %d \n\n", state, currentFloorTest, nextstate, order_control_get_direction());
		switch (nextstate) {

			case IDLE:
				break;

			case OPEN:
				elev_set_stop_lamp(0);
				order_control_clear_order_on_floor();
				motor_control_stop_motor(order_control_get_direction(), currentFloorTest);
				motor_control_open_door();
				timer_control_start();
				break;

			case CLOSE:
				motor_control_close_door();
				break;

			case DRIVE:
				if	(elev_get_obstruction_signal()) {
					nextstate = IDLE;
					break;
				}
				motor_control_close_door();
				elev_set_stop_lamp(0);
				motor_control_start_motor(MOTORSPEED, order_control_get_direction());
				break;

			case EMERGENCY:
				motor_control_stop_motor(order_control_get_direction(), currentFloorTest);
				elev_set_stop_lamp(1);
				order_control_init_order();
				break;
			default:
				break;
			}
	}

	state = nextstate;

}


void elevator_control_init_elevator(void){ // initialiserer heisen.
	
	order_control_init_order();
	elevator_control_floor_light();
	elev_set_stop_lamp(0);
	order_control_set_direction(UP);
	if (elev_get_floor_sensor_signal() == -1){
		currentFloorTest = 0;
		while (elev_get_floor_sensor_signal() == -1){
		elev_set_speed(MOTORSPEED);
		}
	motor_control_stop_motor(order_control_get_direction(), currentFloorTest);
	}
	elev_set_speed(0);
	state = IDLE;
	nextstate = IDLE;
}

void elevator_control_update(void){ // oppdaterer variabler, lys og sjekker for nye ordre.
	order_control_set_order();
	currentFloorTest = elevator_control_current_floor(currentFloorTest);
	order_control_light_update();
	elevator_control_floor_light();
	order_control_lastDirection_update();
}

int elevator_control_safety_check(void){
        if ( (elev_get_floor_sensor_signal() == 0 && currentFloorTest != elev_get_floor_sensor_signal()) ||  (elev_get_floor_sensor_signal() == (N_FLOORS - 1) && currentFloorTest != elev_get_floor_sensor_signal()) ) {
				motor_control_stop_motor(order_control_get_direction(), currentFloorTest);
                return 1;
        }
	return 0;
}


void elevator_control_floor_light(void){
	int floor = elev_get_floor_sensor_signal();
	if (N_FLOORS > elev_get_floor_sensor_signal() && floor >= 0)
		elev_set_floor_indicator(floor);
}

int elevator_control_current_floor(int current){
	if (elev_get_floor_sensor_signal() >= 0)
		return elev_get_floor_sensor_signal();
	else
		return current;
}





