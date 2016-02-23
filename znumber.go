package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"
import (
	"errors"
	//f "fmt"
	"runtime"
	"unsafe"
)

/* An OCI representation of an Oracle Number type. This supports all of the
   features of Oracle numbers, as they can easily exceed the limitations
   of most float/int operations. */
type Number struct {
	err    *C.OCIError
	number C.OCINumber
}

func finalizeNumber(n *Number) {
	ociHandleFree((unsafe.Pointer)(n.err), htypeError)
}

func makeNumberInstance() (rslt *Number) {
	rslt = &Number{}
	ociHandleAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.err)), htypeError)
	runtime.SetFinalizer(rslt, finalizeNumber)
	return
}

// Convert from a native integer type to Oracle Number
func NumberFromInt(val interface{}) (*Number, error) {

	rslt := makeNumberInstance()

	var valtyp C.uword
	var val_u uint64
	var val_s int64
	var valptr unsafe.Pointer
	var valsz C.uword

	switch val.(type) {
	case uint8:
		valtyp = C.OCI_NUMBER_UNSIGNED
		val_u = uint64(val.(uint8))
		valptr = (unsafe.Pointer)(&val_u)
		valsz = (C.uword)(unsafe.Sizeof(val_u))
	case uint16:
		valtyp = C.OCI_NUMBER_UNSIGNED
		val_u = uint64(val.(uint16))
		valptr = (unsafe.Pointer)(&val_u)
		valsz = (C.uword)(unsafe.Sizeof(val_u))
	case uint32:
		valtyp = C.OCI_NUMBER_UNSIGNED
		val_u = uint64(val.(uint32))
		valptr = (unsafe.Pointer)(&val_u)
		valsz = (C.uword)(unsafe.Sizeof(val_u))
	case uint64:
		valtyp = C.OCI_NUMBER_UNSIGNED
		val_u = val.(uint64)
		valptr = (unsafe.Pointer)(&val_u)
		valsz = (C.uword)(unsafe.Sizeof(val_u))
	case uint:
		valtyp = C.OCI_NUMBER_UNSIGNED
		val_u = uint64(val.(uint))
		valptr = (unsafe.Pointer)(&val_u)
		valsz = (C.uword)(unsafe.Sizeof(val_u))
	case int8:
		valtyp = C.OCI_NUMBER_SIGNED
		val_s = int64(val.(int8))
		valptr = (unsafe.Pointer)(&val_s)
		valsz = (C.uword)(unsafe.Sizeof(val_s))
	case int16:
		valtyp = C.OCI_NUMBER_SIGNED
		val_s = int64(val.(int16))
		valptr = (unsafe.Pointer)(&val_s)
		valsz = (C.uword)(unsafe.Sizeof(val_s))
	case int32:
		valtyp = C.OCI_NUMBER_SIGNED
		val_s = int64(val.(int32))
		valptr = (unsafe.Pointer)(&val_s)
		valsz = (C.uword)(unsafe.Sizeof(val_s))
	case int64:
		valtyp = C.OCI_NUMBER_SIGNED
		val_s = val.(int64)
		valptr = (unsafe.Pointer)(&val_s)
		valsz = (C.uword)(unsafe.Sizeof(val_s))
	case int:
		valtyp = C.OCI_NUMBER_SIGNED
		val_s = int64(val.(int))
		valptr = (unsafe.Pointer)(&val_s)
		valsz = (C.uword)(unsafe.Sizeof(val_s))
	default:
		return nil, errors.New("Invalid integer type for conversion")
	}

	vErr := checkError(
		C.OCINumberFromInt(
			rslt.err,
			valptr,
			valsz,
			valtyp,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)

}

// Convert from a native float type to Oracle Number
func NumberFromFloat(val interface{}) (*Number, error) {

	rslt := makeNumberInstance()

	var v64 float64

	switch val.(type) {
	case float32:
		v64 = float64(val.(float32))
	case float64:
		v64 = val.(float64)
	default:
		return nil, errors.New("Invalid float type for conversion")
	}

	vErr := checkError(
		C.OCINumberFromReal(
			rslt.err,
			unsafe.Pointer(&v64),
			(C.uword)(unsafe.Sizeof(v64)),
			&rslt.number), rslt.err)

	return rslt, processError(vErr)

}

// Convert a string to an Oracle Number using formatting/nls.
// val is the string to convert
// fmt is the format string. OCI is a bit brain-dead about this. You *must* provide a valid format string here.
// nls is the NLS parameter settings string. You can pass an empty string here to use the default settings.
func NumberFromStringFmt(val string, fmt string, nls string) (*Number, error) {

	rslt := makeNumberInstance()

	str := []byte(val)

	var format []byte = []byte(fmt)
	var fmtlen C.ub4 = (C.ub4)(len(format))
	var formatp *C.oratext
	if fmtlen > 0 {
		formatp = (*C.oratext)(unsafe.Pointer(&format[0]))
	}

	var nlsparams []byte = []byte(nls)
	var nlsparlen C.ub4 = (C.ub4)(len(nlsparams))
	var nlsparamsp *C.oratext
	if nlsparlen > 0 {
		nlsparamsp = (*C.oratext)(unsafe.Pointer(&nlsparams[0]))
	}

	vErr := checkError(
		C.OCINumberFromText(
			rslt.err,
			(*C.oratext)(unsafe.Pointer(&str[0])),
			(C.ub4)(len(val)),
			formatp,
			fmtlen,
			nlsparamsp,
			nlsparlen,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)

}

// Convert a basic numerical string to an Oracle Number
func NumberFromString(val string) (*Number, error) {
	lvl := []byte(val)
	fmt := make([]byte, len(lvl))
	for i := 0; i < len(fmt); i++ {
		if (lvl[i] >= 0x30) && (lvl[i] <= 0x39) {
			// any digits become the letter "9"
			fmt[i] = 0x39
		} else if (lvl[i] == 0x2B) || (lvl[i] == 0x2D) {
			// a "+" or "-" becomes an "S" (sign)
			fmt[i] = 0x53
		} else {
			// any other characters get passed through, such as "," "." "$"
			fmt[i] = lvl[i]
		}
	}
	return NumberFromStringFmt(val, string(fmt), "")
}

// convert an Oracle Number to native integer
func (num *Number) ToInt() (int64, error) {

	var rslt int64

	vErr := checkError(
		C.OCINumberToInt(
			num.err,
			&num.number,
			(C.uword)(unsafe.Sizeof(rslt)),
			C.OCI_NUMBER_SIGNED,
			unsafe.Pointer(&rslt)), num.err)

	return rslt, processError(vErr)

}

// convert an Oracle Number to native float
func (num *Number) ToFloat() (float64, error) {

	var rslt float64

	vErr := checkError(
		C.OCINumberToReal(
			num.err,
			&num.number,
			(C.uword)(unsafe.Sizeof(rslt)),
			unsafe.Pointer(&rslt)), num.err)

	return rslt, processError(vErr)

}

// convert an Oracle Number to string
// The fmt and nls parameters may be empty if desired.
// Otherwise provide a format string and/or an NLS parameter string
func (num *Number) ToString(fmt, nls string) (string, error) {

	var buflen C.ub4 = 64
	buf := make([]byte, buflen)

	var format []byte
	var fmtlen C.ub4

	if len(fmt) == 0 {
		format = []byte("TM")
		fmtlen = (C.ub4)(len(format))
	} else {
		format = []byte(fmt)
		fmtlen = (C.ub4)(len(format))
	}

	var nlsparams []byte
	var nlsparlen C.ub4 = 0
	var nlsparamsp *C.oratext

	if len(nls) > 0 {
		nlsparams = []byte(nls)
		nlsparlen = (C.ub4)(len(nlsparams))
		nlsparamsp = (*C.oratext)(&nlsparams[0])
	}

	vErr := checkError(
		C.OCINumberToText(
			num.err,
			&num.number,
			(*C.oratext)(&format[0]), fmtlen,
			nlsparamsp, nlsparlen,
			&buflen,
			(*C.oratext)(unsafe.Pointer(&buf[0]))), num.err)

	return nulTerminatedByteToString(buf), processError(vErr)

}

// convert an Oracle Number to a string using default settings
func (num *Number) String() string {
	rslt, err := num.ToString("", "")
	if err != nil {
		return err.Error()
	}
	return rslt
}

func (num *Number) Abs() (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberAbs(
			rslt.err,
			&num.number,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}

func (num *Number) Add(number *Number) (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberAdd(
			rslt.err,
			&num.number,
			&number.number,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}

func (num *Number) Cmp(number *Number) (int, error) {

	var rslt C.sword

	vErr := checkError(
		C.OCINumberCmp(
			num.err,
			&num.number,
			&number.number,
			(*C.sword)(unsafe.Pointer(&rslt))), num.err)

	return (int)(rslt), processError(vErr)
}

func (num *Number) Div(number *Number) (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberDiv(
			rslt.err,
			&num.number,
			&number.number,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}

func (num *Number) Mod(number *Number) (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberMod(
			rslt.err,
			&num.number,
			&number.number,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}

func (num *Number) Mul(number *Number) (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberMul(
			rslt.err,
			&num.number,
			&number.number,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}

func (num *Number) Round(decplaces int) (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberRound(
			rslt.err,
			&num.number,
			(C.sword)(decplaces),
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}

func (num *Number) Sub(number *Number) (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberSub(
			rslt.err,
			&num.number,
			&number.number,
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}

func (num *Number) Trunc(decplaces int) (*Number, error) {

	rslt := makeNumberInstance()

	vErr := checkError(
		C.OCINumberTrunc(
			rslt.err,
			&num.number,
			(C.sword)(decplaces),
			&rslt.number), rslt.err)

	return rslt, processError(vErr)
}
