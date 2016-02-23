package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"

import (
	//"crypto/sha256"
	//"encoding/hex"
	"errors"
	"fmt"
	"unsafe"
	//"time"
)

type tCharSemantics int8

const (
	charSemanticsByte tCharSemantics = 0
	charSemanticsChar tCharSemantics = 1
)

type Column struct {
	datatype      ociSqlType
	name          string
	sizeBytes     int16
	sizeChars     int16
	charSemantics tCharSemantics
	precision     int16
	scale         int8
	nullable      bool
	nTypeName     string
	nTypeSchema   string
}

type ResultSet struct {
	stmt    *Statement
	columns []*Column
}

func (stmt *Statement) query(count uint32) (*ResultSet, error) {

	if stmt.stmtype != StmtSelect {
		return nil, errors.New("Statement type must be a query")
	}

	err := stmt.exec(count, false)
	if err != nil {
		return nil, processError(err)
	}

	var parmcnt uint32

	err = ociAttrGet(unsafe.Pointer(stmt.stm), htypeStatement, unsafe.Pointer(&parmcnt), nil, attrParamCount, stmt.err)

	if err != nil {
		return nil, processError(err)
	}

	fmt.Println("param count: ", parmcnt)

	rslt := &ResultSet{stmt: stmt}
	rslt.columns = make([]*Column, parmcnt)

	// 8==1, 16==2, 32==4, 64==8

	var colIndx uint32 = 1
	var arrIndx uint32 = 0

	for colIndx <= parmcnt {

		paramp, err := stmt.getParameter(colIndx)

		if err != nil {
			return nil, processError(err)
		}

		rslt.columns[arrIndx], err = getColumnInfo(paramp, stmt.err)

		if err != nil {
			return nil, processError(err)
		}

		fmt.Println(rslt.columns[arrIndx])

		colIndx++
		arrIndx++
	}

	// TODO get result meta-data
	return nil, nil
}

func (stmt *Statement) getParameter(indx uint32) (rslt *C.OCIParam, err *OciError) {
	err = checkError(
		C.OCIParamGet(
			(unsafe.Pointer)(stmt.stm),
			(C.ub4)(htypeStatement),
			stmt.err,
			(*unsafe.Pointer)(unsafe.Pointer(&rslt)),
			(C.ub4)(indx)), stmt.err)
	return
}

func getColumnInfo(paramp *C.OCIParam, errhndl *C.OCIError) (rslt *Column, err *OciError) {
	rslt = &Column{}

	// data type (C.ub2)
	err = ociAttrGet(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		(unsafe.Pointer)(&rslt.datatype),
		nil, attrDataType, errhndl)

	if err != nil {
		return
	}

	// name (C.OraText *)
	rslt.name, err = ociAttrGetString(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		attrName,
		errhndl)

	if err != nil {
		return
	}

	// data size (C.ub2) // num bytes needed
	err = ociAttrGet(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		(unsafe.Pointer)(&rslt.sizeBytes),
		nil, attrDataSize, errhndl)

	if err != nil {
		return
	}

	// char_size (C.ub2) // num chars allowed
	err = ociAttrGet(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		(unsafe.Pointer)(&rslt.sizeChars),
		nil, attrCharSize, errhndl)

	if err != nil {
		return
	}

	// char_used (C.ub1) (1==char-semantics, 0==byte-semantics)
	var bb C.ub1

	err = ociAttrGet(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		(unsafe.Pointer)(&bb),
		nil, attrCharUsed, errhndl)

	if err != nil {
		return
	}

	if bb == 0 {
		rslt.charSemantics = charSemanticsByte
	} else {
		rslt.charSemantics = charSemanticsChar
	}

	// precision (C.sb2) // The precision of numeric columns. If the precision is nonzero and scale is -127, then it is a FLOAT; otherwise, it is a NUMBER(precision, scale). When precision is 0, NUMBER(precision, scale) can be represented simply as NUMBER.
	err = ociAttrGet(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		(unsafe.Pointer)(&rslt.precision),
		nil, attrPrecision, errhndl)

	if err != nil {
		return
	}

	// scale (C.sb1)
	err = ociAttrGet(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		(unsafe.Pointer)(&rslt.scale),
		nil, attrScale, errhndl)

	if err != nil {
		return
	}

	// nullable (C.sb1) // if == 0, not nullable ; <> 0 nullable
	var nn C.sb1

	err = ociAttrGet(
		(unsafe.Pointer)(paramp),
		(ociHandleType)(dtypeParam),
		(unsafe.Pointer)(&nn),
		nil, attrIsNull, errhndl)

	if err != nil {
		return
	}

	if nn == 0 {
		rslt.nullable = false
	} else {
		rslt.nullable = true
	}

	if rslt.datatype == sqltObject {

		rslt.nTypeName, err = ociAttrGetString(
			(unsafe.Pointer)(paramp),
			(ociHandleType)(dtypeParam),
			attrTypeName, errhndl)

		if err != nil {
			return
		}
		/*
			rslt.nTypeSchema, err = ociAttrGetString(
				(unsafe.Pointer)(paramp),
				(ociHandleType)(dtypeParam),
				attrTypeSchema, errhndl)
		*/
	}

	return

}
