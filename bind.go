package oci

/*
#cgo pkg-config: oci
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <oci.h>
*/
import "C"

/*
import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
)
*/
type Bind struct {
	bindhndl *C.OCIBind
}

type DType C.ub2

const (
	DTypeString        DType = C.SQLT_CHR           /* (ORANET TYPE) character string */
	DTypeSignedInt     DType = C.SQLT_INT           /* (ORANET TYPE) integer */
	DTypeUnsignedInt   DType = C.SQLT_UIN           /* unsigned integer */
	DTypeFloat         DType = C.SQLT_FLT           /* (ORANET TYPE) Floating point number */
	DTypeNumber        DType = C.SQLT_VNU           /* NUM with preceding length byte */
	DTypeLong          DType = C.SQLT_LNG           /* long */
	DTypeNON           DType = C.SQLT_NON           /* Null/empty PCC Descriptor entry */
	DTypeRowID         DType = C.SQLT_RID           /* rowid */
	DTypeVarBinary     DType = C.SQLT_VBI           /* binary in VCS format */
	DTypeVarBinary2    DType = C.SQLT_BIN           /* binary data(DTYBIN) */
	DTypeLongBinary    DType = C.SQLT_LBI           /* long binary */
	DTypeCursor        DType = C.SQLT_CUR           /* cursor  type */
	DTypeRowIDDesc     DType = C.SQLT_RDD           /* rowid descriptor */
	DTypeLabel         DType = C.SQLT_LAB           /* label type */
	DTypeOSLabel       DType = C.SQLT_OSL           /* oslabel type */
	DTypeNamedType     DType = C.SQLT_NTY           /* named object type */
	DTypeRef           DType = C.SQLT_REF           /* ref type */
	DTypeCLOB          DType = C.SQLT_CLOB          /* character lob */
	DTypeBLOB          DType = C.SQLT_BLOB          /* binary lob */
	DTypeBFile         DType = C.SQLT_BFILEE        /* binary file lob */
	DTypeCFile         DType = C.SQLT_CFILEE        /* character file lob */
	DTypeResultSet     DType = C.SQLT_RSET          /* result set type */
	DTypeNamedColl     DType = C.SQLT_NCO           /* named collection type (varray or nested table) */
	DTypeDate          DType = C.SQLT_ODT           /* OCIDate type */
	DTypeTimestamp     DType = C.SQLT_TIMESTAMP     /* TIMESTAMP */
	DTypeTimestampTZ   DType = C.SQLT_TIMESTAMP_TZ  /* TIMESTAMP WITH TIME ZONE */
	DTypeIntervalYM    DType = C.SQLT_INTERVAL_YM   /* INTERVAL YEAR TO MONTH */
	DTypeIntervalDS    DType = C.SQLT_INTERVAL_DS   /* INTERVAL DAY TO SECOND */
	DTypeTimestampLTZ  DType = C.SQLT_TIMESTAMP_LTZ /* TIMESTAMP WITH LOCAL TZ */
	DTypePSQLNamedType DType = C.SQLT_PNTY          /* pl/sql representation of named types */
)

type BindArray []interface{}

