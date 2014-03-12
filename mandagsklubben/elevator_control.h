#ifndef __INCLUDE_ELEVATOR_CONTROL_H__
#define __INCLUDE_ELEVATOR_CONTROL_H__

typedef enum {
	IDLE 		= 0,
	OPEN		= 1,
	CLOSE 		= 2,
	DRIVE		= 3,
	EMERGENCY 	= 4,
} states_enum;

void elevator_control_state_machine(void);

void elevator_control_init_elevator(void);

void elevator_control_update(void);

int elevator_control_safety_check(void);

void elevator_control_floor_light(void);

int elevator_control_current_floor(int current);



#endif //#ifndef __ELEVATOR_CONTROL_H__
