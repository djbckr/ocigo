package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"regexp"
	"unsafe"
)

// package global OCI environment and error handles
var (
	genv *C.OCIEnv
	gerr *C.OCIError
)

type Pool struct {
	pool        *C.OCISPool // session pool handle
	err         *C.OCIError // pool error handle
	poolName    *C.OraText  // name of pool, assigned by OCI
	poolNameLen C.ub4
	username    []byte
	password    []byte
	database    []byte
}

func CreatePool(connectString string, minSessions, maxSessions, incrStep int) (*Pool, error) {

	// validate inputs...
	if minSessions < 1 {
		panic("minSessions must be 1 or more")
	}
	if maxSessions < 1 {
		panic("maxSessions must be 1 or more")
	}
	if incrStep < 1 {
		panic("incrStep must be 1 or more")
	}

	re := regexp.MustCompile(`^(.*?)\/(.*)@(.*)$`)
	lst := re.FindStringSubmatch(connectString)

	if len(lst) != 4 {
		return nil, errors.New(`Connect String must be of the format: "username/password@<host/dbname> or <TNS-name>" (without quotes, <>)`)
	}

	// create pool object
	rslt := &Pool{
		username: []byte(lst[1]),
		password: []byte(lst[2]),
		database: []byte(lst[3])}

	ociHandleAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.err)), htypeError)
	ociHandleAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.pool)), htypeSessionPool)

	// set pool inactive timeout for 90 seconds
	var timeout C.ub4 = 90

	err := ociAttrSet((unsafe.Pointer)(rslt.pool), htypeSessionPool, (unsafe.Pointer)(&timeout), 0, attrSessPoolTimeout, gerr)
	if err != nil {
		return nil, processError(err)
	}

	err = checkError(
		C.OCISessionPoolCreate(
			genv, gerr, rslt.pool,
			&rslt.poolName, &rslt.poolNameLen,
			(*C.OraText)(&rslt.database[0]), (C.ub4)(len(rslt.database)),
			(C.ub4)(minSessions), (C.ub4)(maxSessions), (C.ub4)(incrStep),
			(*C.OraText)(&rslt.username[0]), (C.ub4)(len(rslt.username)),
			(*C.OraText)(&rslt.password[0]), (C.ub4)(len(rslt.password)),
			// homogenous and statement caching
			C.OCI_SPC_HOMOGENEOUS+C.OCI_SPC_STMTCACHE), gerr)

	return rslt, processError(err)

}

func (pool *Pool) Destroy() {
	err := checkError(C.OCISessionPoolDestroy(pool.pool, pool.err, C.OCI_SPD_FORCE), pool.err)
	ociHandleFree((unsafe.Pointer)(pool.pool), htypeSessionPool)
	if err != nil {
		panic(err)
	}
}

type Session struct {
	svc *C.OCISvcCtx  // service context handle (associates connection with session)
	err *C.OCIError   // session error handle
	ses *C.OCISession // session handle - used for date/time/number functions
}

func (pool *Pool) GetSession() (*Session, error) {

	rslt := &Session{}
	ociHandleAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.err)), htypeError)

	// get the session (which actually returns the service handle, not the session... )
	err := checkError(
		C.OCISessionGet(
			genv, rslt.err, &rslt.svc,
			nil, pool.poolName, pool.poolNameLen,
			nil, 0, nil, nil, nil, C.OCI_SESSGET_SPOOL), rslt.err)

	if err != nil {
		if err.IsError() {
			return nil, err.err
		} else {
			fmt.Println(err.Warning())
		}
	}

	err = ociAttrGet(
		(unsafe.Pointer)(rslt.svc),
		htypeSvcCtx,
		(unsafe.Pointer)(&rslt.ses),
		nil,
		attrSession,
		rslt.err)

	return rslt, processError(err)

}

func (sess *Session) Commit() error {
	err := checkError(
		C.OCITransCommit(
			sess.svc,
			sess.err,
			C.OCI_DEFAULT), sess.err)

	return processError(err)
}

func (sess *Session) Rollback() error {
	err := checkError(
		C.OCITransRollback(
			sess.svc,
			sess.err,
			C.OCI_DEFAULT), sess.err)

	return processError(err)
}

type TxnType C.ub4

const (
	TxnTypeNormal       TxnType = C.OCI_TRANS_NEW
	TxnTypeReadOnly     TxnType = C.OCI_TRANS_READONLY
	TxnTypeSerializable TxnType = C.OCI_TRANS_SERIALIZABLE
)

func (sess *Session) StartTransaction(txType TxnType) error {
	err := checkError(
		C.OCITransStart(
			sess.svc,
			sess.err,
			0,
			(C.ub4)(txType)), sess.err)

	return processError(err)
}

func (sess *Session) FreeSession() error {

	err := checkError(
		C.OCISessionRelease(sess.svc, sess.err, nil, 0, C.OCI_DEFAULT), sess.err)

	sess.svc = nil
	sess.err = nil
	sess.ses = nil

	return processError(err)

}

func (sess *Session) SetClientIdentifier(value string) {
	//OCI_ATTR_CLIENT_IDENTIFIER

}

func (sess *Session) SetCurrentSchema(value string) {
	//OCI_ATTR_CURRENT_SCHEMA

}

func (sess *Session) SetModule(value string) {
	//OCI_ATTR_MODULE

}

func (sess *Session) SetAction(value string) {
	//OCI_ATTR_ACTION

}

func (sess *Session) SetClientInfo(value string) {
	//OCI_ATTR_CLIENT_INFO

}

func (sess *Session) SetLobPrefetchSize(value string) {
	//OCI_ATTR_DEFAULT_LOBPREFETCH_SIZE

}

func init() {

	errcode := C.OCIEnvCreate(&genv, C.OCI_THREADED+C.OCI_OBJECT+C.OCI_EVENTS, nil, nil, nil, nil, 0, nil)

	if errcode != 0 {
		panic(fmt.Errorf("OCIEnvCreate failed with errcode = %d.\n", errcode))
	}

	ociHandleAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&gerr)), htypeError)

}