/*
func (stmt *Statement) BindByPos(position uint, valuep unsafe.Pointer, val_sz int32, typ DType) (*Bind, error) {

	rslt := &Bind{}

	var ind C.sb2

	if valuep == nil {
		ind = -1
	}

	stmt.ses.err.maybeError(
		C.OCIBindByPos(
			stmt.stm,               //OCIStmt      *stmtp,
			&rslt.bindhndl,         //OCIBind      **bindpp,
			stmt.ses.err.hndl,      //OCIError     *errhp,
			(C.ub4)(position),      //ub4          position,
			valuep,                 //void         *valuep,
			(C.sb4)(val_sz),        //sb4          value_sz,
			(C.ub2)(typ),           //ub2          dty,
			(unsafe.Pointer)(&ind), //void         *indp,
			nil,            //ub2          *alenp,
			nil,            //ub2          *rcodep,
			0,              //ub4          maxarr_len,
			nil,            //ub4          *curelep,
			C.OCI_DEFAULT)) //ub4          mode

	return rslt, stmt.ses.err
}

func (stmt *Statement) BindByPos(position uint, bindval interface{}) (*Bind, error) {

	rslt := &Bind{}
	var valuep unsafe.Pointer
	var val_sz C.sb4
	var dty C.ub2
	var ind C.sb2

	if bindval == nil {
		ind = -1
	} else {
		ind = 0
		switch xx := bindval.(type) {
		case string:
			zz := []byte(xx)
			valuep = (unsafe.Pointer)(&zz[0])
			val_sz = (C.sb4)(len(zz))
			dty = C.SQLT_CHR
		case uint, uint8, uint16, uint32, uint64:
			zz := reflect.ValueOf(xx).Uint()
			valuep = (unsafe.Pointer)(&zz)
			val_sz = (C.sb4)(unsafe.Sizeof(zz))
			dty = C.SQLT_UIN
		case int, int8, int16, int32, int64:
			zz := reflect.ValueOf(xx).Int()
			valuep = (unsafe.Pointer)(&zz)
			val_sz = (C.sb4)(unsafe.Sizeof(zz))
			dty = C.SQLT_INT
		case float32:
			valuep = (unsafe.Pointer)(&xx)
			val_sz = (C.sb4)(unsafe.Sizeof(xx))
			dty = C.SQLT_FLT
		case float64:
			valuep = (unsafe.Pointer)(&xx)
			val_sz = (C.sb4)(unsafe.Sizeof(xx))
			dty = C.SQLT_FLT
		case time.Time:
			fmt.Println("Binding time.Time")
			zz, _ := stmt.ses.TimeStampFromGoTime(TypeTimestampTZ, xx)
			valuep = (unsafe.Pointer)(zz.datetime)
			val_sz = (C.sb4)(unsafe.Sizeof(zz.datetime))
			dty = C.SQLT_TIMESTAMP_TZ
		case time.Duration:
			// skip this for now
		case *Number:
			valuep = (unsafe.Pointer)(&xx.number)
			val_sz = (C.sb4)(unsafe.Sizeof(xx.number))
			dty = C.SQLT_VNU
		case *TimeStamp:
			fmt.Println("Binding TimeStamp")
			valuep = (unsafe.Pointer)(xx.datetime)
			val_sz = (C.sb4)(unsafe.Sizeof(xx.datetime))
			switch xx.tstype {
			case TypeTimestamp:
				dty = C.SQLT_TIMESTAMP
			case TypeTimestampTZ:
				dty = C.SQLT_TIMESTAMP_TZ
			case TypeTimestampLTZ:
				dty = C.SQLT_TIMESTAMP_LTZ
			}
		case *Interval:
			fmt.Println("Binding Interval")
			valuep = (unsafe.Pointer)(xx.interval)
			val_sz = (C.sb4)(unsafe.Sizeof(xx.interval))
			switch xx.intype {
			case TypeIntervalDS:
				dty = C.SQLT_INTERVAL_DS
			case TypeIntervalYM:
				dty = C.SQLT_INTERVAL_YM
			}
		case bool:
			fmt.Println("Binding bool")
			var zz []byte
			if xx {
				zz = []byte("*")
			} else {
				zz = []byte(" ")
			}
			valuep = (unsafe.Pointer)(&zz[0])
			val_sz = (C.sb4)(len(zz))
			dty = C.SQLT_CHR
		case BindArray:
			// skip this for now
		case []byte:
			fmt.Println("Binding RAW")
			valuep = (unsafe.Pointer)(&xx[0])
			val_sz = (C.sb4)(len(xx))
			dty = C.SQLT_BIN
		default:
			panic("Could not determine data type for BindByPos")
		}
	}

	stmt.ses.err.maybeError(
		C.OCIBindByPos(
			stmt.stm,          //OCIStmt      *stmtp,
			&rslt.bindhndl,    //OCIBind      **bindpp,
			stmt.ses.err.hndl, //OCIError     *errhp,
			(C.ub4)(position), //ub4          position,
			valuep,            //void         *valuep,
			val_sz,            //sb4          value_sz,
			dty,               //ub2          dty,
			(unsafe.Pointer)(&ind), //void         *indp,
			nil,            //ub2          *alenp,
			nil,            //ub2          *rcodep,
			0,              //ub4          maxarr_len,
			nil,            //ub4          *curelep,
			C.OCI_DEFAULT)) //ub4          mode

	return rslt
}

func (stmt *Statement) Bind(values ...interface{}) {
	//  for key, value := range values {

	//  }
}
*/
