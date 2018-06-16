package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"

import (
	// "crypto/sha256"
	// "encoding/hex"
	"errors"
	"unsafe"
	// "time"
	"fmt"
)

type tCharSemantics int8

const (
	charSemanticsByte tCharSemantics = 0
	charSemanticsChar tCharSemantics = 1
)

type Column struct {
	datatype      ociSqlType
	name          string
	sizeBytes     int32
	sizeChars     uint16
	charSemantics tCharSemantics
	precision     int16
	scale         int8
	nullable      bool
	nTypeName     string
	nTypeSchema   string
	defnptr       *C.OCIDefine
	buffer        interface{}
	ind           *int16
}

func charSemantics(c tCharSemantics) string {
	switch c {
	case charSemanticsByte:
		return "B"
	case charSemanticsChar:
		return "C"
	default:
		return "-"
	}
}

type ResultSet struct {
	stmt    *Statement
	columns []*Column
}

func (rs *ResultSet) GetColumns() []*Column {
	return rs.columns
}

func (rs *ResultSet) Fetch() (rslt bool, err *OciError) {

	err = checkError(
		C.OCIStmtFetch2(
			rs.stmt.stm,
			rs.stmt.err,
			1, C.OCI_FETCH_NEXT, 0, C.OCI_DEFAULT), rs.stmt.err)

	rslt = err == nil

	if !rslt {
		if err.code == 1403 {
			err = nil
		}
	}

	return
}

func (col *Column) Print() string {
	return fmt.Sprintf("Name: %v ~ Type: %v ~ SizeBytes: %v ~ SizeChars: %v ~ Char/Byte: %v ~ Prec: %v ~ Scale: %v ~ Nullable: %v ~ ObjSchema: %v ~ ObjName: %v",
		col.name, SqlTypeName(col.datatype), col.sizeBytes, col.sizeChars, charSemantics(col.charSemantics), col.precision, col.scale, col.nullable, col.nTypeSchema, col.nTypeName)
}

func (col *Column) IsNull() bool {
	return col.ind != nil && *col.ind != 0
}

func (col *Column) IsNotNull() bool {
	return col.ind != nil && *col.ind == 0
}

func (col *Column) Get() interface{} {

	if col.ind != nil && *col.ind != 0 {
		return nil
	}

	switch v := col.buffer.(type) {
	case []byte:
		return nulTerminatedByteToString(v)
	case *Number:
		return v
	case *TimeStamp:
		return v
	case *float64:
		return *v
	case *float32:
		return *v
	case *Interval:
		return v
	case *Raw:
		return v
	default:
		return nil
	}
}

func (stmt *Statement) query(count uint32) (*ResultSet, error) {

	if stmt.stmtype != StmtSelect {
		return nil, errors.New("statement type must be a query")
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

		stmt.doDefine(rslt.columns[arrIndx], colIndx)

		// fmt.Println(rslt.columns[arrIndx])

		colIndx++
		arrIndx++
	}

	return rslt, nil
}

func (stmt *Statement) doDefine(column *Column, colIndx uint32) (err *OciError) {

	var sqlType ociSqlType
	var sizeBytes int32
	var buffer interface{}
	var bufptr unsafe.Pointer

	switch column.datatype {
	case sqltVarchar, sqltVarchar2, sqltChar :
		sizeBytes = column.sizeBytes + 1
		buf := make([]byte, sizeBytes)
		buffer = buf
		bufptr = unsafe.Pointer(&buf[0])
		sqlType = C.SQLT_STR

	case sqltNumber:
		sizeBytes = column.sizeBytes
		num := makeNumberInstance()
		buffer = num
		bufptr = unsafe.Pointer(&num.number)
		sqlType = C.SQLT_VNU

	case sqltBFloat:
		sizeBytes = 4
		sqlType = C.SQLT_BFLOAT
		var flt32 float32
		buffer = &flt32
		bufptr = unsafe.Pointer(&flt32)

	case sqltBDouble:
		sizeBytes = 8
		sqlType = C.SQLT_BDOUBLE
		var flt64 float64
		buffer = &flt64
		bufptr = unsafe.Pointer(&flt64)

	case sqltDate, sqltTimestamp, sqltTimestampTZ, sqltTimestampLTZ:
		var tstype TimestampType
		switch column.datatype {
		case sqltDate, sqltTimestamp:
			tstype = TypeTimestamp
			sqlType = sqltTimestamp
		case sqltTimestampTZ:
			tstype = TypeTimestampTZ
			sqlType = sqltTimestampTZ
		case sqltTimestampLTZ:
			tstype = TypeTimestampLTZ
			sqlType = sqltTimestampLTZ
		}
		sizeBytes = 0
		date := makeTimestampInstance(stmt.ses, tstype)
		buffer = date
		bufptr = date.ptrdt

	case sqltIntervalDS, sqltIntervalYM:
		var itype IntervalType
		switch column.datatype {
		case sqltIntervalDS:
			itype = TypeIntervalDS
			sqlType = sqltIntervalDS
		case sqltIntervalYM:
			itype = TypeIntervalYM
			sqlType = sqltIntervalYM
		}
		sizeBytes = 0
		interval := makeIntervalInstance(stmt.ses, itype)
		buffer = interval
		bufptr = interval.ptrintvl

	case sqltUnsigned8 /* aka RAW */:
		sizeBytes = column.sizeBytes

		raw := MakeRawWithSize(int(sizeBytes))

		buffer = raw
		bufptr = unsafe.Pointer(raw.data)

		sqlType = sqltRaw

	default:
		// do nothing for now
	}

	var pdefnptr *C.OCIDefine
	var pind int16

	if sqlType != 0 {
		fmt.Println("defining " + column.name)
		err = checkError(C.OCIDefineByPos(
			stmt.stm,
			&pdefnptr,
			stmt.err,
			C.ub4(colIndx),
			bufptr,
			C.sb4(sizeBytes),
			C.ub2(sqlType),
			unsafe.Pointer(&pind),
			nil, nil, C.OCI_DEFAULT), stmt.err)

		column.buffer = buffer
		column.defnptr = pdefnptr
		column.ind = &pind
		fmt.Println("defined column " + column.name)
	}

	return
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
	var sb4 int32
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
	sb4, err = ociAttrGetSB4(
		(unsafe.Pointer)(paramPtr),
		(ociHandleType)(dtypeParam),
		attrDataSize, errhndl)

	if err != nil {
		return
	}

	rslt.sizeBytes = sb4

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
			attrSchemaName, errhndl)

	}

	return

}
