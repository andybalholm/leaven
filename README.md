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
	$ clang -S -emit-llvm -Os -fno-discard-value-names strcmp.c
	$ leaven strcmp.ll
	$ goimports -w strcmp.go
	$ cat strcmp.go
	package main

	import "unsafe"

	func strcmp(l *byte, r *byte) int32 {
		var r_addr_017, l_addr_016, incdec_ptr, incdec_ptr4 *byte
		var cmp13, tobool14, or_cond15, cmp, tobool, or_cond bool
		var v0, v1, v2, v3, _lcssa12, _lcssa byte
		var conv5, conv6, sub int32

		_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v0, v1, cmp13, tobool14, or_cond15, r_addr_017, l_addr_016, incdec_ptr, incdec_ptr4, v2, v3, cmp, tobool, or_cond, _lcssa12, _lcssa, conv5, conv6, sub

		v0 = *l
		v1 = *r
		cmp13 = v0 != v1
		tobool14 = v0 == 0
		or_cond15 = tobool14 || cmp13
		if or_cond15 {
			_lcssa12, _lcssa = v0, v1
			goto for_end
		} else {
			r_addr_017, l_addr_016 = r, l
			goto for_inc
		}

	for_inc:
		incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(l_addr_016)) + 1*unsafe.Sizeof(*(*byte)(nil))))
		incdec_ptr4 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(r_addr_017)) + 1*unsafe.Sizeof(*(*byte)(nil))))
		v2 = *incdec_ptr
		v3 = *incdec_ptr4
		cmp = v2 != v3
		tobool = v2 == 0
		or_cond = tobool || cmp
		if or_cond {
			_lcssa12, _lcssa = v2, v3
			goto for_end
		} else {
			r_addr_017, l_addr_016 = incdec_ptr4, incdec_ptr
			goto for_inc
		}

	for_end:
		conv5 = int32(uint32(_lcssa12))
		conv6 = int32(uint32(_lcssa))
		sub = conv5 - conv6
		return sub
	}
