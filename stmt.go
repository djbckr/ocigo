package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"

import (
	"crypto/sha256"
	"encoding/hex"
	//"errors"
	//"fmt"
	"unsafe"
	//"time"
)

type StmtType uint16

const (
	StmtSelect  StmtType = C.OCI_STMT_SELECT
	StmtUpdate  StmtType = C.OCI_STMT_UPDATE
	StmtDelete  StmtType = C.OCI_STMT_DELETE
	StmtInsert  StmtType = C.OCI_STMT_INSERT
	StmtCreate  StmtType = C.OCI_STMT_CREATE
	StmtDrop    StmtType = C.OCI_STMT_DROP
	StmtAlter   StmtType = C.OCI_STMT_ALTER
	StmtBegin   StmtType = C.OCI_STMT_BEGIN
	StmtDeclare StmtType = C.OCI_STMT_DECLARE
)

type Statement struct {
	ses     *Session
	err     *C.OCIError
	stm     *C.OCIStmt
	key     []byte
	qry     []byte
	stmtype StmtType
}

func (stmt Statement) StatementType() StmtType {
	return stmt.stmtype
}

func stmtFinalizer(stmt *Statement) {
	stmt.Release(false)
}

func (sess *Session) Prepare(sql string) (*Statement, error) {
	rslt := &Statement{ses: sess, qry: []byte(sql)}

	// hash the query, turn to slice, output to hex string, convert to []byte
	hash := sha256.Sum256(rslt.qry)
	rslt.key = []byte(hex.EncodeToString(hash[:]))

	rslt.err = (*C.OCIError)(ociHandleAlloc(unsafe.Pointer(genv), htypeError))

	vErr := checkError(
		C.OCIStmtPrepare2(
			sess.svc,
			&rslt.stm,
			rslt.err,
			nil, 0,
			(*C.OraText)(&rslt.key[0]),
			(C.ub4)(len(rslt.key)),
			C.OCI_NTV_SYNTAX,
			C.OCI_PREP2_CACHE_SEARCHONLY), rslt.err)

	if vErr == nil {
		stype, vErr := ociAttrGetUB2(unsafe.Pointer(rslt.stm), htypeStatement, attrStmtType, rslt.err)
		rslt.stmtype = StmtType(stype)
		return rslt, vErr
	}

	if vErr.code == 24431 {
		// it's not in the cache (and hope to Gods that number never changes), so we need to create it

		vErr = checkError(
			C.OCIStmtPrepare2(
				sess.svc,
				&rslt.stm,
				rslt.err,
				(*C.OraText)(&rslt.qry[0]),
				(C.ub4)(len(rslt.qry)),
				nil, 0,
				C.OCI_NTV_SYNTAX,
				C.OCI_DEFAULT), rslt.err)

		if vErr == nil {
			stype, err := ociAttrGetUB2(
				unsafe.Pointer(rslt.stm),
				htypeStatement,
				attrStmtType,
				rslt.err)

			rslt.stmtype = StmtType(stype)

			return rslt, processError(err)
		}
		return nil, processError(vErr)
	}
	return nil, processError(vErr)
}

func (stmt *Statement) exec(iterations uint32, commit bool) *OciError {
	// 8=1  16=2  32=4  64=8
	var flags C.ub4 = C.OCI_DEFAULT
	var iters C.ub4 = (C.ub4)(iterations)

	if commit {
		flags = C.OCI_COMMIT_ON_SUCCESS
	}

	vErr := checkError(
		C.OCIStmtExecute(
			stmt.ses.svc,
			stmt.stm,
			stmt.err,
			iters, 0, nil, nil, flags), stmt.err)

	return vErr

}

func (stmt *Statement) Execute() error {
	return processError(stmt.exec(1, false))
}

func (stmt *Statement) ExecuteAndCommit() error {
	return processError(stmt.exec(1, true))
}

func (stmt *Statement) Query() (*ResultSet, error) {
	return stmt.query(0) // in resultset.go
}

func (stmt *Statement) QueryLimit(count uint32) (*ResultSet, error) {
	return stmt.query(count) // in resultset.go
}

func (stmt *Statement) Release(KeepInCache bool) error {

	if stmt.stm != nil {

		var mode C.ub4

		if KeepInCache {
			mode = C.OCI_DEFAULT
		} else {
			mode = C.OCI_STRLS_CACHE_DELETE
		}

		vErr := checkError(
			C.OCIStmtRelease(
				stmt.stm,
				stmt.err,
				(*C.OraText)(&stmt.key[0]),
				(C.ub4)(len(stmt.key)),
				mode), stmt.err)

		rslt := processError(vErr)

		stmt.stm = nil
		stmt.ses = nil
		ociHandleFree(unsafe.Pointer(stmt.err), htypeError)
		stmt.err = nil

		return rslt
	}

	return nil

}

func (rset *ResultSet) Next() []interface{} {
	return nil
}
