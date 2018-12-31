package libc

import (
	"reflect"
	"unsafe"
)

// VAArg returns a pointer to the next argument in a varargs list. The actual
// type of list is *[]interface{}, but it is declared as void * in C.
func VAArg(list *byte) *byte {
	vl := (*[]interface{})(unsafe.Pointer(list))
	arg := (*vl)[0]
	*vl = (*vl)[1:]

	var intVal int32
	switch arg := arg.(type) {
	case byte:
		intVal = int32(arg)
		return (*byte)(unsafe.Pointer(&intVal))
	case int16:
		intVal = int32(arg)
		return (*byte)(unsafe.Pointer(&intVal))
	case int32:
		intVal = arg
		return (*byte)(unsafe.Pointer(&intVal))
	}

	// Use reflect to make a copy of arg that we can take the address of.
	av := reflect.ValueOf(arg)
	p := reflect.New(av.Type())
	p.Elem().Set(av)
	return (*byte)(unsafe.Pointer((p.Elem().UnsafeAddr())))
}
