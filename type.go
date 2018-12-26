package main

import (
	"bytes"
	"fmt"

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
		return name, nil
	}
	return TypeDefinition(t)
}
