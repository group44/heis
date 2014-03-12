
#include <stdio.h>
#include "motor_control.h"
#include "elev.h"
#include "timer_control.h"



int motorDirection;


void motor_control_start_motor(int speed, int direction){ // tar inn speed og direction og setter igang motoren. setter også motorDirection.
	if (direction == UP){
		elev_set_speed(speed);
	}
	else if (direction == DOWN){
		elev_set_speed(-speed);
	}
	else
		elev_set_speed(0);
}

void motor_control_stop_motor(int direction, int currentFloor){ // Stopper heisen ved å kjøre den litt i motsatt retning
	if (elev_get_obstruction_signal() == 0 && currentFloor != elev_get_floor_sensor_signal()){
		if (direction == UP){
			elev_set_speed(-300);
		}
		if (direction == DOWN){
			elev_set_speed(300);
		}
	timer_control_sleep(10);
	elev_set_speed(0);
	}
	elev_set_speed(0);
}

void motor_control_open_door(void){ // setter lyset som representerer åpen dør
	elev_set_door_open_lamp(1);
	timer_control_start();
}

void motor_control_close_door(void){ // slukker lyset som representerer åpen dør
	if (elev_get_obstruction_signal() == 0)
		elev_set_door_open_lamp(0);
}



