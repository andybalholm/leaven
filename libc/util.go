// The libc package implements various functions from the C standard library
// in Go.
package libc

import "unsafe"

// byteSlice returns a slice of n bytes, starting at p.
func byteSlice(p *byte, n int) []byte {
	return (*[1 << 30]byte)(unsafe.Pointer(p))[:n:n]
}

// GoString returns s converted from a C string to a Go string.
func GoString(s *byte) string {
	return string(byteSlice(s, int(Strlen(s))))
}
