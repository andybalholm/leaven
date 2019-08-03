package libc

import "unsafe"

// The functions in this file are transpiled from
// the Public Domain C Library (PDCLib).

func Memcmp(s1 *byte, s2 *byte, n int64) int32 {
	var p2_016, p1_015, incdec_ptr, incdec_ptr5 *byte
	var tobool14, cmp, tobool bool
	var v0, v1 byte
	var conv, conv1, sub, retval_0 int32
	var dec17_in, dec17 int64

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = tobool14, dec17_in, p2_016, p1_015, dec17, v0, v1, cmp, conv, conv1, sub, incdec_ptr, incdec_ptr5, tobool, retval_0

	tobool14 = n == 0
	if tobool14 {
		retval_0 = 0
		goto cleanup
	} else {
		dec17_in, p2_016, p1_015 = n, s2, s1
		goto while_body
	}

while_body:
	dec17 = dec17_in - 1
	v0 = *p1_015
	v1 = *p2_016
	cmp = v0 == v1
	if cmp {
		goto if_end
	} else {
		goto if_then
	}

if_then:
	conv = int32(uint32(v0))
	conv1 = int32(uint32(v1))
	sub = conv - conv1
	retval_0 = sub
	goto cleanup

if_end:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p1_015)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	incdec_ptr5 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p2_016)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	tobool = dec17 == 0
	if tobool {
		retval_0 = 0
		goto cleanup
	} else {
		dec17_in, p2_016, p1_015 = dec17, incdec_ptr5, incdec_ptr
		goto while_body
	}

cleanup:
	return retval_0
}

func Strchr(s *byte, c int32) *byte {
	var s_addr_0, incdec_ptr, retval_0 *byte
	var cmp, tobool bool
	var v0 byte
	var sext, conv2, conv int32

	_, _, _, _, _, _, _, _, _ = sext, conv2, s_addr_0, v0, conv, cmp, incdec_ptr, tobool, retval_0

	sext = c << 24
	conv2 = sext >> 24
	s_addr_0 = s
	goto do_body

do_body:
	v0 = *s_addr_0
	conv = int32(int8(v0))
	cmp = conv2 == conv
	if cmp {
		retval_0 = s_addr_0
		goto _return
	} else {
		goto do_cond
	}

do_cond:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s_addr_0)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	tobool = v0 == 0
	if tobool {
		retval_0 = nil
		goto _return
	} else {
		s_addr_0 = incdec_ptr
		goto do_body
	}

_return:
	return retval_0
}

func Strcmp(s1 *byte, s2 *byte) int32 {
	var s2_addr_013, s1_addr_012, incdec_ptr, incdec_ptr4, s2_addr_0_lcssa *byte
	var tobool11, cmp, tobool bool
	var v0, v1, v2, v3, _lcssa, v4 byte
	var conv5, conv6, sub int32

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v0, tobool11, v1, s2_addr_013, s1_addr_012, v2, cmp, incdec_ptr, incdec_ptr4, v3, tobool, s2_addr_0_lcssa, _lcssa, conv5, v4, conv6, sub

	v0 = *s1
	tobool11 = v0 == 0
	if tobool11 {
		s2_addr_0_lcssa, _lcssa = s2, 0
		goto while_end
	} else {
		v1, s2_addr_013, s1_addr_012 = v0, s2, s1
		goto land_rhs
	}

land_rhs:
	v2 = *s2_addr_013
	cmp = v1 == v2
	if cmp {
		goto while_body
	} else {
		s2_addr_0_lcssa, _lcssa = s2_addr_013, v1
		goto while_end
	}

while_body:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1_addr_012)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	incdec_ptr4 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s2_addr_013)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	v3 = *incdec_ptr
	tobool = v3 == 0
	if tobool {
		s2_addr_0_lcssa, _lcssa = incdec_ptr4, 0
		goto while_end
	} else {
		v1, s2_addr_013, s1_addr_012 = v3, incdec_ptr4, incdec_ptr
		goto land_rhs
	}

