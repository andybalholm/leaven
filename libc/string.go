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

// StrcatChk concatenates two C strings. It panics if the space required is
// more than destlen bytes.
func StrcatChk(dest *byte, src *byte, destlen int64) *byte {
	d := byteSlice(dest, int(destlen))
	s := byteSlice(src, int(destlen+1)) // The +1 ensures a panic if dest is an empty string, and src is too long.

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
	copy(byteSlice(dest, int(length)), byteSlice(src, int(length)))
	return dest
}
