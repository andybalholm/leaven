package libc

var allocated = make(map[*byte][]byte)

// Malloc allocates size bytes of memory, and returns a pointer to the
// allocated memory. The memory will not be garbage-collected; it must be
// release by a call to free.
func Malloc(size int64) *byte {
	if size == 0 {
		return nil
	}
	b := make([]byte, size)
	p := &b[0]
	allocated[p] = b
	return p
}

// Free releases memory allocated by Malloc.
func Free(p *byte) {
	delete(allocated, p)
}

// Calloc allocates a block of memory for count objects of size bytes each.
func Calloc(count, size int64) *byte {
	return Malloc(count * size)
}
