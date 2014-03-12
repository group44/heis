// REPLACE THIS FILE WITH YOUR OWN CODE.
// READ ELEV.H FOR INFORMATION ON HOW TO USE THE ELEVATOR FUNCTIONS.

#include <stdio.h>
#include "elev.h"
#include "elevator_control.h"


int main()
{
	 if (!elev_init()) {
        printf(__FILE__ ": Unable to initialize elevator hardware\n");
        return 1;
    }
    elevator_control_init_elevator();

	while (1) {
		elevator_control_update();
		elevator_control_state_machine();
	}
    return 0;
}

