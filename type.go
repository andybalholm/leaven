package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/types"
)

// TypeDefinition returns the definition (not just the name) of t.
func TypeDefinition(t types.Type) (string, error) {
	switch t := t.(type) {
	case *types.ArrayType:
		elemType, err := TypeSpec(t.ElemType)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%d]%s", t.Len, elemType), nil

	case *types.FloatType:
		switch t.Kind {
		case types.FloatKindFloat:
			return "float32", nil
		case types.FloatKindDouble, types.FloatKindX86_FP80:
			return "float64", nil
		default:
			return "", fmt.Errorf("unsupported floating-point type: %v", t.Kind)
		}

	case *types.FuncType:
		b := new(bytes.Buffer)
		b.WriteString("func(")
		for i, p := range t.Params {
			if i != 0 {
				b.WriteString(", ")
			}
			pt, err := TypeSpec(p)
			if err != nil {
				return "", fmt.Errorf("error converting type of parameter %d (%v): %v", i, p, err)
			}
			b.WriteString(pt)
		}
		b.WriteString(")")
		if !types.Equal(t.RetType, types.Void) {
			b.WriteString(" ")
			rt, err := TypeSpec(t.RetType)
			if err != nil {
				return "", fmt.Errorf("error converting return type (%v): %v", t.RetType, err)
			}
			b.WriteString(rt)
		}
		return b.String(), nil

	case *types.IntType:
		switch t.BitSize {
		case 1:
			return "bool", nil
		case 8:
			return "byte", nil
		default:
			return fmt.Sprintf("int%d", t.BitSize), nil
		}

	case *types.PointerType:
		if _, ok := t.ElemType.(*types.FuncType); ok {
			// Translate a C function pointer type as a Go function type.
			return TypeDefinition(t.ElemType)
		}
		elemType, err := TypeSpec(t.ElemType)
		if err != nil {
			return "", err
		}
		return "*" + elemType, nil

	case *types.StructType:
		b := new(bytes.Buffer)
		b.WriteString("struct {\n")
		for i, field := range t.Fields {
			fieldType, err := TypeSpec(field)
			if err != nil {
				return "", fmt.Errorf("error converting type of field %d: %v", i, err)
			}
			fmt.Fprintf(b, "\tf%d %s\n", i, fieldType)
		}
		b.WriteString("}")
		return b.String(), nil

	default:
		return "", fmt.Errorf("unsupported type %T", t)
	}
}

// TypeSpec returns the name (if it has one) or the definition of t.
func TypeSpec(t types.Type) (string, error) {
	if name := t.Name(); name != "" {
		if name == "union.anon" {
			return TypeDefinition(t)
		}
		name = strings.TrimPrefix(name, "struct.")
		name = strings.TrimPrefix(name, "union.")
		return name, nil
	}
	return TypeDefinition(t)
}
