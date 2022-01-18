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
	$ clang -S -emit-llvm -fno-discard-value-names strcmp.c
	$ leaven strcmp.ll
	$ goimports -w strcmp.go
	$ cat strcmp.go
	package main

	import "unsafe"

	func strcmp(l *byte, r *byte) int32 {
		var l_addr, r_addr **byte
		var v0, v2, v4, v7, incdec_ptr, v8, incdec_ptr4, v9, v11 *byte
		var cmp, tobool, v6 bool
		var v1, v3, v5, v10, v12 byte
		var conv, conv1, conv3, conv5, conv6, sub int32

		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = l_addr, r_addr, v0, v1, conv, v2, v3, conv1, cmp, v4, v5, conv3, tobool, v6, v7, incdec_ptr, v8, incdec_ptr4, v9, v10, conv5, v11, v12, conv6, sub

		l_addr = new(*byte)
		r_addr = new(*byte)
		*l_addr = l
		*r_addr = r
		goto for_cond

	for_cond:
		v0 = *l_addr
		v1 = *v0
		conv = int32(int8(v1))
		v2 = *r_addr
		v3 = *v2
		conv1 = int32(int8(v3))
		cmp = conv == conv1
		if cmp {
			goto land_rhs
		} else {
			v6 = false
			goto land_end
		}

	land_rhs:
		v4 = *l_addr
		v5 = *v4
		conv3 = int32(int8(v5))
		tobool = conv3 != 0
		v6 = tobool
		goto land_end

	land_end:
		if v6 {
			goto for_body
		} else {
			goto for_end
		}

	for_body:
		goto for_inc

	for_inc:
		v7 = *l_addr
		incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(v7)) + 1*unsafe.Sizeof(*(*byte)(nil))))
		*l_addr = incdec_ptr
		v8 = *r_addr
		incdec_ptr4 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + 1*unsafe.Sizeof(*(*byte)(nil))))
		*r_addr = incdec_ptr4
		goto for_cond

	for_end:
		v9 = *l_addr
		v10 = *v9
		conv5 = int32(uint32(v10))
		v11 = *r_addr
		v12 = *v11
		conv6 = int32(uint32(v12))
		sub = conv5 - conv6
		return sub
	}

