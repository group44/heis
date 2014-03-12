
#include <stdio.h>
#include "order_control.h"
#include "elev.h"


int order_table[3][N_FLOORS];

int direction, lastDirection;

void order_control_init_order(void){  // resetter ordretabellen
	int i, j;
	for (i = 0; i < 3; ++i){
		for (j = 0; j < N_FLOORS; ++j){
			order_table[i][j] = 0;
		}
	}
}

void order_control_set_order(void){ // setter ordre i ordretabellen
	int i, j;
	for (i = 0; i < 3; ++i){
		for (j = 0; j < N_FLOORS; ++j){
			if (!(i == 0 && j == (N_FLOORS - 1)) && !(i == 1 && j == 0) && elev_get_button_signal(i,j))
				order_table[i][j] = 1;
		}
	}
}

void order_control_clear_order_on_floor(void){ // resetter ordre i etasje og retning som heis er i
	if(elev_get_floor_sensor_signal() != -1){
		order_table[direction][elev_get_floor_sensor_signal()] = 0;
		order_table[2][elev_get_floor_sensor_signal()] = 0;
		if (elev_get_floor_sensor_signal() == 0)
			order_table[0][0] = 0;
		else if (elev_get_floor_sensor_signal() == (N_FLOORS-1))
			order_table[1][(N_FLOORS - 1)] = 0;
	}
}

int order_control_order_on_floor(void){ // Returnerer 1 om det er ordre i etasjen uansett retning, returnerer 0 hvis ikke
	int i;
	for(i = 0; i < 3; ++i){
		if (order_table[i][elev_get_floor_sensor_signal()])
			return 1;
	}
	return 0;
}

int order_control_order_in_direction(int currentFloor){ // Sjekker om det er ordre i retningen heisen kjører i. Både etasjen den er i, og over/under.

	int i;
	if (elev_get_floor_sensor_signal() != -1){
		if (order_table[2][elev_get_floor_sensor_signal()] == 1){
			return 1;
		}
		if (direction == UP){
			if (order_table[UP][elev_get_floor_sensor_signal()] == 1){
				return 1;
			}
			for (i = (currentFloor + 1); i < N_FLOORS; ++i){
				if (order_table[UP][i] == 1){
					return 0;
				}
			}
			
			for (i = (N_FLOORS - 1); i > currentFloor; --i){
				if (order_table[DOWN][i] == 1 || order_table[2][i] == 1){
					return 0;
				}
			}
			
			if (order_table[DOWN][elev_get_floor_sensor_signal()] == 1){
				order_control_change_direction();
				return 1;
			}
		}
		if (direction == DOWN){
			if (order_table[DOWN][elev_get_floor_sensor_signal()] == 1){
				return 1;
			}
			for (i = (currentFloor - 1); i > 0; --i){
				if (order_table[DOWN][i] == 1)
					return 0;
			}
			
			for (i = 0; i < currentFloor; ++i){
				if (order_table[UP][i] == 1 || order_table[2][i] == 1){
					return 0;
				}
			}
			
			if (order_table[UP][elev_get_floor_sensor_signal()] == 1){
				order_control_change_direction();
				return 1;
			}
		}		
	}
	return 0;
}

/*
int order_control_order_in_direction(int currentFloor) {
int i;
if (elev_get_floor_sensor_signal() != (-1)){
if (order_table[2][elev_get_floor_sensor_signal()] == 1){
return 1;
}

if (direction == UP){
if (order_table[UP][elev_get_floor_sensor_signal()] == 1){
return 1;
}
for (i = (currentFloor + 1); i < N_FLOORS; ++i){
if (order_table[UP][i] == 1){
return 0;
}
}
for (i = (N_FLOORS - 1); i > currentFloor; --i){
if (order_table[DOWN][i] == 1 || order_table[2][i] == 1){
return 0;
}
}
if (order_table[DOWN][elev_get_floor_sensor_signal()] == 1){
order_control_change_direction();
return 1;
}

}

if (direction == DOWN){
if (order_table[DOWN][elev_get_floor_sensor_signal()] == 1){
return 1;
}
for (i = (currentFloor - 1); i > 0; --i){
if (order_table[DOWN][i] == 1){
return 0;
}
}
for (i = 0; i < currentFloor; ++i){
if (order_table[UP][i] == 1 || order_table[2][i] == 1){
return 0;
}
}
if (order_table[UP][elev_get_floor_sensor_signal()] == 1){
order_control_change_direction();
return 1;
}
}
}
return 0;
}

*/


