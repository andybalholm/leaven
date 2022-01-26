package libc

import "os"

func Getchar() int32 {
	var buf [1]byte
	_, err := os.Stdin.Read(buf[:])
	if err != nil {
		return -1
	}
	return int32(buf[0])
}
