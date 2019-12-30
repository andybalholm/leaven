package main

import (
	"fmt"
	"strings"

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
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v + %s[i] }", x, VariableName(inst), y), nil
		}
		if ciy, ok := inst.Y.(*constant.Int); ok && ciy.X.Sign() == -1 {
			return fmt.Sprintf("%s = %s %s", VariableName(inst), x, ciy.X), nil // Use the constant's own minus sign.
		}
		return fmt.Sprintf("%s = %s + %s", VariableName(inst), x, y), nil

	case *ir.InstAlloca:
		t, err := TypeSpec(inst.ElemType)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.ElemType, err)
		}
		if inst.NElems == nil {
			return fmt.Sprintf("%s = (*%s)(unsafe.Pointer(libc.Malloc(int64(unsafe.Sizeof(*(*%s)(nil)))))); defer libc.Free((*byte)(unsafe.Pointer(%s)))", VariableName(inst), t, t, VariableName(inst)), nil
		}
		nElems, err := FormatValue(inst.NElems)
		if err != nil {
			return "", fmt.Errorf("error translating NElems (%v): %v", inst.NElems, err)
		}
		return fmt.Sprintf("%s = (*%s)(unsafe.Pointer(libc.Malloc(int64(%s * unsafe.Sizeof(*(*%s)(nil)))))); defer libc.Free((*byte)(unsafe.Pointer(%s)))", VariableName(inst), t, nElems, t, VariableName(inst)), nil

	case *ir.InstAnd:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v & %s[i] }", x, VariableName(inst), y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return fmt.Sprintf("%s = %s && %s", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s & %s", VariableName(inst), x, y), nil

	case *ir.InstAShr:
		x, err := FormatSigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize == 8 {
			return fmt.Sprintf("%s = byte(%s >> %s)", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s >> %s", VariableName(inst), x, y), nil

	case *ir.InstBitCast:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		return fmt.Sprintf("%s = (%s)(unsafe.Pointer(%s))", VariableName(inst), to, from), nil

	case *ir.InstCall:
		callee, err := FormatValue(inst.Callee)
		if err != nil {
			return "", fmt.Errorf("error translating callee (%v): %v", inst.Callee, err)
		}
		args := make([]string, len(inst.Args))
		for i, a := range inst.Args {
			v, err := FormatValue(a)
			if err != nil {
				return "", fmt.Errorf("error translating argument %d (%v): %v", i, a, err)
			}
			args[i] = v
		}
		if renamed, ok := libraryFunctions[callee]; ok {
			callee = renamed
		}
		switch callee {
		case "leaven_va_start":
			if len(args) == 1 {
				return fmt.Sprintf("*%s = (*byte)(unsafe.Pointer(&varargs))", args[0]), nil
			}
		case "ldexp":
			if len(args) == 2 {
				return fmt.Sprintf("%s = math.Ldexp(%s, int(%s))", VariableName(inst), args[0], args[1]), nil
			}
		case "llvm_fabs_f32":
			if len(args) == 1 {
				return fmt.Sprintf("%s = float32(math.Abs(float64(%s)))", VariableName(inst), args[0]), nil
			}
		case "llvm_lifetime_start", "llvm_lifetime_end":
			return ";", nil
		case "llvm_memcpy_p0i8_p0i8_i64":
			return fmt.Sprintf("libc.Memmove(%s, %s, %s)", args[0], args[1], args[2]), nil
		case "llvm_memset_p0i8_i64":
			return fmt.Sprintf("libc.Memset(%s, %s, %s)", args[0], args[1], args[2]), nil
		case "llvm_objectsize_i64_p0i8":
			// Use -1 for unknown size.
			return fmt.Sprintf("%s = -1", VariableName(inst)), nil
		case "putchar":
			if len(args) == 1 {
				return fmt.Sprintf("if _, err := os.Stdout.Write([]byte{byte(%s)}); err != nil { %s = -1 } else { %s = %s }", args[0], VariableName(inst), VariableName(inst), args[0]), nil
			}
		case "__sprintf_chk":
			return fmt.Sprintf("%s = noarch.Snprintf(%s, %s)", VariableName(inst), args[0], strings.Join(args[2:], ", ")), nil
		}
		if types.Equal(inst.Type(), types.Void) {
			return fmt.Sprintf("%s(%s)", callee, strings.Join(args, ", ")), nil
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), callee, strings.Join(args, ", ")), nil

	case *ir.InstExtractElement:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating vector (%v): %v", inst.X, err)
		}
		index, err := FormatValue(inst.Index)
		if err != nil {
			return "", fmt.Errorf("error translating index (%v): %v", inst.Index, err)
		}
		return fmt.Sprintf("%s = %s[%s]", VariableName(inst), x, index), nil

	case *ir.InstFCmp:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}

		var op string
		switch inst.Pred {
		case enum.FPredOEQ:
			op = "=="
		case enum.FPredOGE:
			op = ">="
		case enum.FPredOGT:
			op = ">"
		case enum.FPredOLE:
			op = "<="
		case enum.FPredOLT:
			op = "<"
		case enum.FPredUNE:
			op = "!="
		case enum.FPredORD:
			return fmt.Sprintf("%s = %s == %s && %s == %s", VariableName(inst), x, x, y, y), nil
		case enum.FPredUNO:
			return fmt.Sprintf("%s = %s != %s || %s != %s", VariableName(inst), x, x, y, y), nil
		case enum.FPredUEQ:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s == %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredUGT:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s > %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredUGE:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s >= %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredULT:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s < %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredULE:
			return fmt.Sprintf("%s = %s != %s || %s != %s || %s <= %s", VariableName(inst), x, x, y, y, x, y), nil
		case enum.FPredONE:
			return fmt.Sprintf("%s = %s == %s && %s == %s && %s != %s", VariableName(inst), x, x, y, y, x, y), nil
		default:
			return "", fmt.Errorf("unsupported comparison predicate: %v", inst.Pred)
		}

		return fmt.Sprintf("%s = %s %s %s", VariableName(inst), x, op, y), nil

	case *ir.InstFDiv:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s / %s", VariableName(inst), x, y), nil

	case *ir.InstFMul:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s * %s", VariableName(inst), x, y), nil

	case *ir.InstFPExt:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstFPToSI:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		if to == "byte" {
			return fmt.Sprintf("%s = byte(int8(%s))", VariableName(inst), from), nil
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstFPTrunc:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstFSub:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s - %s", VariableName(inst), x, y), nil

	case *ir.InstGetElementPtr:
		result, err := GetElementPtr(inst.ElemType, inst.Src, inst.Indices)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s = %s", VariableName(inst), result), nil

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
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s %s %s", VariableName(inst), x, op, y), nil

	case *ir.InstInsertElement:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating initial vector (%v): %v", inst.X, err)
		}
		elem, err := FormatValue(inst.Elem)
		if err != nil {
			return "", fmt.Errorf("error translating new element (%v): %v", inst.Elem, err)
		}
		index, err := FormatValue(inst.Index)
		if err != nil {
			return "", fmt.Errorf("error translating index (%v): %v", inst.Index, err)
		}
		if _, ok := inst.X.(*constant.Undef); ok {
			return fmt.Sprintf("%s[%s] = %s", VariableName(inst), index, elem), nil
		}
		return fmt.Sprintf("%s = %s; %s[%s] = %s", VariableName(inst), x, VariableName(inst), index, elem), nil

	case *ir.InstIntToPtr:
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		return fmt.Sprintf("%s = (%s)(unsafe.Pointer(uintptr(%s)))", VariableName(inst), to, from), nil

	case *ir.InstLoad:
		src, err := FormatValue(inst.Src)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.Src, err)
		}
		if strings.HasPrefix(src, "&") {
			return fmt.Sprintf("%s = %s", VariableName(inst), strings.TrimPrefix(src, "&")), nil
		}
		return fmt.Sprintf("%s = *%s", VariableName(inst), src), nil

	case *ir.InstLShr:
		x, err := FormatUnsigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		if t, ok := inst.Typ.(*types.IntType); ok && t.BitSize > 8 {
			return fmt.Sprintf("%s = int%d(%s >> %s)", VariableName(inst), t.BitSize, x, y), nil
		}
		return fmt.Sprintf("%s = %s >> %s", VariableName(inst), x, y), nil

	case *ir.InstMul:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s * %s", VariableName(inst), x, y), nil

	case *ir.InstOr:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v | %s[i] }", x, VariableName(inst), y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return fmt.Sprintf("%s = %s || %s", VariableName(inst), x, y), nil
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

	case *ir.InstSDiv:
		x, err := FormatSigned(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatSigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 8 {
			return fmt.Sprintf("%s = byte(%s / %s)", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s / %s", VariableName(inst), x, y), nil

	case *ir.InstSelect:
		cond, err := FormatValue(inst.Cond)
		if err != nil {
			return "", fmt.Errorf("error translating condition (%v): %v", inst.Cond, err)
		}
		valueTrue, err := FormatValue(inst.ValueTrue)
		if err != nil {
			return "", fmt.Errorf("error translating first operand (%v): %v", inst.ValueTrue, err)
		}
		valueFalse, err := FormatValue(inst.ValueFalse)
		if err != nil {
			return "", fmt.Errorf("error translating second operand (%v): %v", inst.ValueFalse, err)
		}
		name := VariableName(inst)
		return fmt.Sprintf("if %s { %s = %s } else { %s = %s }", cond, name, valueTrue, name, valueFalse), nil

	case *ir.InstSExt:
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return "", fmt.Errorf("unsupported To type for zext: %T", inst.To)
		}
		from, err := FormatSigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		return fmt.Sprintf("%s = int%d(%s)", VariableName(inst), toType.BitSize, from), nil

	case *ir.InstShl:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatUnsigned(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s << %s", VariableName(inst), x, y), nil

	case *ir.InstShuffleVector:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		mask, err := FormatValue(inst.Mask)
		if err != nil {
			return "", fmt.Errorf("error translating mask (%v): %v", inst.Mask, err)
		}
		length := inst.Typ.Len
		return fmt.Sprintf("for i, m := range %s { if m < %d { %s[i] = %s[m] } else { %s[i] = %s[m - %d] } }", mask, length, VariableName(inst), x, VariableName(inst), y, length), nil

	case *ir.InstSIToFP:
		from, err := FormatSigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstStore:
		dest, err := FormatValue(inst.Dst)
		if err != nil {
			return "", fmt.Errorf("error translating destination (%v): %v", inst.Dst, err)
		}
		src, err := FormatValue(inst.Src)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.Src, err)
		}
		if strings.HasPrefix(dest, "&") {
			return fmt.Sprintf("%s = %s", strings.TrimPrefix(dest, "&"), src), nil
		}
		return fmt.Sprintf("*%s = %s", dest, src), nil

	case *ir.InstSub:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		return fmt.Sprintf("%s = %s - %s", VariableName(inst), x, y), nil

	case *ir.InstTrunc:
		if vt, ok := inst.To.(*types.VectorType); ok {
			toType, ok := vt.ElemType.(*types.IntType)
			if !ok {
				return "", fmt.Errorf("unsupported To type for zext: %v", inst.To)
			}
			to, err := TypeSpec(toType)
			if err != nil {
				return "", fmt.Errorf("error translating To type (%v): %v", toType, err)
			}
			from, err := FormatValue(inst.From)
			if err != nil {
				return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
			}
			return fmt.Sprintf("for i, v := range %s { %s[i] = %s(v) }", from, VariableName(inst), to), nil
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating To type (%v): %v", inst.To, err)
		}
		from, err := FormatValue(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		if intType, ok := inst.To.(*types.IntType); ok && intType.BitSize < 8 {
			return fmt.Sprintf("%s = byte(%s & %d)", VariableName(inst), from, 255>>(8-intType.BitSize)), nil
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstUIToFP:
		from, err := FormatUnsigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		to, err := TypeSpec(inst.To)
		if err != nil {
			return "", fmt.Errorf("error translating type (%v): %v", inst.To, err)
		}
		return fmt.Sprintf("%s = %s(%s)", VariableName(inst), to, from), nil

	case *ir.InstXor:
		x, err := FormatValue(inst.X)
		if err != nil {
			return "", fmt.Errorf("error translating left operand (%v): %v", inst.X, err)
		}
		y, err := FormatValue(inst.Y)
		if err != nil {
			return "", fmt.Errorf("error translating right operand (%v): %v", inst.Y, err)
		}
		if _, ok := inst.Typ.(*types.VectorType); ok {
			return fmt.Sprintf("for i, v := range %s { %s[i] = v ^ %s[i] }", x, VariableName(inst), y), nil
		}
		if intType, ok := inst.Typ.(*types.IntType); ok && intType.BitSize == 1 {
			return fmt.Sprintf("%s = %s != %s", VariableName(inst), x, y), nil
		}
		return fmt.Sprintf("%s = %s ^ %s", VariableName(inst), x, y), nil

	case *ir.InstZExt:
		if vt, ok := inst.To.(*types.VectorType); ok {
			toType, ok := vt.ElemType.(*types.IntType)
			if !ok {
				return "", fmt.Errorf("unsupported To type for zext: %v", inst.To)
			}
			ft, ok := inst.From.Type().(*types.VectorType)
			if !ok {
				return "", fmt.Errorf("mismatched types for zext: %v and %v", inst.To, inst.From.Type())
			}
			fromType, ok := ft.ElemType.(*types.IntType)
			if !ok {
				return "", fmt.Errorf("unsupported From type for zext: %v", inst.From.Type())
			}
			from, err := FormatValue(inst.From)
			if err != nil {
				return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
			}
			return fmt.Sprintf("for i, v := range %s { %s[i] = int%d(uint%d(uint%d(v))) }", from, VariableName(inst), toType.BitSize, toType.BitSize, fromType.BitSize), nil
		}
		toType, ok := inst.To.(*types.IntType)
		if !ok {
			return "", fmt.Errorf("unsupported To type for zext: %T", inst.To)
		}
		from, err := FormatUnsigned(inst.From)
		if err != nil {
			return "", fmt.Errorf("error translating source (%v): %v", inst.From, err)
		}
		if fromType, ok := inst.From.Type().(*types.IntType); ok && fromType.BitSize == 1 {
			return fmt.Sprintf("if %s { %s = 1 } else { %s = 0 }", from, VariableName(inst), VariableName(inst)), nil
		}
		return fmt.Sprintf("%s = int%d(uint%d(%s))", VariableName(inst), toType.BitSize, toType.BitSize, from), nil

	default:
		return "", fmt.Errorf("unsupported instruction type: %T", inst)
	}
}

var libraryFunctions = map[string]string{
	"calloc":           "libc.Calloc",
	"fabs":             "math.Abs",
	"free":             "libc.Free",
	"leaven_va_arg":    "libc.VAArg",
	"llvm_fabs_f64":    "math.Abs",
	"llvm_fabs_f80":    "math.Abs",
	"llvm_pow_f64":     "math.Pow",
	"malloc":           "libc.Malloc",
	"__memcpy_chk":     "libc.MemcpyChk",
	"memset_pattern16": "libc.MemsetPattern16",
	"__memset_chk":     "libc.MemsetChk",
	"printf":           "noarch.Printf",
	"__strcat_chk":     "libc.StrcatChk",
	"strcmp":           "libc.Strcmp",
}