while_end:
	conv5 = int32(uint32(_lcssa))
	v4 = *s2_addr_0_lcssa
	conv6 = int32(uint32(v4))
	sub = conv5 - conv6
	return sub
}

func Strcpy(s1 *byte, s2 *byte) *byte {
	var s2_addr_0, s1_addr_0, incdec_ptr, incdec_ptr1 *byte
	var tobool bool
	var v0 byte

	_, _, _, _, _, _ = s2_addr_0, s1_addr_0, incdec_ptr, v0, incdec_ptr1, tobool

	s2_addr_0, s1_addr_0 = s2, s1
	goto while_cond

while_cond:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s2_addr_0)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	v0 = *s2_addr_0
	incdec_ptr1 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1_addr_0)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	*s1_addr_0 = v0
	tobool = v0 == 0
	if tobool {
		goto while_end
	} else {
		s2_addr_0, s1_addr_0 = incdec_ptr, incdec_ptr1
		goto while_cond
	}

while_end:
	return s1
}

func Strcspn(s1 *byte, s2 *byte) int64 {
	var p_0, incdec_ptr, arrayidx *byte
	var tobool20, tobool2, cmp, tobool bool
	var v0, v1, v2, v3 byte
	var len_021, inc, len_019 int64

	_, _, _, _, _, _, _, _, _, _, _, _, _, _ = v0, tobool20, v1, len_021, p_0, v2, tobool2, incdec_ptr, cmp, inc, arrayidx, v3, tobool, len_019

	v0 = *s1
	tobool20 = v0 == 0
	if tobool20 {
		len_019 = 0
		goto cleanup
	} else {
		v1, len_021 = v0, 0
		goto while_cond1_preheader
	}

while_cond1_preheader:
	p_0 = s2
	goto while_cond1

while_cond1:
	v2 = *p_0
	tobool2 = v2 == 0
	if tobool2 {
		goto while_end
	} else {
		goto while_body3
	}

while_body3:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p_0)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	cmp = v1 == v2
	if cmp {
		len_019 = len_021
		goto cleanup
	} else {
		p_0 = incdec_ptr
		goto while_cond1
	}

while_end:
	inc = len_021 + 1
	arrayidx = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1)) + uintptr(int64(inc))*unsafe.Sizeof(*(*byte)(nil))))
	v3 = *arrayidx
	tobool = v3 == 0
	if tobool {
		len_019 = inc
		goto cleanup
	} else {
		v1, len_021 = v3, inc
		goto while_cond1_preheader
	}

cleanup:
	return len_019
}

func Strlen(s *byte) int64 {
	var arrayidx *byte
	var tobool bool
	var v0 byte
	var rc_0, inc int64

	_, _, _, _, _ = rc_0, arrayidx, v0, tobool, inc

	rc_0 = 0
	goto while_cond

while_cond:
	arrayidx = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(int64(rc_0))*unsafe.Sizeof(*(*byte)(nil))))
	v0 = *arrayidx
	tobool = v0 == 0
	inc = rc_0 + 1
	if tobool {
		goto while_end
	} else {
		rc_0 = inc
		goto while_cond
	}

while_end:
	return rc_0
}

func Strncat(s1 *byte, s2 *byte, n int64) *byte {
	var s1_addr_0, incdec_ptr, s1_addr_121, s2_addr_019, incdec_ptr4, incdec_ptr3, s1_addr_1_lcssa *byte
	var tobool, tobool218, tobool5, tobool2 bool
	var v0, v1 byte
	var n_addr_020, dec int64

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = s1_addr_0, v0, tobool, incdec_ptr, tobool218, s1_addr_121, n_addr_020, s2_addr_019, v1, tobool5, incdec_ptr4, incdec_ptr3, dec, tobool2, s1_addr_1_lcssa

	s1_addr_0 = s1
	goto while_cond

while_cond:
	v0 = *s1_addr_0
	tobool = v0 == 0
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1_addr_0)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	if tobool {
		goto while_cond1_preheader
	} else {
		s1_addr_0 = incdec_ptr
		goto while_cond
	}

