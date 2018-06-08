package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"

/*
   The OCI error management is a bit of a headache. It's somewhat rare to get warnings,
   but they happen - such as password expiration warnings on login, or PL/SQL compilation
   errors. These generally don't require special handling, but should be checked by
   the client.

   The programs using this library should not see the OciError struct, but should instead
   see error objects. This is used for this library processing purposes.
   -- perhaps OciError should be ociError? --
*/

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

type OciError struct {
	code int32
	err  error
	inf  string
}

func (err *OciError) IsWarning() bool {
	return err.inf != ""
}

func (err *OciError) IsError() bool {
	return err.err != nil
}

func (err *OciError) Warning() string {
	return err.inf
}

func (err *OciError) Error() string {
	return err.err.Error()
}

func processError(err *OciError) error {
	if err == nil {
		return nil
	}

	if err.IsWarning() {
		fmt.Println(err.Warning())
	}

	if err.IsError() {
		return err.err
	}

	return nil
}

// this function is called on nearly every OCI call
func checkError(errval C.sword, errhndl *C.OCIError) (result *OciError) {

	if errval == C.OCI_SUCCESS {
		result = nil
		return
	}

	if errval == C.OCI_INVALID_HANDLE {
		result = &OciError{err: errors.New("Invalid Handle!")}
		return
	}

	rsltstr, eCode := ociGetError(errhndl)

	result = &OciError{code: eCode}

	if errval == C.OCI_SUCCESS_WITH_INFO {
		result.inf = rsltstr
	} else {
		result.err = errors.New(rsltstr)
	}

	return
}

// low-level call into OCI; Oracle allows multiple errors
// to be reported in one error handle, hence the loop
func ociGetError(errh *C.OCIError) (string, int32) {

	BUFSIZE := 1024 // 1kb ought to be enough for anybody

	rslts := make([]string, 0, 10)
	buffer := make([]byte, BUFSIZE)
	indx := 0
	var p1 int32 // the first error code
	var p2 int32 // each subsequent error code (discarded)
	var pp *int32

MyLoop:
	for {
		indx++
		if indx == 1 {
			pp = &p1
		} else {
			pp = &p2
		}

		callResult := C.OCIErrorGet(
			unsafe.Pointer(errh),
			(C.ub4)(indx),
			nil,
			(*C.sb4)(pp),
			(*C.OraText)(unsafe.Pointer(&buffer[0])),
			(C.ub4)(BUFSIZE), C.OCI_HTYPE_ERROR)

		switch callResult {
		case C.OCI_SUCCESS:
			rslts = append(rslts, nulTerminatedByteToString(buffer))
			for i := 0; i < BUFSIZE; i++ {
				buffer[i] = 0
			}
		case C.OCI_NO_DATA:
			break MyLoop
		default: // this should *never* happen!
			rslts = append(rslts, fmt.Sprintf("Error retrieving error: code %v", callResult))
			break MyLoop
		}
	}

	return strings.Join(rslts, "\n"), p1
}
