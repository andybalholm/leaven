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
	if invalidNames[name] {
		name = "_" + name
	}
	return name
}

func BlockName(v value.Value) string {
	block := v.(*ir.Block)
	name := block.Name()
	if name == "" {
		return "block" + strings.TrimPrefix(block.Ident(), "%")
	}
	if c := name[0]; '0' <= c && c <= '9' {
		name = "block" + name
	}
	name = strings.Replace(name, ".", "_", -1)
	if invalidNames[name] {
		name = "_" + name
	}
	return name
}

var invalidNames = map[string]bool{
	"return": true,
}

// FormatValue formats a constant or variable as it should appear in an expression.
func FormatValue(v value.Value) (string, error) {
	switch v := v.(type) {
	case *ir.Global:
		if types.IsFunc(v.ContentType) {
			return VariableName(v), nil
		}
		return "&" + VariableName(v), nil

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
		var value int64
		switch {
		case v.X.IsInt64():
			value = v.X.Int64()
		case v.X.IsUint64():
			value = int64(v.X.Uint64())
		default:
			return "", fmt.Errorf("integer constant too large: %v", v.X)
		}

		switch v.Typ.BitSize {
		case 1:
			if value != 0 {
				return "true", nil
			}
			return "false", nil
		case 8:
			return fmt.Sprint(byte(value)), nil
		case 16:
			return fmt.Sprint(int16(value)), nil
		case 32:
			return fmt.Sprint(int32(value)), nil
		default:
			return fmt.Sprint(value), nil
		}

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

	case *constant.Undef:
		switch v.Typ.(type) {
		case *types.ArrayType, *types.StructType, *types.VectorType:
			t, err := TypeSpec(v.Typ)
			if err != nil {
				return "", fmt.Errorf("error translating type (%v): %v", v.Typ, err)
			}
			return t + "{}", nil
		case *types.IntType, *types.FloatType:
			return "0", nil
		case *types.PointerType:
			return "nil", nil
		default:
			return "", fmt.Errorf("unsupported type for undefined constant: %v", v.Typ)
		}

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

	if ci, ok := v.(*constant.Int); ok {
		if ci.Typ.BitSize == 8 {
			return fmt.Sprint(int8(ci.X.Int64())), nil
		}
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

	if ci, ok := v.(*constant.Int); ok {
		var value uint64
		switch {
		case ci.X.IsUint64():
			value = ci.X.Uint64()
		case ci.X.IsInt64():
			value = uint64(ci.X.Int64())
			switch ci.Typ.BitSize {
			case 8:
				return fmt.Sprintf("byte(%d)", byte(value)), nil
			case 16:
				return fmt.Sprintf("uint16(%d)", uint16(value)), nil
			case 32:
				return fmt.Sprintf("uint32(%d)", uint32(value)), nil
			default:
				return fmt.Sprint(value), nil
			}
		default:
			return "", fmt.Errorf("integer constant too large: %v", ci.X)
		}

		switch ci.Typ.BitSize {
		case 1:
			if value != 0 {
				return "true", nil
			}
			return "false", nil
		case 8:
			return fmt.Sprint(byte(value)), nil
		case 16:
			return fmt.Sprint(uint16(value)), nil
		case 32:
			return fmt.Sprint(uint32(value)), nil
		default:
			return fmt.Sprint(value), nil
		}
	}

	if t, ok := v.Type().(*types.IntType); ok && t.BitSize > 8 {
		return fmt.Sprintf("uint%d(%s)", t.BitSize, result), nil
	}
	return result, nil
}
