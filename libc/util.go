// The libc package implements various functions from the C standard library
// in Go.
package libc

import "unsafe"

// byteSlice returns a slice of n bytes, starting at p.
func byteSlice(p *byte, n int) []byte {
	return (*[2 << 30]byte)(unsafe.Pointer(p))[:n:n]
}
