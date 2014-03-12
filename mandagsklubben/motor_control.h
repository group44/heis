#ifndef __INCLUDE_MOTOR_CONTROL_H__
#define __INCLUDE_MOTOR_CONTROL_H__


void motor_control_start_motor(int speed, int direction);

void motor_control_stop_motor(int direction, int currentFloor);

void motor_control_open_door(void);

void motor_control_close_door(void);



#endif //#ifndef __MOTOR_CONTROL_H__