while_cond1_preheader:
	tobool218 = n == 0
	if tobool218 {
		s1_addr_1_lcssa = s1_addr_0
		goto if_then
	} else {
		s1_addr_121, n_addr_020, s2_addr_019 = s1_addr_0, n, s2
		goto land_end
	}

land_end:
	v1 = *s2_addr_019
	*s1_addr_121 = v1
	tobool5 = v1 == 0
	if tobool5 {
		goto if_end
	} else {
		goto while_body6
	}

while_body6:
	incdec_ptr4 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1_addr_121)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	incdec_ptr3 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s2_addr_019)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	dec = n_addr_020 - 1
	tobool2 = dec == 0
	if tobool2 {
		s1_addr_1_lcssa = incdec_ptr4
		goto if_then
	} else {
		s1_addr_121, n_addr_020, s2_addr_019 = incdec_ptr4, dec, incdec_ptr3
		goto land_end
	}

if_then:
	*s1_addr_1_lcssa = 0
	goto if_end

if_end:
	return s1
}

func Strncmp(s1 *byte, s2 *byte, n int64) int32 {
	var s2_addr_021, s1_addr_020, incdec_ptr, incdec_ptr5 *byte
	var cond19, tobool1, cmp, or_cond, cond bool
	var v0, _pre byte
	var conv8, conv9, sub, retval_0 int32
	var n_addr_022, dec int64

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = cond19, n_addr_022, s2_addr_021, s1_addr_020, v0, tobool1, _pre, cmp, or_cond, incdec_ptr, incdec_ptr5, dec, cond, conv8, conv9, sub, retval_0

	cond19 = n == 0
	if cond19 {
		retval_0 = 0
		goto _return
	} else {
		n_addr_022, s2_addr_021, s1_addr_020 = n, s2, s1
		goto land_lhs_true
	}

land_lhs_true:
	v0 = *s1_addr_020
	tobool1 = v0 != 0
	_pre = *s2_addr_021
	cmp = v0 == _pre
	or_cond = tobool1 && cmp
	if or_cond {
		goto while_body
	} else {
		goto if_else
	}

while_body:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1_addr_020)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	incdec_ptr5 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s2_addr_021)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	dec = n_addr_022 - 1
	cond = dec == 0
	if cond {
		retval_0 = 0
		goto _return
	} else {
		n_addr_022, s2_addr_021, s1_addr_020 = dec, incdec_ptr5, incdec_ptr
		goto land_lhs_true
	}

if_else:
	conv8 = int32(uint32(v0))
	conv9 = int32(uint32(_pre))
	sub = conv8 - conv9
	retval_0 = sub
	goto _return

_return:
	return retval_0
}

func Strncpy(s1 *byte, s2 *byte, n int64) *byte {
	var s1_addr_020, s2_addr_018, incdec_ptr1, incdec_ptr *byte
	var tobool17, tobool2, tobool, cmp14 bool
	var v0 byte
	var n_addr_019, dec, v1 int64

	_, _, _, _, _, _, _, _, _, _, _, _ = tobool17, s1_addr_020, n_addr_019, s2_addr_018, v0, incdec_ptr1, tobool2, incdec_ptr, dec, tobool, cmp14, v1

	tobool17 = n == 0
	if tobool17 {
		goto while_end8
	} else {
		s1_addr_020, n_addr_019, s2_addr_018 = s1, n, s2
		goto land_rhs
	}

land_rhs:
	v0 = *s2_addr_018
	incdec_ptr1 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1_addr_020)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	*s1_addr_020 = v0
	tobool2 = v0 == 0
	if tobool2 {
		goto while_end
	} else {
		goto while_body
	}

