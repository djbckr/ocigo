package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"

import (
	"errors"
	"fmt"
	//"runtime"
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

func ociGetError(errh *C.OCIError) (string, int32) {

	BUFSIZE := 1024

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
		callresult := C.OCIErrorGet(unsafe.Pointer(errh), (C.ub4)(indx), nil, (*C.sb4)(pp), (*C.OraText)(unsafe.Pointer(&buffer[0])), (C.ub4)(BUFSIZE), C.OCI_HTYPE_ERROR)
		switch callresult {
		case C.OCI_SUCCESS:
			rslts = append(rslts, nulTerminatedByteToString(buffer))
			for i := 0; i < BUFSIZE; i++ {
				buffer[i] = 0
			}
		case C.OCI_NO_DATA:
			break MyLoop
		default:
			rslts = append(rslts, fmt.Sprintf("Error retrieving error: code %v", callresult))
			break MyLoop
		}
	}

	return strings.Join(rslts, "\n"), p1
}
