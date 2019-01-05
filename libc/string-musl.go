package libc

import "unsafe"

// The functions in this file are transpiled from musl libc, which has the
// following license:

/*
Copyright © 2005-2014 Rich Felker, et al.

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

func Strcmp(v0 *byte, v1 *byte) int32 {
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

func Strlen(v0 *byte) int64 {
	var v8, v12, v18, v31, v35, v36, v41 *byte
	var v19, v21, v28 *int64
	var v4, v10, v15, v27, v32, v38 bool
	var v9, v30, v37 byte
	var v2, v3, v7, v13, v14, v22, v23, v24, v25, v26, v42, v45, v46 int64

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v2, v3, v4, v7, v8, v9, v10, v12, v13, v14, v15, v18, v19, v21, v22, v23, v24, v25, v26, v27, v28, v30, v31, v32, v35, v36, v37, v38, v41, v42, v45, v46

	v2 = int64(uintptr(unsafe.Pointer(v0)))
	v3 = v2 & 7
	v4 = v3 == 0
	if v4 {
		v18 = v0
		goto block17
	} else {
		goto block5
	}

block5:
	v7, v8 = v2, v0
	goto block6

block6:
	v9 = *v8
	v10 = v9 == 0
	if v10 {
		goto block43
	} else {
		goto block11
	}

block11:
	v12 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(v8)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	v13 = int64(uintptr(unsafe.Pointer(v12)))
	v14 = v13 & 7
	v15 = v14 == 0
	if v15 {
		goto block16
	} else {
		v7, v8 = v13, v12
		goto block6
	}

block16:
	v18 = v12
	goto block17

block17:
	v19 = (*int64)(unsafe.Pointer(v18))
	v21 = v19
	goto block20

block20:
	v22 = *v21
	v23 = v22 - 72340172838076673
	v24 = v22 & -9187201950435737472
	v25 = v24 ^ -9187201950435737472
	v26 = v25 & v23
	v27 = v26 == 0
	v28 = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v21)) + 1*unsafe.Sizeof(*(*int64)(nil))))
	if v27 {
		v21 = v28
		goto block20
	} else {
		goto block29
	}

block29:
	v30 = byte(v22)
	v31 = (*byte)(unsafe.Pointer(v21))
	v32 = v30 == 0
	if v32 {
		v41 = v31
		goto block40
	} else {
		goto block33
	}

block33:
	v35 = v31
	goto block34

block34:
	v36 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(v35)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	v37 = *v36
	v38 = v37 == 0
	if v38 {
		goto block39
	} else {
		v35 = v36
		goto block34
	}

block39:
	v41 = v36
	goto block40

block40:
	v42 = int64(uintptr(unsafe.Pointer(v41)))
	v45 = v42
	goto block44

block43:
	v45 = v7
	goto block44

block44:
	v46 = v45 - v2
	return v46
}
