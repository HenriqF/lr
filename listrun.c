#include <stdio.h>
#include <windows.h>
#include <sys/time.h>
#include "readwrite.h"

int startsWith(char* a, char* b){
    size_t la = strlen(a);
    size_t lb = strlen(b);

    if (lb > la) return 0;

    return memcmp(a, b, lb) == 0;
}

void runLr(DWORD wait){
    char ppath[MAX_PATH];
    GetCurrentDirectoryA(MAX_PATH, ppath);
    char lr_path[MAX_PATH+3];
    snprintf(lr_path, MAX_PATH+3, "%s\\lr", ppath);

    size_t size;
    char* content;

    FILE* f = fopen(lr_path, "rb");
    if (!f){
        printf("sem arquivo lr");
        return;
    }

    readFile(f, &size, &content);
    fclose(f);

    char command[1024];
    int ll = 0;

    for(int i = 0; i <= size; i++){
        if(content[i] == '\n' || i == size){
            snprintf(command, (size_t)i-ll+1, "%s", content+ll);

            if (startsWith(command, "//") == 1 || startsWith(command, "lr")){}
            else if (startsWith(command, "#") == 1){
                printf("Aperte enter para continuar...");
                int c;
                while ((c = getchar()) != '\n' && c != EOF) {}
            }
            else{
                system(command);
                if (i != size) Sleep(wait);
            }

            ll = i+1;
        }
    }
}

int main(int argc, char** argv){
    DWORD wait_time = 0;

    if (argc != 1){
        for (int i = 1 ; i < argc; i++){
            if (startsWith(argv[i], "-t")){
                char time_str[10];
                snprintf(time_str, 10, "%s", argv[i]+2);
                wait_time = (DWORD)atoi(time_str);
            }
        }
    }
    
    runLr(wait_time);
    return 0;
}