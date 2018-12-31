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