while_body:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s2_addr_018)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	dec = n_addr_019 - 1
	tobool = dec == 0
	if tobool {
		goto while_end8
	} else {
		s1_addr_020, n_addr_019, s2_addr_018 = incdec_ptr1, dec, incdec_ptr
		goto land_rhs
	}

while_end:
	cmp14 = uint64(n_addr_019) > 1
	if cmp14 {
		goto while_body6_preheader
	} else {
		goto while_end8
	}

while_body6_preheader:
	v1 = n_addr_019 - 1
	Memset(incdec_ptr1, 0, v1)
	goto while_end8

while_end8:
	return s1
}

func Strpbrk(s1 *byte, s2 *byte) *byte {
	var p1_017, p2_0, incdec_ptr, incdec_ptr6, retval_0 *byte
	var tobool16, tobool2, cmp, tobool bool
	var v0, v1, v2, v3 byte

	_, _, _, _, _, _, _, _, _, _, _, _, _ = v0, tobool16, v1, p1_017, p2_0, v2, tobool2, incdec_ptr, cmp, incdec_ptr6, v3, tobool, retval_0

	v0 = *s1
	tobool16 = v0 == 0
	if tobool16 {
		retval_0 = nil
		goto cleanup
	} else {
		v1, p1_017 = v0, s1
		goto while_cond1_preheader
	}

while_cond1_preheader:
	p2_0 = s2
	goto while_cond1

while_cond1:
	v2 = *p2_0
	tobool2 = v2 == 0
	if tobool2 {
		goto while_end
	} else {
		goto while_body3
	}

while_body3:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p2_0)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	cmp = v1 == v2
	if cmp {
		retval_0 = p1_017
		goto cleanup
	} else {
		p2_0 = incdec_ptr
		goto while_cond1
	}

while_end:
	incdec_ptr6 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p1_017)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	v3 = *incdec_ptr6
	tobool = v3 == 0
	if tobool {
		retval_0 = nil
		goto cleanup
	} else {
		v1, p1_017 = v3, incdec_ptr6
		goto while_cond1_preheader
	}

cleanup:
	return retval_0
}

func Strrchr(s *byte, c int32) *byte {
	var arrayidx, arrayidx1, arrayidx1_le, retval_0 *byte
	var tobool, cmp, tobool5 bool
	var v0, v1 byte
	var sext, conv3, conv int32
	var i_0, inc, i_1, dec int64

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = i_0, inc, arrayidx, v0, tobool, sext, conv3, i_1, dec, arrayidx1, v1, conv, cmp, tobool5, arrayidx1_le, retval_0

	i_0 = 0
	goto while_cond

while_cond:
	inc = i_0 + 1
	arrayidx = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(int64(i_0))*unsafe.Sizeof(*(*byte)(nil))))
	v0 = *arrayidx
	tobool = v0 == 0
	if tobool {
		goto do_body_preheader
	} else {
		i_0 = inc
		goto while_cond
	}

do_body_preheader:
	sext = c << 24
	conv3 = sext >> 24
	i_1 = inc
	goto do_body

do_body:
	dec = i_1 - 1
	arrayidx1 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(int64(dec))*unsafe.Sizeof(*(*byte)(nil))))
	v1 = *arrayidx1
	conv = int32(int8(v1))
	cmp = conv3 == conv
	if cmp {
		goto cleanup_split_loop_exit15
	} else {
		goto do_cond
	}

do_cond:
	tobool5 = dec == 0
	if tobool5 {
		retval_0 = nil
		goto cleanup
	} else {
		i_1 = dec
		goto do_body
	}

cleanup_split_loop_exit15:
	arrayidx1_le = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + uintptr(int64(dec))*unsafe.Sizeof(*(*byte)(nil))))
	retval_0 = arrayidx1_le
	goto cleanup

cleanup:
	return retval_0
}

