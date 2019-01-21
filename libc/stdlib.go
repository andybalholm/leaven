package libc

import (
	"sync"

	"golang.org/x/sys/unix"
)

var mallocLock sync.Mutex
var allocated = make(map[*byte][]byte)

// Malloc allocates size bytes of memory, and returns a pointer to the
// allocated memory. The memory will not be garbage-collected; it must be
// released by a call to free.
func Malloc(size int64) *byte {
	if size == 0 {
		return nil
	}
	mallocLock.Lock()
	defer mallocLock.Unlock()

	b, err := unix.Mmap(0, 0, int(size), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON)
	if err != nil {
		panic(err)
	}
	p := &b[0]
	allocated[p] = b
	return p
}

// Free releases memory allocated by Malloc.
func Free(p *byte) {
	mallocLock.Lock()
	defer mallocLock.Unlock()

	if b, ok := allocated[p]; ok {
		if err := unix.Munmap(b); err != nil {
			panic(err)
		}
		delete(allocated, p)
	}
}

// Calloc allocates a block of memory for count objects of size bytes each.
func Calloc(count, size int64) *byte {
	return Malloc(count * size)
}
