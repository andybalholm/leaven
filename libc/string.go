package libc

// Memmove copies length bytes from src to dst. The blocks of memory may
// overlap.
func Memmove(dst *byte, src *byte, length int64) *byte {
	copy(byteSlice(dst, int(length)), byteSlice(src, int(length)))
	return dst
}

// Memset fills the memory at b with bytes of the value c.
func Memset(b *byte, c byte, length int64) *byte {
	dest := byteSlice(b, int(length))
	for i := range dest {
		dest[i] = c
	}
	return b
}

// MemsetPattern16 fills the memory at b with a 16-byte pattern.
func MemsetPattern16(b *byte, pattern16 *byte, length int64) {
	memsetPattern(byteSlice(b, int(length)), byteSlice(pattern16, 16))
}

func memsetPattern(dest, pattern []byte) {
	for i := 0; i < len(dest); i += len(pattern) {
		copy(dest[i:], pattern)
	}
}
