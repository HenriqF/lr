#include <stdio.h>
#include <windows.h>
#include <sys/time.h>
#include "readwrite.h"

#define F_TIME 0x01

int startsWith(char* a, char* b){
    size_t la = strlen(a);
    size_t lb = strlen(b);

    if (lb > la) return 0;

    return memcmp(a, b, lb) == 0;
}

void runLr(DWORD wait){
    char ppath[MAX_PATH];
    GetCurrentDirectoryA(MAX_PATH, ppath);

    char ls_path[MAX_PATH+3];
    snprintf(ls_path, MAX_PATH+3, "%s\\lr", ppath);

    size_t size;
    char* content;

    FILE* f = fopen(ls_path, "rb");
    if (!f){
        printf("sem arquivo lr");
        return;
    }

    readFile(f, &size, &content);
    fclose(f);

    char command[512];
    int ll = 0;


    for(int i = 0; i <= size; i++){
        if(content[i] == '\n' || i == size){
            snprintf(command, (size_t)i-ll+1, "%s", content+ll);
            system(command);
            ll = i+1;
            if (i != size) Sleep(wait);
        }
    }
}

int main(int argc, char** argv){

    int flags = 0;
    DWORD wait_time = 0;

    if (argc != 1){
        for (int i = 1 ; i < argc; i++){
            if (startsWith(argv[i], "-t")){
                flags |= F_TIME;
                char time_str[10];
                snprintf(time_str, 10, "%s", argv[i]+2);
                wait_time = (DWORD)atoi(time_str);
            }
        }
    }
    
    runLr(wait_time);

    return 0;
}