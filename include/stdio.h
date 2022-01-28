#pragma once

typedef struct FILE {} FILE;

extern FILE *stdin;
extern FILE *stdout;
extern FILE *stderr;

int getchar(void);
int putchar(int c);
int putc(int c, FILE *stream);

int printf(const char * restrict format, ...);
int puts(const char *str);

int scanf(const char * restrict format, ...);