void order_control_change_direction(void){ // Endrer direction variabelen
	if (direction == UP)
		direction = DOWN;
	else if (direction == DOWN)
		direction = UP;
}

int order_control_get_direction(void){ // returnerer retning 0 = opp og 1 = ned
	return direction;
}

void order_control_emergency_direction(int currentFloor){ // Bestemmer retning på heis etter nødstopp
	int i, j;
	for (i = 0; i < 3; ++i){
		for (j = 0; j < N_FLOORS; ++j){
			if (order_table[i][j]){
				if (j < currentFloor)
					order_control_set_direction(DOWN);
				else if (j > currentFloor)
					order_control_set_direction(UP);
				else if (lastDirection == UP && j == currentFloor)
					order_control_set_direction(DOWN);
				else if (lastDirection == DOWN && j == currentFloor)
					order_control_set_direction(UP);
			}
		}
	}
}

void order_control_lastDirection_update(void){ // oppdaterer lastDirection
	lastDirection = direction;
}

void order_control_set_direction(int dir){ // setter direction..
	direction = dir;
}


int order_control_is_order(int currentFloor){ // Sjekker ordretabell om det er ordre og returnerer slik at heisen kan begynne å kjøre, returnerer -1 om det ikke er ordre
	
	int j;

	if(direction == UP) {
		for(j = currentFloor+1; j < N_FLOORS; ++j){
			if(order_table[UP][j] == 1 || order_table[2][j] == 1){
				return j;
			}
		}
		for(j = N_FLOORS-1; j >= 0; --j){
			if(order_table[DOWN][j] == 1){
				if(j < currentFloor){
					order_control_change_direction();
				}
				return j;
			}
		}
		for(j = 0; j <= currentFloor; ++j){
			if(order_table[UP][j] == 1 || order_table[2][j] == 1){
				order_control_change_direction();
				return j;
			}
		}
	}

	if(direction == DOWN) {
		for(j = currentFloor-1; j > 0; --j){
			if(order_table[DOWN][j] == 1 || order_table[2][j] == 1){
				return 2;
			}
		}
		for(j = 0; j < N_FLOORS; ++j){
			if(order_table[UP][j] == 1){
				if(j > currentFloor){
					order_control_change_direction();
				}
				return 3;
			}
		}
		for(j = N_FLOORS-1 ; j >= currentFloor; --j){
			if(order_table[DOWN][j] == 1 || order_table[2][j] == 1){
				order_control_change_direction();
				return 4;
			}
		}
	}
	return -1;
}

int order_control_emergency_start(void) {
	int j;
	for (j = 0; j < N_FLOORS; ++j) {
		if (order_table[2][j])
			return 1;
	}
	return 0;
}


int order_control_emergency_floor_check(int currentFloor) {
	return (order_table[2][currentFloor] && currentFloor == elev_get_floor_sensor_signal());
}


void order_control_light_update(void){
	int i, j;
	for (i = 0; i < 3; ++i){
		for (j = 0; j < N_FLOORS; ++j){
			if (!(i == 0 && j == (N_FLOORS - 1)) && !(i == 1 && j == 0) && order_table[i][j] == 1)
				elev_set_button_lamp(i, j, 1);
			else if (!(i == 0 && j == (N_FLOORS - 1)) && !(i == 1 && j == 0) && order_table[i][j] == 0)
				elev_set_button_lamp(i, j, 0);
		}
	}
}



