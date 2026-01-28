#ifndef READWRITE_H
#define READWRITE_H

#include <stdio.h>
#include <stdlib.h>

int readFile(FILE* f, size_t* size, char** content);
int writeFile(FILE* f, char* content);

#endif