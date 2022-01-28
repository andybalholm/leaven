#pragma once

#include <stddef.h>

void *memchr(const void *b, int c, size_t len);
int memcmp(const void *b1, const void *b2, size_t len);

char *strchr(const char *s, int c);
int strcmp(const char *s1, const char *s2);
char *strcpy(char * restrict dst, const char * restrict src);
size_t strcspn(const char *s, const char *charset);
char *strncat(char * restrict s, const char * restrict append, size_t count);
int strncmp(const char *s1, const char *s2, size_t len);
char *strncpy(char * restrict dst, const char * restrict src, size_t len);
char *strrchr(const char *s, int c);
size_t strspn(const char *s, const char *charset);
char *strstr(const char *big, const char *little);
