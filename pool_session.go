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
	"time"
	"unsafe"
)

// package global OCI environment and error handles
var (
	genv *C.OCIEnv
	gerr *C.OCIError
)

// Pool is an opaque structure that manages a connection pool to an Oracle database.
type Pool struct {
	pool        *C.OCISPool // session pool handle
	err         *C.OCIError // pool error handle
	poolName    *C.OraText  // name of pool, assigned by OCI
	poolNameLen C.ub4
	username    []byte
	password    []byte
	database    []byte
}

// CreatePool initializes a connection to a database and returns a Pool structure.
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

	rslt.err = (*C.OCIError)(ociHandleAlloc((unsafe.Pointer)(genv), htypeError))
	rslt.pool = (*C.OCISPool)(ociHandleAlloc((unsafe.Pointer)(genv), htypeSessionPool))

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

// Destroy shuts down all connections to the pool.
func (pool *Pool) Destroy() {
	err := checkError(C.OCISessionPoolDestroy(pool.pool, pool.err, C.OCI_SPD_FORCE), pool.err)
	ociHandleFree((unsafe.Pointer)(pool.pool), htypeSessionPool)
	if err != nil {
		panic(err)
	}
}

func (pool *Pool) SetConnectionTimeout(duration time.Duration) {
	//OCI_ATTR_SPOOL_TIMEOUT
}

func (pool *Pool) GetConnectionTimeout() time.Duration {
	//OCI_ATTR_SPOOL_TIMEOUT
	return 0
}

type AcquireMode C.ub1

const (
	AcquireModeWait   AcquireMode = C.OCI_SPOOL_ATTRVAL_WAIT
	AcquireModeNoWait AcquireMode = C.OCI_SPOOL_ATTRVAL_NOWAIT
	AcquireModeForce  AcquireMode = C.OCI_SPOOL_ATTRVAL_FORCEGET
)

func (pool *Pool) SetAcquireNoWait(value AcquireMode) {
	//OCI_ATTR_SPOOL_GETMODE
}

func (pool *Pool) GetAcquireMode() AcquireMode {
	//OCI_ATTR_SPOOL_GETMODE
	return AcquireModeWait
}

func (pool *Pool) GetNumBusyConnections() int32 {
	//OCI_ATTR_SPOOL_BUSY_COUNT
	return 0
}

func (pool *Pool) GetNumOpenConnections() int32 {
	//OCI_ATTR_SPOOL_OPEN_COUNT
	return 0
}

func (pool *Pool) SetStatementCacheSize(value uint32) {
	//OCI_ATTR_SPOOL_STMTCACHESIZE
}

type Session struct {
	svc *C.OCISvcCtx  // service context handle (associates connection with session)
	err *C.OCIError   // session error handle
	ses *C.OCISession // session handle - used for date/time/number functions
}

// Acquire gets a session from the pool in order to execute SQL against the database.
func (pool *Pool) Acquire() (*Session, error) {

	rslt := &Session{}
	rslt.err = (*C.OCIError)(ociHandleAlloc((unsafe.Pointer)(genv), htypeError))

	// get the session (which actually returns the service handle, not the session... )
	err := checkError(
		C.OCISessionGet(
			genv, rslt.err, &rslt.svc,
			nil, pool.poolName, pool.poolNameLen,
			nil, 0, nil, nil, nil, C.OCI_SESSGET_SPOOL), rslt.err)

	if err != nil {
		if err.IsError() {
			return nil, err.err
		}

		fmt.Println(err.Warning())
	}

	// now actually get the session handle (I know, right?)
	ses, err := ociAttrGetPointer(
		(unsafe.Pointer)(rslt.svc),
		htypeSvcCtx,
		attrSession,
		rslt.err)

	rslt.ses = (*C.OCISession)(ses)

	return rslt, processError(err)

}

// Commit issues a commit to the database.
func (sess *Session) Commit() error {
	err := checkError(
		C.OCITransCommit(
			sess.svc,
			sess.err,
			C.OCI_DEFAULT), sess.err)

	return processError(err)
}

