// The libc package implements various functions from the C standard library
// in Go.
package libc

import "unsafe"

// GoString returns s converted from a C string to a Go string.
func GoString(s *byte) string {
	return string(unsafe.Slice(s, Strlen(s)))
}

// AddPointer does C-style pointer addition: it multiplies offset by
// sizeof(*ptr) and adds it to ptr.
func AddPointer[T any](ptr *T, offset int) *T {
	return (*T)(unsafe.Add(unsafe.Pointer(ptr), offset*int(unsafe.Sizeof(*ptr))))
}
