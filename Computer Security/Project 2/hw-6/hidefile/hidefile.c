#define _GNU_SOURCE
#include <stdio.h>
#include <stdlib.h>
#include <dirent.h>
#include <dlfcn.h>
#include <string.h>

// Name: Ajay Vejendla  
// netID: av591
// RUID: 178007104


static struct dirent *(*realReaddir)(DIR *dirp) = NULL;

struct dirent *readdir(DIR *dirp){
    //Load enviroment variable
    //Filenames are case sensitive
    const char* hidden = getenv("HIDDEN");
    
    //Load actual readdir function
    realReaddir = dlsym(RTLD_NEXT,"readdir");

    //Load file being requested 
    struct dirent * targetFile = NULL;
    targetFile = realReaddir(dirp);

    //Load filename
    char * fileName = NULL;
    if (targetFile != NULL){
        fileName = targetFile->d_name;
    }
    

    //If both the file requested and the enviroment variable exist
    //Compare filename and HIDDEN
    //If equal, return NULL dirent
    if((targetFile != NULL) && (hidden != NULL) && (fileName != NULL)){
        int compare = strcmp(hidden,fileName);
        if(compare == 0){
            targetFile = realReaddir(dirp);
        }
    }
   return targetFile;
}