package libc

import (
	"bytes"
	"unsafe"
)

// Memmove copies length bytes from src to dst. The blocks of memory may
// overlap.
func Memmove(dst *byte, src *byte, length int64) *byte {
	copy(unsafe.Slice(dst, int(length)), unsafe.Slice(src, int(length)))
	return dst
}

// Memset fills the memory at b with bytes of the value c.
func Memset(b *byte, c byte, length int64) *byte {
	dest := unsafe.Slice(b, int(length))
	for i := range dest {
		dest[i] = c
	}
	return b
}

// MemsetChk fills the memory at b with bytes of the value c. If length is
// greater than destlen (interpreted as unsigned integers), it will panic.
func MemsetChk(b *byte, c byte, length int64, destlen int64) *byte {
	if uint64(length) > uint64(destlen) {
		panic("buffer overflow")
	}
	return Memset(b, c, length)
}

// MemsetPattern16 fills the memory at b with a 16-byte pattern.
func MemsetPattern16(b *byte, pattern16 *byte, length int64) {
	memsetPattern(unsafe.Slice(b, int(length)), unsafe.Slice(pattern16, 16))
}

func memsetPattern(dest, pattern []byte) {
	for i := 0; i < len(dest); i += len(pattern) {
		copy(dest[i:], pattern)
	}
}

// StrcatChk concatenates two C strings. It panics if the space required is
// more than destlen bytes.
func StrcatChk(dest *byte, src *byte, destlen int64) *byte {
	d := unsafe.Slice(dest, int(destlen))
	s := unsafe.Slice(src, int(destlen+1)) // The +1 ensures a panic if dest is an empty string, and src is too long.

	// Find the end of dest.
	i := 0
	for d[i] != 0 {
		i++
	}

	for _, c := range s {
		d[i] = c
		i++
		if c == 0 {
			break
		}
	}

	return dest
}

// MemcpyChk copies length bytes from src to dest. If length is greater than
// destlen (interpreted as unsigned integers), it will panic.
func MemcpyChk(dest *byte, src *byte, length int64, destlen int64) *byte {
	if uint64(length) > uint64(destlen) {
		panic("buffer overflow")
	}
	copy(unsafe.Slice(dest, int(length)), unsafe.Slice(src, int(length)))
	return dest
}

// Memchr returns a pointer to the first occurrence of c in string s.
// It returns nil if no such byte exists within n bytes.
func Memchr(s *byte, c int32, n int64) *byte {
	b := unsafe.Slice(s, int(n))
	i := bytes.IndexByte(b, byte(c))
	if i == -1 {
		return nil
	}
	return &b[i]
}
