# Leaven: Compile LLVM IR to Go

Leaven translates LLVM intermediate representation to Go. 
In theory, it should be able to transpile any language that has an LLVM-based compiler to Go.
But so far I’ve only used it for C.

Each LLVM instruction is translated to an equivalent statement in Go.
This produces very verbose code;
if you are looking for a tool that will convert a C codebase into maintainable Go,
Leaven isn’t it.

But it does allow you to call C code from Go without using CGo.
And I am hoping that it can produce a working Go translation of a program,
which will be a good starting point for incrementally re-translating it
(probably by hand) into idiomatic Go.

## Warning

This software is incomplete and experimental.
It does not support nearly all LLVM instructions.

I am not using it myself any more. 
The transpiler at github.com/andybalholm/c2go produces much better results
(but it is not as automatic).

## Usage Example
(Translating `strcmp` from musl libc.)

	$ cat strcmp.c
	#include <string.h>

	int strcmp(const char *l, const char *r)
	{
		for (; *l==*r && *l; l++, r++);
		return *(unsigned char *)l - *(unsigned char *)r;
	}
	$ clang -S -emit-llvm -Os strcmp.c
	$ leaven strcmp.ll
	$ goimports -w strcmp.go
	$ cat strcmp.go
	package main

	import "unsafe"

	var _ unsafe.Pointer

	func strcmp(v0 *byte, v1 *byte) int32 {
		var v10, v11, v12, v13 *byte
		var v5, v6, v7, v16, v17, v18 bool
		var v3, v4, v14, v15, v21, v22 byte
		var v23, v24, v25 int32

		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v3, v4, v5, v6, v7, v10, v11, v12, v13, v14, v15, v16, v17, v18, v21, v22, v23, v24, v25

		v3 = *v0
		v4 = *v1
		v5 = v3 != v4
		v6 = v3 == 0
		v7 = v6 || v5
		if v7 {
			v21, v22 = v3, v4
			goto block20
		} else {
			goto block8
		}

	block8:
		v10, v11 = v1, v0
		goto block9

	block9:
		v12 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(v11)) + 1*unsafe.Sizeof(*(*byte)(nil))))
		v13 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(v10)) + 1*unsafe.Sizeof(*(*byte)(nil))))
		v14 = *v12
		v15 = *v13
		v16 = v14 != v15
		v17 = v14 == 0
		v18 = v17 || v16
		if v18 {
			goto block19
		} else {
			v10, v11 = v13, v12
			goto block9
		}

	block19:
		v21, v22 = v14, v15
		goto block20

	block20:
		v23 = int32(uint32(v21))
		v24 = int32(uint32(v22))
		v25 = v23 - v24
		return v25
	}
