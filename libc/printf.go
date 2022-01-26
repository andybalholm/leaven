package libc

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// fixPrintfFormat converts a printf format string from C-style to Go-style,
// and makes needed changes to the other arguments as well.
//
// It is based on the function of the same name in github.com/andybalholm/c2go.
func fixPrintfFormat(f *byte, args []any) string {
	format := unsafe.Slice(f, Strlen(f))
	var buf strings.Builder
	narg := 0
	start := 0
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		buf.Write(format[start:i])
		start = i
		i++
		if i < len(format) && format[i] == '%' {
			buf.WriteByte('%')
			buf.WriteByte('%')
			start = i + 1
			continue
		}
		for i < len(format) {
			c := format[i]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '#', '-', '.', ',', ' ', 'h', 'j', 'l':
				i++
				continue
			}
			break
		}
		if i >= len(format) {
			// The format string ends with an invalid verb.
			return string(format)
		}
		flags, verb := string(format[start:i]), format[i]
		start = i + 1

		allFlags := flags
		_ = allFlags

		flags = strings.Replace(flags, "h", "", -1)
		flags = strings.Replace(flags, "j", "", -1)
		flags = strings.Replace(flags, "l", "", -1)

		if j := strings.Index(flags, "#0"); j >= 0 && verb == 'x' {
			k := j + 2
			for k < len(flags) && '0' <= flags[k] && flags[k] <= '9' {
				k++
			}
			n, _ := strconv.Atoi(flags[j+2 : k])
			flags = flags[:j+2] + fmt.Sprint(n-2) + flags[k:]
		}

		switch verb {
		default:
			buf.WriteString("%")
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 's':
			if narg < len(args) {
				switch a := args[narg].(type) {
				case *byte:
					args[narg] = unsafe.Slice(a, Strlen(a))
				}
			}
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'f', 'e', 'g', 'c', 'p':
			// usual meanings
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'x', 'X', 'o', 'd', 'b', 'i':
			if verb == 'i' {
				verb = 'd'
			}
			buf.WriteString(flags)
			buf.WriteString(string(verb))

		case 'u':
			buf.WriteString(flags)
			buf.WriteString("d")
			if narg >= len(args) {
				break
			}
			switch a := args[narg].(type) {
			case int16:
				args[narg] = uint16(a)
			case int32:
				args[narg] = uint32(a)
			case int64:
				args[narg] = uint64(a)
			}
		}

		narg++
	}
	buf.Write(format[start:])

	return buf.String()
}

func Printf(format *byte, args ...any) int32 {
	f := fixPrintfFormat(format, args)
	n, err := fmt.Printf(f, args...)
	if err != nil {
		return -1
	}
	return int32(n)
}

func Puts(s *byte) int32 {
	n, err := fmt.Printf("%s\n", unsafe.Slice(s, Strlen(s)))
	if err != nil {
		return -1
	}
	return int32(n)
}