// Rollback issues a rollback to the database.
func (sess *Session) Rollback() error {
	err := checkError(
		C.OCITransRollback(
			sess.svc,
			sess.err,
			C.OCI_DEFAULT), sess.err)

	return processError(err)
}

// TxnType type defines a Go type for transaction types.
type TxnType C.ub4

// The available transaction types.
const (
	TxnTypeNormal       TxnType = C.OCI_TRANS_NEW
	TxnTypeReadOnly     TxnType = C.OCI_TRANS_READONLY
	TxnTypeSerializable TxnType = C.OCI_TRANS_SERIALIZABLE
)

// StartTransaction allows you to manually start a transaction in one of the TxnType modes.
func (sess *Session) StartTransaction(txType TxnType) error {

	err := checkError(
		C.OCITransStart(
			sess.svc,
			sess.err,
			0,
			(C.ub4)(txType)), sess.err)

	return processError(err)

}

// Release puts the session back in the pool for reuse.
func (sess *Session) Release() error {

	err := checkError(
		C.OCISessionRelease(sess.svc, sess.err, nil, 0, C.OCI_DEFAULT), sess.err)

	sess.svc = nil
	sess.err = nil
	sess.ses = nil

	return processError(err)

}

func (sess *Session) SetClientIdentifier(value string) {
	//OCI_ATTR_CLIENT_IDENTIFIER up to 64 bytes

}

func (sess *Session) SetCurrentSchema(value string) {
	//OCI_ATTR_CURRENT_SCHEMA

}

func (sess *Session) GetCurrentSchema() string {
	//OCI_ATTR_CURRENT_SCHEMA
	return ""
}

func (sess *Session) SetModule(value string) {
	//OCI_ATTR_MODULE

}

func (sess *Session) SetAction(value string) {
	//OCI_ATTR_ACTION

}

func (sess *Session) SetClientInfo(value string) {
	//OCI_ATTR_CLIENT_INFO up to 64 bytes

}

func (sess *Session) SetLobPrefetchSize(value uint32) {
	//OCI_ATTR_DEFAULT_LOBPREFETCH_SIZE

}

func (sess *Session) GetLobPrefetchSize() uint32 {
	//OCI_ATTR_DEFAULT_LOBPREFETCH_SIZE
	return 0
}

func (sess *Session) SetCollectCallTime(value bool) {
	//OCI_ATTR_COLLECT_CALL_TIME
}

func (sess *Session) GetCollectCallTime() bool {
	//OCI_ATTR_COLLECT_CALL_TIME
	return false
}

func (sess *Session) GetCallTime() time.Duration {
	//OCI_ATTR_CALL_TIME returns ub8
	return 0
}

func (sess *Session) SetEdition(value string) {
	//OCI_ATTR_EDITION
}

func (sess *Session) GetEdition() string {
	//OCI_ATTR_EDITION
	return ""
}

type SessionState C.ub1

const (
	Stateful  SessionState = C.OCI_SESSION_STATEFUL
	Stateless SessionState = C.OCI_SESSION_STATELESS
)

func (sess *Session) SetState(value SessionState) {
	//OCI_ATTR_SESSION_STATE
}

func (sess *Session) GetState() SessionState {
	//OCI_ATTR_SESSION_STATE
	return Stateful
}

func (sess *Session) IsTransactionInProgress() bool {
	//OCI_ATTR_TRANSACTION_IN_PROGRESS
	return false
}

// package initialization - there is only one environment for whoever uses this library.

func init() {

	errcode := C.OCIEnvCreate(&genv, C.OCI_THREADED+C.OCI_OBJECT+C.OCI_EVENTS, nil, nil, nil, nil, 0, nil)

	if errcode != 0 {
		panic(fmt.Errorf("OCIEnvCreate failed with errcode = %d.\n", errcode))
	}

	gerr = (*C.OCIError)(ociHandleAlloc((unsafe.Pointer)(genv), htypeError))

}
