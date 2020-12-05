#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <dlfcn.h>

// Name: Ajay Vejendla
// netID: av591
// RUID: 178007104

//Global variable to keep track of how many times time() is called
//Takes advantage of the fact that newtime stays loaded until the calling program is done with it
int numRuns = 0;

//Create function prototype for the actual time() function
static time_t (*realTime)(time_t *tloc) = NULL;

time_t time(time_t *tloc){
    //Load the real time function as realTime
    realTime = dlsym(RTLD_NEXT,"time");

    //If time hasn't been called before in this run, set time to Sept. 1st 2020, I don't remember exactly what time on that day
    if (numRuns == 0){
        *tloc = 1598988279;
    }

    //If time has been called before in this run, call the real time() function
    else{
        *tloc = realTime(NULL);
    }
    
    //Increment counting variable keeping track of the number of time() calls
    numRuns = numRuns + 1;

    return *tloc;
 }