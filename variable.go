package main

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// VariableName returns the name to use for a local variable or parameter.
func VariableName(v value.Named) string {
	if name := v.Name(); name != "" {
		return name
	}
	return "v" + strings.TrimPrefix(v.Ident(), "%")
}

// FormatValue formats a constant or variable as it should appear in an expression.
func FormatValue(v value.Value) (string, error) {
	switch v := v.(type) {
	case value.Named:
		return VariableName(v), nil

	case *constant.Int:
		return v.X.String(), nil

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
