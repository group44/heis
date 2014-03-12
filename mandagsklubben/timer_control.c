

#include <stdio.h>
#include "time.h"
#include "timer_control.h"

time_t timestart, timestop;




void timer_control_start(void){ // starter referanseklokke
	timestart = time(NULL);	
}


int timer_control_stop(int sec) { // sjekker forskjellen mellom referanse og en ny klokka
	timestop = time(NULL);
	int difference = (int)difftime(timestop, timestart);
	if (difference < sec) {
		return 0;
	}
	else {
		return 1;
	}
}

void timer_control_sleep(int milisec) {
	struct timespec req = {0};
	req.tv_sec = 0;
	req.tv_nsec = milisec * 1000000L;
	nanosleep(&req, (struct timespec *)NULL);
}
