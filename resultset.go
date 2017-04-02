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
	sizeBytes     uint16
	sizeChars     uint16
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

	var paramCount uint32

	paramCount, err = ociAttrGetUB4(unsafe.Pointer(stmt.stm), htypeStatement, attrParamCount, stmt.err)

	if err != nil {
		return nil, processError(err)
	}

	rslt := &ResultSet{stmt: stmt}
	rslt.columns = make([]*Column, paramCount)

	// 8==1, 16==2, 32==4, 64==8

	var colIndx uint32 = 1
	var arrIndx uint32

	for colIndx <= paramCount {

		paramPtr, err := stmt.getParameter(colIndx)

		if err != nil {
			return nil, processError(err)
		}

		rslt.columns[arrIndx], err = getColumnInfo(paramPtr, stmt.err)

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

func getColumnInfo(paramPtr *C.OCIParam, errhndl *C.OCIError) (rslt *Column, err *OciError) {
	rslt = &Column{}

	// data type (C.ub2)
	var ub2 uint16
	ub2, err = ociAttrGetUB2(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrDataType, errhndl)

	if err != nil {
		return
	}
	rslt.datatype = ociSqlType(ub2)

	// name (C.OraText *)
	rslt.name, err = ociAttrGetString(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrName,
		errhndl)

	if err != nil {
		return
	}

	// data size (C.ub2) // num bytes needed
	ub2, err = ociAttrGetUB2(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrDataSize, errhndl)

	if err != nil {
		return
	}

	rslt.sizeBytes = ub2

	// char_size (C.ub2) // num chars allowed
	ub2, err = ociAttrGetUB2(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrCharSize, errhndl)

	if err != nil {
		return
	}

	rslt.sizeChars = ub2

	// char_used (C.ub1) (1==char-semantics, 0==byte-semantics)
	var bb uint8

	bb, err = ociAttrGetUB1(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrCharUsed, errhndl)

	if err != nil {
		return
	}

	if bb == 0 {
		rslt.charSemantics = charSemanticsByte
	} else {
		rslt.charSemantics = charSemanticsChar
	}

	// precision (C.sb2) // The precision of numeric columns. If the precision is nonzero and scale is -127, then it is a FLOAT; otherwise, it is a NUMBER(precision, scale).
	// When precision is 0, NUMBER(precision, scale) can be represented simply as NUMBER.
	var sb2 int16

	sb2, err = ociAttrGetSB2(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrPrecision, errhndl)

	if err != nil {
		return
	}

	rslt.precision = sb2

	var sb1 int8

	// scale (C.sb1)
	sb1, err = ociAttrGetSB1(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrScale, errhndl)

	if err != nil {
		return
	}

	rslt.scale = sb1

	// nullable (C.sb1) // if == 0, not nullable ; <> 0 nullable
	var nn int8

	nn, err = ociAttrGetSB1(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrIsNull, errhndl)

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
			(unsafe.Pointer)(paramPtr),
			(ociHandleType)(dtypeParam),
			attrTypeName, errhndl)

		if err != nil {
			return
		}

		rslt.nTypeSchema, err = ociAttrGetString(
			(unsafe.Pointer)(paramPtr),
			(ociHandleType)(dtypeParam),
			attrTypeSchema, errhndl)

	}

	return

}
