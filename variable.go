package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// VariableName returns the name to use for a local variable or parameter.
func VariableName(v value.Named) string {
	name := v.Name()
	if name == "" {
		return "v" + strings.TrimPrefix(v.Ident(), "%")
	}
	if c := name[0]; '0' <= c && c <= '9' {
		name = "v" + name
	}
	name = strings.Replace(name, ".", "_", -1)
	return name
}

// FormatValue formats a constant or variable as it should appear in an expression.
func FormatValue(v value.Value) (string, error) {
	switch v := v.(type) {
	case value.Named:
		return VariableName(v), nil

	case *ir.Arg:
		return FormatValue(v.Value)

	case *constant.Array:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", v.Typ, err)
		}
		b := new(bytes.Buffer)
		if len(v.Elems) < 16 {
			b.WriteString(t)
			b.WriteByte('{')
			for i, c := range v.Elems {
				if i > 0 {
					b.WriteString(", ")
				}
				e, err := FormatValue(c)
				if err != nil {
					return "", fmt.Errorf("error translating element %d (%v): %v", i, c, err)
				}
				fmt.Fprint(b, e)
			}
			b.WriteByte('}')
		} else {
			b.WriteString(t)
			b.WriteString("{\n\t")
			for i, c := range v.Elems {
				if i > 0 {
					if i%16 == 0 {
						b.WriteString(",\n\t")
					} else {
						b.WriteString(", ")
					}
				}
				e, err := FormatValue(c)
				if err != nil {
					return "", fmt.Errorf("error translating element %d (%v): %v", i, c, err)
				}
				fmt.Fprint(b, e)
			}
			b.WriteString(",\n}")
		}
		return b.String(), nil

	case *constant.CharArray:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", v.Typ, err)
		}
		b := new(bytes.Buffer)
		if len(v.X) < 16 {
			b.WriteString(t)
			b.WriteByte('{')
			for i, c := range v.X {
				if i > 0 {
					b.WriteString(", ")
				}
				fmt.Fprintf(b, "%d", c)
			}
			b.WriteByte('}')
		} else {
			b.WriteString(t)
			b.WriteString("{\n\t")
			for i, c := range v.X {
				if i > 0 {
					if i%16 == 0 {
						b.WriteString(",\n\t")
					} else {
						b.WriteString(", ")
					}
				}
				fmt.Fprintf(b, "%d", c)
			}
			b.WriteString(",\n}")
		}
		return b.String(), nil

	case *constant.ExprBitCast:
		from, err := FormatValue(v.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", v.From, err)
		}
		switch v.From.(type) {
		case *ir.Global:
			from = "&" + from
		}
		to, err := TypeSpec(v.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", v.To, err)
		}
		return fmt.Sprintf("(%s)(unsafe.Pointer(%s))", to, from), nil

	case *constant.ExprGetElementPtr:
		indices := make([]value.Value, len(v.Indices))
		for i, index := range v.Indices {
			indices[i] = index
		}
		return GetElementPtr(v.ElemType, v.Src, indices)

	case *constant.Float:
		result := v.X.String()
		special := false
		switch result {
		case "+Inf":
			result = "math.Inf(1)"
			special = true
		case "-Inf":
			result = "math.Inf(-1)"
			special = true
		case "NaN":
			result = "math.NaN()"
			special = true
		}
		if special && v.Typ.Kind == types.FloatKindFloat {
			result = fmt.Sprintf("float32(%s)", result)
		}
		return result, nil

	case *constant.Index:
		return FormatValue(v.Constant)

	case *constant.Int:
		result := v.X.String()
		if v.Typ.BitSize == 1 {
			if result == "1" {
				result = "true"
			} else {
				result = "false"
			}
		}
		return result, nil

	case *constant.Null:
		return "nil", nil

	case *constant.Struct:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", v.Typ, err)
		}
		b := new(bytes.Buffer)
		b.WriteString(t)
		b.WriteByte('{')
		for i, c := range v.Fields {
			if i > 0 {
				b.WriteString(", ")
			}
			e, err := FormatValue(c)
			if err != nil {
				return "", fmt.Errorf("error translating field %d (%v): %v", i, c, err)
			}
			fmt.Fprint(b, e)
		}
		b.WriteByte('}')
		return b.String(), nil

	case *constant.Vector:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", v.Typ, err)
		}
		b := new(bytes.Buffer)
		b.WriteString(t)
		b.WriteByte('{')
		for i, c := range v.Elems {
			if i > 0 {
				b.WriteString(", ")
			}
			e, err := FormatValue(c)
			if err != nil {
				return "", fmt.Errorf("error translating element %d (%v): %v", i, c, err)
			}
			fmt.Fprint(b, e)
		}
		b.WriteByte('}')
		return b.String(), nil

	case *constant.ZeroInitializer:
		t, err := TypeSpec(v.Typ)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", v.Typ, err)
		}
		return t + "{}", nil

	default:
		return "", fmt.Errorf("unsupported type of value to translate: %T", v)
	}
}

// FormatSigned is like FormatValue, except that it converts "byte" to "int8".
func FormatSigned(v value.Value) (string, error) {
	result, err := FormatValue(v)
	if err != nil {
		return "", err
	}

	if _, ok := v.(*constant.Int); ok {
		return result, nil
	}

	if t, ok := v.Type().(*types.IntType); ok && t.BitSize == 8 {
		return fmt.Sprintf("int8(%s)", result), nil
	}
	return result, nil
}

// FormatUnsigned is like FormatValue, except that it converts integer types to
// unsigned.
func FormatUnsigned(v value.Value) (string, error) {
	result, err := FormatValue(v)
	if err != nil {
		return "", err
	}

	if _, ok := v.(*constant.Int); ok {
		return result, nil
	}

	if t, ok := v.Type().(*types.IntType); ok && t.BitSize > 8 {
		return fmt.Sprintf("uint%d(%s)", t.BitSize, result), nil
	}
	return result, nil
}
