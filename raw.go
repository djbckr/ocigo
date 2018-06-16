package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"
import (
	"unsafe"
	"encoding/hex"
	"runtime"
)

type Raw struct {
	data *C.OCIRaw
	dataptr unsafe.Pointer
}

func rawFinalizer(r *Raw)  {
	err := checkError(C.OCIRawResize(genv, gerr, 0, (**C.OCIRaw)(unsafe.Pointer(r.dataptr))), gerr)
	if err != nil {
		panic(err.Error())
	}
}

func MakeRaw() *Raw {
	var d *C.OCIRaw
	rslt := &Raw{}
	rslt.data = d
	rslt.dataptr = unsafe.Pointer(&d)
	runtime.SetFinalizer(rslt, rawFinalizer)
	return rslt
}

func MakeRawWithSize(size int) *Raw {

	var d *C.OCIRaw

	err := checkError(C.OCIRawResize(genv, gerr, C.uint(size), (**C.OCIRaw)(unsafe.Pointer(&d))), gerr)
	if err != nil {
		panic(err.Error())
	}

	rslt := &Raw{}
	rslt.data = d
	rslt.dataptr = unsafe.Pointer(&d)
	return rslt
}

func (r *Raw) Data() []byte {
	sz := int(C.OCIRawSize(genv, r.data))
	pt := unsafe.Pointer(C.OCIRawPtr(genv, r.data))
	return C.GoBytes(pt, C.int(sz))
}

func (r *Raw) String() string {
	return hex.EncodeToString(r.Data())
}