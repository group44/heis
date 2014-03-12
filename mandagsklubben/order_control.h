#ifndef __INCLUDE_ORDER_CONTROL_H__
#define __INCLUDE_ORDER_CONTROL_H__


void order_control_init_order(void);

void order_control_set_order(void);

void order_control_clear_order_on_floor(void);

int order_control_order_on_floor(void);

int order_control_order_in_direction(int currentFloor);

void order_control_change_direction(void);

int order_control_get_direction(void);

void order_control_emergency_direction(int currentFloor);

void order_control_lastDirection_update(void);

void order_control_set_direction(int dir);

int order_control_is_order(int currentFloor);

int order_control_emergency_start(void);

int order_control_emergency_floor_check(int currentFloor);

void order_control_light_update(void);


#endif //#ifndef __ORDER_CONTROL_H__