func Strspn(s1 *byte, s2 *byte) int64 {
	var p_025, incdec_ptr, arrayidx *byte
	var tobool27, tobool224, tobool2, cmp, tobool bool
	var v0, v1, v2, v3, v4, v5 byte
	var len_028, inc, len_023 int64

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v0, tobool27, v1, tobool224, v2, len_028, v3, tobool2, v4, p_025, cmp, incdec_ptr, inc, arrayidx, v5, tobool, len_023

	v0 = *s1
	tobool27 = v0 == 0
	if tobool27 {
		len_023 = 0
		goto cleanup
	} else {
		goto while_cond1_preheader_lr_ph
	}

while_cond1_preheader_lr_ph:
	v1 = *s2
	tobool224 = v1 == 0
	v2, len_028 = v0, 0
	goto while_cond1_preheader

while_cond1_preheader:
	if tobool224 {
		len_023 = 0
		goto cleanup
	} else {
		v4, p_025 = v1, s2
		goto while_body3
	}

while_cond1:
	v3 = *incdec_ptr
	tobool2 = v3 == 0
	if tobool2 {
		len_023 = len_028
		goto cleanup
	} else {
		v4, p_025 = v3, incdec_ptr
		goto while_body3
	}

while_body3:
	cmp = v2 == v4
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p_025)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	if cmp {
		goto if_end9
	} else {
		goto while_cond1
	}

if_end9:
	inc = len_028 + 1
	arrayidx = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1)) + uintptr(int64(inc))*unsafe.Sizeof(*(*byte)(nil))))
	v5 = *arrayidx
	tobool = v5 == 0
	if tobool {
		len_023 = inc
		goto cleanup
	} else {
		v2, len_028 = v5, inc
		goto while_cond1_preheader
	}

cleanup:
	return len_023
}

func Strstr(s1 *byte, s2 *byte) *byte {
	var s1_addr_028, p1_12437, p2_02536, incdec_ptr, incdec_ptr7, incdec_ptr9, retval_0 *byte
	var tobool27, tobool223, cmp35, tobool2, cmp, tobool bool
	var v0, v1, v2, v3, _pre, v4 byte

	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _ = v0, tobool27, v1, tobool223, v2, s1_addr_028, cmp35, p1_12437, p2_02536, incdec_ptr, incdec_ptr7, v3, tobool2, _pre, cmp, incdec_ptr9, v4, tobool, retval_0

	v0 = *s1
	tobool27 = v0 == 0
	if tobool27 {
		retval_0 = nil
		goto cleanup
	} else {
		goto while_cond1_preheader_lr_ph
	}

while_cond1_preheader_lr_ph:
	v1 = *s2
	tobool223 = v1 == 0
	v2, s1_addr_028 = v0, s1
	goto while_cond1_preheader

while_cond1_preheader:
	if tobool223 {
		retval_0 = s1
		goto cleanup
	} else {
		goto land_rhs_preheader
	}

land_rhs_preheader:
	cmp35 = v2 == v1
	if cmp35 {
		p1_12437, p2_02536 = s1_addr_028, s2
		goto while_body6
	} else {
		goto if_end
	}

while_body6:
	incdec_ptr = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p1_12437)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	incdec_ptr7 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p2_02536)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	v3 = *incdec_ptr7
	tobool2 = v3 == 0
	if tobool2 {
		retval_0 = s1_addr_028
		goto cleanup
	} else {
		goto while_body6_land_rhs_crit_edge
	}

while_body6_land_rhs_crit_edge:
	_pre = *incdec_ptr
	cmp = _pre == v3
	if cmp {
		p1_12437, p2_02536 = incdec_ptr, incdec_ptr7
		goto while_body6
	} else {
		goto if_end
	}

if_end:
	incdec_ptr9 = (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(s1_addr_028)) + 1*unsafe.Sizeof(*(*byte)(nil))))
	v4 = *incdec_ptr9
	tobool = v4 == 0
	if tobool {
		retval_0 = nil
		goto cleanup
	} else {
		v2, s1_addr_028 = v4, incdec_ptr9
		goto while_cond1_preheader
	}

cleanup:
	return retval_0
}
