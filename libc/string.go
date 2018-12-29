package libc

// MemsetPattern16 fills the memory at b with a 16-byte pattern.
func MemsetPattern16(b *byte, pattern16 *byte, length int64) {
	memsetPattern(byteSlice(b, int(length)), byteSlice(pattern16, 16))
}

func memsetPattern(dest, pattern []byte) {
	for i := 0; i < len(dest); i += len(pattern) {
		copy(dest[i:], pattern)
	}
}
