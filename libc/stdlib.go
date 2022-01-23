package libc

import "unsafe"

// Malloc allocates n bytes of memory. It informs the garbage collector that
// the memory will be used to store objects of type T.
func Malloc[T any](n int64) *T {
	var p *T
	if uintptr(n) == unsafe.Sizeof(*p) {
		return new(T)
	}
	// Allocate one extra element to allow indexing off the end, like C tends
	// to do.
	count := uintptr(n)/unsafe.Sizeof(*p) + 1
	return &make([]T, count)[0]
}

// Calloc allocates a block of memory for count objects of size bytes each.
func Calloc[T any](count, size int64) *T {
	return Malloc[T](count * size)
}
