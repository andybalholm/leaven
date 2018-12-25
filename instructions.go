package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

// TranslateInstruction translates an LLVM instruction to Go.
func TranslateInstruction(inst ir.Instruction) (string, error) {
	switch inst := inst.(type) {
	case *ir.InstAdd:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.X, err)
		}
		return fmt.Sprintf("%s = %s + %s", VariableName(inst), x, y), nil

	case *ir.InstAlloca:
		t, err := TypeSpec(inst.ElemType)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.ElemType, err)
		}
		if inst.NElems == nil {
			return fmt.Sprintf("%s = new(%s)", VariableName(inst), t), nil
		}
		nElems, err := FormatValue(inst.NElems)
		if err != nil {
			return "", fmt.Errorf("error translating NElems (%v): %v", inst.NElems, err)
		}
		return fmt.Sprintf("%s = &make([]%s, %s)[0]", VariableName(inst), t, nElems), nil

	case *ir.InstAnd:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.X, err)
		}
		return fmt.Sprintf("%s = %s & %s", VariableName(inst), x, y), nil

	case *ir.InstGetElementPtr:
		srcPointerType, ok := inst.Src.Type().(*types.PointerType)
		if !ok {
			return "", fmt.Errorf("non-pointer source parameter: %v", inst.Src.Type())
		}
		if !types.Equal(srcPointerType.ElemType, inst.ElemType) {
			return "", fmt.Errorf("mismatched source and element types")
		}

		zeroFirstIndex := false
		positiveFirstIndex := false
		if firstIndex, ok := inst.Indices[0].(*constant.Int); ok {
			switch firstIndex.X.Sign() {
			case 0:
				zeroFirstIndex = true
			case 1:
				positiveFirstIndex = true
			}
		}
		takeAddress := false

		source, err := FormatValue(inst.Src)
		if err != nil {
			return "", fmt.Errorf("error translating source pointer (%q): %v", inst.Src, err)
		}
		result := source

		if !zeroFirstIndex {
			firstIndex, err := FormatValue(inst.Indices[0])
			if err != nil {
				return "", fmt.Errorf("error translating first index (%v): %v", inst.Indices[0], err)
			}
			elemType, err := TypeSpec(inst.ElemType)
			if err != nil {
				return "", fmt.Errorf("error translating element type (%v): %v", inst.ElemType, err)
			}
			offset := fmt.Sprintf("uintptr(int64(%s)) * unsafe.Sizeof(*(*%s)(nil))", firstIndex, elemType)
			if positiveFirstIndex {
				offset = fmt.Sprintf("%s * unsafe.Sizeof(*(*%s)(nil))", firstIndex, elemType)
			}
			result = fmt.Sprintf("uintptr(unsafe.Pointer(%s)) + %s", source, offset)
			result = fmt.Sprintf("(*%s)(unsafe.Pointer(%s))", elemType, result)
		}

		currentType := inst.ElemType

		for _, index := range inst.Indices[1:] {
			switch ct := currentType.(type) {
			case *types.StructType:
				ci, ok := index.(*constant.Int)
				if !ok {
					return "", fmt.Errorf("non-constant index into struct: %v", index)
				}
				result = fmt.Sprintf("%s.f%v", result, ci.X)
				currentType = ct.Fields[ci.X.Int64()]
				takeAddress = true

			default:
				return "", fmt.Errorf("unsupported type to index into: %v", currentType)
			}
		}

		if takeAddress {
			result = fmt.Sprintf("%s = &%s", VariableName(inst), result)
		} else {
			result = fmt.Sprintf("%s = %s", VariableName(inst), result)
		}

		return result, nil

	case *ir.InstICmp:
		var op string
		format := FormatValue
		switch inst.Pred {
		case enum.IPredEQ:
			op = "=="
		case enum.IPredNE:
			op = "!="
		case enum.IPredSGE:
			op = ">="
			format = FormatSigned
		case enum.IPredSGT:
			op = ">"
			format = FormatSigned
		case enum.IPredSLE:
			op = "<="
			format = FormatSigned
		case enum.IPredSLT:
			op = "<"
			format = FormatSigned
		case enum.IPredUGE:
			op = ">="
			format = FormatUnsigned
		case enum.IPredUGT:
			op = ">"
			format = FormatUnsigned
		case enum.IPredULE:
			op = "<="
			format = FormatUnsigned
		case enum.IPredULT:
			op = "<"
			format = FormatUnsigned
		default:
			return "", fmt.Errorf("unsupported comparison predicate: %v", inst.Pred)
		}

		x, err := format(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := format(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.X, err)
		}
		return fmt.Sprintf("%s = %s %s %s", VariableName(inst), x, op, y), nil

	case *ir.InstLoad:
		src, err := FormatValue(inst.Src)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.Src, err)
		}
		return fmt.Sprintf("%s = *%s", VariableName(inst), src), nil

	case *ir.InstLShr:
		x, err := FormatUnsigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.X, err)
		}
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize > 8 {
			return fmt.Sprintf("%s = int%d(%s >> %s)", VariableName(inst), t.BitSize, x, y), nil
		}
		return fmt.Sprintf("%s = %s >> %s", VariableName(inst), x, y), nil

	case *ir.InstOr:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.X, err)
		}
		return fmt.Sprintf("%s = %s | %s", VariableName(inst), x, y), nil

	case *ir.InstPtrToInt:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(uintptr(unsafe.Pointer(%s)))", VariableName(inst), to, from), nil

	case *ir.InstShl:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.X, err)
		}
		return fmt.Sprintf("%s = %s << %s", VariableName(inst), x, y), nil

	case *ir.InstStore:
		dest, err := FormatValue(inst.Dst)
		if err != nil {
			return "", fmt.Errorf("error translating destination (%v): %v", inst.Dst, err)
		}
		src, err := FormatValue(inst.Src)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.Src, err)
		}
		return fmt.Sprintf("*%s = %s", dest, src), nil

	case *ir.InstSub:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.X, err)
		}
		return fmt.Sprintf("%s = %s - %s", VariableName(inst), x, y), nil

	case *ir.InstZExt:
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return "", fmt.Errorf("unsupported To type for zext: %T", inst.To)
		}
		from, err := FormatUnsigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		return fmt.Sprintf("%s = int%d(uint%d(%s))", VariableName(inst), toType.BitSize, toType.BitSize, from), nil

	default:
		return "", fmt.Errorf("unsupported instruction type: %T", inst)
	}
}
