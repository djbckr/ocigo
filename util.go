package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"

import (
	"bytes"
	"fmt"
	"unsafe"
)

func nulTerminatedByteToString(buf []byte) string {
	// find null terminator
	z := bytes.IndexByte(buf, 0)
	return string(buf[:z])
}

type ociSqlType C.ub2

const (
	sqltRef             ociSqlType = C.OCI_TYPECODE_REF             /* SQL/OTS OBJECT REFERENCE */
	sqltDate            ociSqlType = C.OCI_TYPECODE_DATE            /* SQL DATE  OTS DATE */
	sqltSigned8         ociSqlType = C.OCI_TYPECODE_SIGNED8         /* SQL SIGNED INTEGER(8)  OTS SINT8 */
	sqltSigned16        ociSqlType = C.OCI_TYPECODE_SIGNED16        /* SQL SIGNED INTEGER(16)  OTS SINT16 */
	sqltSigned32        ociSqlType = C.OCI_TYPECODE_SIGNED32        /* SQL SIGNED INTEGER(32)  OTS SINT32 */
	sqltReal            ociSqlType = C.OCI_TYPECODE_REAL            /* SQL REAL  OTS SQL_REAL */
	sqltDouble          ociSqlType = C.OCI_TYPECODE_DOUBLE          /* SQL DOUBLE PRECISION  OTS SQL_DOUBLE */
	sqltBFloat          ociSqlType = C.OCI_TYPECODE_BFLOAT          /* Binary float */
	sqltBDouble         ociSqlType = C.OCI_TYPECODE_BDOUBLE         /* Binary double */
	sqltFloat           ociSqlType = C.OCI_TYPECODE_FLOAT           /* SQL FLOAT(P)  OTS FLOAT(P) */
	sqltNumber          ociSqlType = C.OCI_TYPECODE_NUMBER          /* SQL NUMBER(P S)  OTS NUMBER(P S) */
	sqltDecimal         ociSqlType = C.OCI_TYPECODE_DECIMAL         /* SQL DECIMAL(P S)  OTS DECIMAL(P S) */
	sqltUnsigned8       ociSqlType = C.OCI_TYPECODE_UNSIGNED8       /* SQL UNSIGNED INTEGER(8)  OTS UINT8 */
	sqltUnsigned16      ociSqlType = C.OCI_TYPECODE_UNSIGNED16      /* SQL UNSIGNED INTEGER(16)  OTS UINT16 */
	sqltUnsigned32      ociSqlType = C.OCI_TYPECODE_UNSIGNED32      /* SQL UNSIGNED INTEGER(32)  OTS UINT32 */
	sqltOctet           ociSqlType = C.OCI_TYPECODE_OCTET           /* SQL ???  OTS OCTET */
	sqltSmallint        ociSqlType = C.OCI_TYPECODE_SMALLINT        /* SQL SMALLINT  OTS SMALLINT */
	sqltInteger         ociSqlType = C.OCI_TYPECODE_INTEGER         /* SQL INTEGER  OTS INTEGER */
	sqltRaw             ociSqlType = C.OCI_TYPECODE_RAW             /* SQL RAW(N)  OTS RAW(N) */
	sqltPtr             ociSqlType = C.OCI_TYPECODE_PTR             /* SQL POINTER  OTS POINTER */
	sqltVarchar2        ociSqlType = C.OCI_TYPECODE_VARCHAR2        /* SQL VARCHAR2(N)  OTS SQL_VARCHAR2(N) */
	sqltChar            ociSqlType = C.OCI_TYPECODE_CHAR            /* SQL CHAR(N)  OTS SQL_CHAR(N) */
	sqltVarchar         ociSqlType = C.OCI_TYPECODE_VARCHAR         /* SQL VARCHAR(N)  OTS SQL_VARCHAR(N) */
	sqltMlsLabel        ociSqlType = C.OCI_TYPECODE_MLSLABEL        /* OTS MLSLABEL */
	sqltVArray          ociSqlType = C.OCI_TYPECODE_VARRAY          /* SQL VARRAY  OTS PAGED VARRAY */
	sqltTable           ociSqlType = C.OCI_TYPECODE_TABLE           /* SQL TABLE  OTS MULTISET */
	sqltObject          ociSqlType = C.OCI_TYPECODE_OBJECT          /* SQL/OTS NAMED OBJECT TYPE */
	sqltOpaque          ociSqlType = C.OCI_TYPECODE_OPAQUE          /* SQL/OTS Opaque Types */
	sqltNamedCollection ociSqlType = C.OCI_TYPECODE_NAMEDCOLLECTION /* SQL/OTS NAMED COLLECTION TYPE */
	sqltBLOB            ociSqlType = C.OCI_TYPECODE_BLOB            /* SQL/OTS BINARY LARGE OBJECT */
	sqltBFile           ociSqlType = C.OCI_TYPECODE_BFILE           /* SQL/OTS BINARY FILE OBJECT */
	sqltCLOB            ociSqlType = C.OCI_TYPECODE_CLOB            /* SQL/OTS CHARACTER LARGE OBJECT */
	sqltCFile           ociSqlType = C.OCI_TYPECODE_CFILE           /* SQL/OTS CHARACTER FILE OBJECT */
	sqltTime            ociSqlType = C.OCI_TYPECODE_TIME            /* SQL/OTS TIME */
	sqltTimeTZ          ociSqlType = C.OCI_TYPECODE_TIME_TZ         /* SQL/OTS TIME_TZ */
	sqltTimestamp       ociSqlType = C.OCI_TYPECODE_TIMESTAMP       /* SQL/OTS TIMESTAMP */
	sqltTimestampTZ     ociSqlType = C.OCI_TYPECODE_TIMESTAMP_TZ    /* SQL/OTS TIMESTAMP_TZ */
	sqltTimestampLTZ    ociSqlType = C.OCI_TYPECODE_TIMESTAMP_LTZ   /* TIMESTAMP_LTZ */
	sqltIntervalYM      ociSqlType = C.OCI_TYPECODE_INTERVAL_YM     /* SQL/OTS INTRVL YR-MON */
	sqltIntervalDS      ociSqlType = C.OCI_TYPECODE_INTERVAL_DS     /* SQL/OTS INTRVL DAY-SEC */
	sqltURowID          ociSqlType = C.OCI_TYPECODE_UROWID          /* Urowid type */
	sqltPlsInteger      ociSqlType = C.OCI_TYPECODE_PLS_INTEGER     /* type code for PLS_INTEGER */
	sqltNChar           ociSqlType = C.OCI_TYPECODE_NCHAR
	sqltNVarchar2       ociSqlType = C.OCI_TYPECODE_NVARCHAR2
	sqltNCLOB           ociSqlType = C.OCI_TYPECODE_NCLOB
	sqltNone            ociSqlType = C.OCI_TYPECODE_NONE
)

type ociHandleType C.ub4

const (
	htypeEnv               ociHandleType = C.OCI_HTYPE_ENV                  /* environment handle */
	htypeError             ociHandleType = C.OCI_HTYPE_ERROR                /* error handle */
	htypeSvcCtx            ociHandleType = C.OCI_HTYPE_SVCCTX               /* service handle */
	htypeStatement         ociHandleType = C.OCI_HTYPE_STMT                 /* statement handle */
	htypeBind              ociHandleType = C.OCI_HTYPE_BIND                 /* bind handle */
	htypeDefine            ociHandleType = C.OCI_HTYPE_DEFINE               /* define handle */
	htypeDescribe          ociHandleType = C.OCI_HTYPE_DESCRIBE             /* describe handle */
	htypeServer            ociHandleType = C.OCI_HTYPE_SERVER               /* server handle */
	htypeSession           ociHandleType = C.OCI_HTYPE_SESSION              /* authentication handle */
	htypeAuthInfo          ociHandleType = C.OCI_HTYPE_AUTHINFO             /* SessionGet auth handle */
	htypeTrans             ociHandleType = C.OCI_HTYPE_TRANS                /* transaction handle */
	htypeComplexObject     ociHandleType = C.OCI_HTYPE_COMPLEXOBJECT        /* complex object retrieval handle */
	htypeSecurity          ociHandleType = C.OCI_HTYPE_SECURITY             /* security handle */
	htypeSubscription      ociHandleType = C.OCI_HTYPE_SUBSCRIPTION         /* subscription handle */
	htypeDirpathCtx        ociHandleType = C.OCI_HTYPE_DIRPATH_CTX          /* direct path context */
	htypeDirpathColArray   ociHandleType = C.OCI_HTYPE_DIRPATH_COLUMN_ARRAY /* direct path column array */
	htypeDirpathStream     ociHandleType = C.OCI_HTYPE_DIRPATH_STREAM       /* direct path stream */
	htypeProc              ociHandleType = C.OCI_HTYPE_PROC                 /* process handle */
	htypeDirpathFnCtx      ociHandleType = C.OCI_HTYPE_DIRPATH_FN_CTX       /* direct path function context */
	htypeDirpathFnColArray ociHandleType = C.OCI_HTYPE_DIRPATH_FN_COL_ARRAY /* dp object column array */
	htypeXadSession        ociHandleType = C.OCI_HTYPE_XADSESSION           /* access driver session */
	htypeXadTable          ociHandleType = C.OCI_HTYPE_XADTABLE             /* access driver table */
	htypeXadField          ociHandleType = C.OCI_HTYPE_XADFIELD             /* access driver field */
	htypeXadGranule        ociHandleType = C.OCI_HTYPE_XADGRANULE           /* access driver granule */
	htypeXadRecord         ociHandleType = C.OCI_HTYPE_XADRECORD            /* access driver record */
	htypeXadIO             ociHandleType = C.OCI_HTYPE_XADIO                /* access driver I/O */
	htypeConnectionPool    ociHandleType = C.OCI_HTYPE_CPOOL                /* connection pool handle */
	htypeSessionPool       ociHandleType = C.OCI_HTYPE_SPOOL                /* session pool handle */
	htypeAdmin             ociHandleType = C.OCI_HTYPE_ADMIN                /* admin handle */
	htypeEvent             ociHandleType = C.OCI_HTYPE_EVENT                /* HA event handle */
)

// handle is a **
func ociHandleAlloc(parent unsafe.Pointer, handle *unsafe.Pointer, htype ociHandleType) {

	errcode := C.OCIHandleAlloc(parent, handle, (C.ub4)(htype), 0, nil)

	if errcode != 0 {
		panic(fmt.Errorf("OCIHandleAlloc(HTYPE_ERROR) failed with errcode = %d.\n", errcode))
	}

}

func ociHandleFree(handle unsafe.Pointer, htype ociHandleType) {
	C.OCIHandleFree(handle, (C.ub4)(htype))
}

type ociDescriptorType C.ub4

const (
	dtypeLOB                ociDescriptorType = C.OCI_DTYPE_LOB                  /* lob locator */
	dtypeSnap               ociDescriptorType = C.OCI_DTYPE_SNAP                 /* snapshot descriptor */
	dtypeResultSet          ociDescriptorType = C.OCI_DTYPE_RSET                 /* result set descriptor */
	dtypeParam              ociDescriptorType = C.OCI_DTYPE_PARAM                /* a parameter descriptor obtained from ocigparm */
	dtypeRowID              ociDescriptorType = C.OCI_DTYPE_ROWID                /* rowid descriptor */
	dtypeComplexObject      ociDescriptorType = C.OCI_DTYPE_COMPLEXOBJECTCOMP    /* complex object retrieval descriptor */
	dtypeFile               ociDescriptorType = C.OCI_DTYPE_FILE                 /* File Lob locator */
	dtypeAQEnqOptions       ociDescriptorType = C.OCI_DTYPE_AQENQ_OPTIONS        /* enqueue options */
	dtypeAQDeqOptions       ociDescriptorType = C.OCI_DTYPE_AQDEQ_OPTIONS        /* dequeue options */
	dtypeAQMsgProperties    ociDescriptorType = C.OCI_DTYPE_AQMSG_PROPERTIES     /* message properties */
	dtypeAQAgent            ociDescriptorType = C.OCI_DTYPE_AQAGENT              /* aq agent */
	dtypeLOBLocator         ociDescriptorType = C.OCI_DTYPE_LOCATOR              /* LOB locator */
	dtypeIntervalYM         ociDescriptorType = C.OCI_DTYPE_INTERVAL_YM          /* Interval year month */
	dtypeIntervalDS         ociDescriptorType = C.OCI_DTYPE_INTERVAL_DS          /* Interval day second */
	dtypeAQNotify           ociDescriptorType = C.OCI_DTYPE_AQNFY_DESCRIPTOR     /* AQ notify descriptor */
	dtypeDate               ociDescriptorType = C.OCI_DTYPE_DATE                 /* Date */
	dtypeTime               ociDescriptorType = C.OCI_DTYPE_TIME                 /* Time */
	dtypeTimeTZ             ociDescriptorType = C.OCI_DTYPE_TIME_TZ              /* Time with timezone */
	dtypeTimestamp          ociDescriptorType = C.OCI_DTYPE_TIMESTAMP            /* Timestamp */
	dtypeTimestampTZ        ociDescriptorType = C.OCI_DTYPE_TIMESTAMP_TZ         /* Timestamp with timezone */
	dtypeTimestampLTZ       ociDescriptorType = C.OCI_DTYPE_TIMESTAMP_LTZ        /* Timestamp with local tz */
	dtypeUserCallback       ociDescriptorType = C.OCI_DTYPE_UCB                  /* user callback descriptor */
	dtypeServerDN           ociDescriptorType = C.OCI_DTYPE_SRVDN                /* server DN list descriptor */
	dtypeSignature          ociDescriptorType = C.OCI_DTYPE_SIGNATURE            /* signature */
	dtypeAQListenOptions    ociDescriptorType = C.OCI_DTYPE_AQLIS_OPTIONS        /* AQ listen options */
	dtypeAQListenMsgProps   ociDescriptorType = C.OCI_DTYPE_AQLIS_MSG_PROPERTIES /* AQ listen msg props */
	dtypeChangeNotification ociDescriptorType = C.OCI_DTYPE_CHDES                /* Top level change notification desc */
	dtypeTableChange        ociDescriptorType = C.OCI_DTYPE_TABLE_CHDES          /* Table change descriptor */
	dtypeRowChange          ociDescriptorType = C.OCI_DTYPE_ROW_CHDES            /* Row change descriptor */
	dtypeQueryChange        ociDescriptorType = C.OCI_DTYPE_CQDES                /* Query change descriptor */
	dtypeLOBRegion          ociDescriptorType = C.OCI_DTYPE_LOB_REGION           /* LOB Share region descriptor */
)

// descriptor is a **
func ociDescriptorAlloc(parent unsafe.Pointer, descriptor *unsafe.Pointer, dtype ociDescriptorType) {

	errcode := C.OCIDescriptorAlloc(parent, descriptor, (C.ub4)(dtype), 0, nil)

	if errcode != C.OCI_SUCCESS {
		panic("Unable to allocate descriptor")
	}

}

func ociDescriptorFree(descriptor unsafe.Pointer, dtype ociDescriptorType) {
	C.OCIDescriptorFree(descriptor, (C.ub4)(dtype))
}

type ociAttrType C.ub4

const (
	attrFunctionCode                  ociAttrType = C.OCI_ATTR_FNCODE                  /* the OCI function code */
	attrObject                        ociAttrType = C.OCI_ATTR_OBJECT                  /* is the environment initialized in object mode */
	attrNonBlockingMode               ociAttrType = C.OCI_ATTR_NONBLOCKING_MODE        /* non blocking mode */
	attrSQLCode                       ociAttrType = C.OCI_ATTR_SQLCODE                 /* the SQL verb */
	attrEnvironment                   ociAttrType = C.OCI_ATTR_ENV                     /* the environment handle */
	attrServer                        ociAttrType = C.OCI_ATTR_SERVER                  /* the server handle */
	attrSession                       ociAttrType = C.OCI_ATTR_SESSION                 /* the user session handle */
	attrTrans                         ociAttrType = C.OCI_ATTR_TRANS                   /* the transaction handle */
	attrRowCount                      ociAttrType = C.OCI_ATTR_ROW_COUNT               /* the rows processed so far */
	attrSqlFnCode                     ociAttrType = C.OCI_ATTR_SQLFNCODE               /* the SQL verb of the statement */
	attrPrefetchRows                  ociAttrType = C.OCI_ATTR_PREFETCH_ROWS           /* sets the number of rows to prefetch */
	attrNestedPrefetchRows            ociAttrType = C.OCI_ATTR_NESTED_PREFETCH_ROWS    /* the prefetch rows of nested table*/
	attrPrefetchMemory                ociAttrType = C.OCI_ATTR_PREFETCH_MEMORY         /* memory limit for rows fetched */
	attrNestedPrefetchMemory          ociAttrType = C.OCI_ATTR_NESTED_PREFETCH_MEMORY  /* memory limit for nested rows */
	attrCharCount                     ociAttrType = C.OCI_ATTR_CHAR_COUNT              /* this specifies the bind and define size in characters */
	attrPackedDecimalScale            ociAttrType = C.OCI_ATTR_PDSCL                   /* packed decimal scale */
	attrFSPrecision                   ociAttrType = C.OCI_ATTR_FSPRECISION             /* fs prec for datetime data types */
	attrPackedDecimalFormat           ociAttrType = C.OCI_ATTR_PDPRC                   /* packed decimal format */
	attrLFPrecision                   ociAttrType = C.OCI_ATTR_LFPRECISION             /* fs prec for datetime data types */
	attrParamCount                    ociAttrType = C.OCI_ATTR_PARAM_COUNT             /* number of column in the select list */
	attrRowID                         ociAttrType = C.OCI_ATTR_ROWID                   /* the rowid */
	attrCharset                       ociAttrType = C.OCI_ATTR_CHARSET                 /* the character set value */
	attrNChar                         ociAttrType = C.OCI_ATTR_NCHAR                   /* NCHAR type */
	attrUsername                      ociAttrType = C.OCI_ATTR_USERNAME                /* username attribute */
	attrPassword                      ociAttrType = C.OCI_ATTR_PASSWORD                /* password attribute */
	attrStmtType                      ociAttrType = C.OCI_ATTR_STMT_TYPE               /* statement type */
	attrInternalName                  ociAttrType = C.OCI_ATTR_INTERNAL_NAME           /* user friendly global name */
	attrExternalName                  ociAttrType = C.OCI_ATTR_EXTERNAL_NAME           /* the internal name for global txn */
	attrXID                           ociAttrType = C.OCI_ATTR_XID                     /* XOPEN defined global transaction id */
	attrTransLock                     ociAttrType = C.OCI_ATTR_TRANS_LOCK              /* */
	attrTransName                     ociAttrType = C.OCI_ATTR_TRANS_NAME              /* string to identify a global transaction */
	attrHeapAlloc                     ociAttrType = C.OCI_ATTR_HEAPALLOC               /* memory allocated on the heap */
	attrCharsetID                     ociAttrType = C.OCI_ATTR_CHARSET_ID              /* Character Set ID */
	attrCharsetForm                   ociAttrType = C.OCI_ATTR_CHARSET_FORM            /* Character Set Form */
	attrMaxdataSize                   ociAttrType = C.OCI_ATTR_MAXDATA_SIZE            /* Maximumsize of data on the server  */
	attrCacheOptSize                  ociAttrType = C.OCI_ATTR_CACHE_OPT_SIZE          /* object cache optimal size */
	attrCacheMaxSize                  ociAttrType = C.OCI_ATTR_CACHE_MAX_SIZE          /* object cache maximum size percentage */
	attrPinOption                     ociAttrType = C.OCI_ATTR_PINOPTION               /* object cache default pin option */
	attrAllocDuration                 ociAttrType = C.OCI_ATTR_ALLOC_DURATION          /* object cache default allocation duration */
	attrPinDuration                   ociAttrType = C.OCI_ATTR_PIN_DURATION            /* object cache default pin duration */
	attrFDO                           ociAttrType = C.OCI_ATTR_FDO                     /* Format Descriptor object attribute */
	attrPostprocessingCallback        ociAttrType = C.OCI_ATTR_POSTPROCESSING_CALLBACK /* Callback to process outbind data */
	attrPostprocessingContext         ociAttrType = C.OCI_ATTR_POSTPROCESSING_CONTEXT  /* Callback context to process outbind data */
	attrRowsReturned                  ociAttrType = C.OCI_ATTR_ROWS_RETURNED           /* Number of rows returned in current iter - for Bind handles */
	attrFailoverCallback              ociAttrType = C.OCI_ATTR_FOCBK                   /* Failover Callback attribute */
	attrInV8Mode                      ociAttrType = C.OCI_ATTR_IN_V8_MODE              /* is the server/service context in V8 mode */
	attrLobEmpty                      ociAttrType = C.OCI_ATTR_LOBEMPTY                /* empty lob ? */
	attrSessLang                      ociAttrType = C.OCI_ATTR_SESSLANG                /* session language handle */
	attrVisibility                    ociAttrType = C.OCI_ATTR_VISIBILITY              /* visibility */
	attrRelativeMsgID                 ociAttrType = C.OCI_ATTR_RELATIVE_MSGID          /* relative message id */
	attrSequenceDeviation             ociAttrType = C.OCI_ATTR_SEQUENCE_DEVIATION      /* sequence deviation */
	attrConsumerName                  ociAttrType = C.OCI_ATTR_CONSUMER_NAME           /* consumer name */
	attrDeqMode                       ociAttrType = C.OCI_ATTR_DEQ_MODE                /* dequeue mode */
	attrNavigation                    ociAttrType = C.OCI_ATTR_NAVIGATION              /* navigation */
	attrWait                          ociAttrType = C.OCI_ATTR_WAIT                    /* wait */
	attrDeqMsgID                      ociAttrType = C.OCI_ATTR_DEQ_MSGID               /* dequeue message id */
	attrPriority                      ociAttrType = C.OCI_ATTR_PRIORITY                /* priority */
	attrDelay                         ociAttrType = C.OCI_ATTR_DELAY                   /* delay */
	attrExpiration                    ociAttrType = C.OCI_ATTR_EXPIRATION              /* expiration */
	attrCorrelation                   ociAttrType = C.OCI_ATTR_CORRELATION             /* correlation id */
	attrAttempts                      ociAttrType = C.OCI_ATTR_ATTEMPTS                /* # of attempts */
	attrRecipientList                 ociAttrType = C.OCI_ATTR_RECIPIENT_LIST          /* recipient list */
	attrExceptionQueue                ociAttrType = C.OCI_ATTR_EXCEPTION_QUEUE         /* exception queue name */
	attrEnqTime                       ociAttrType = C.OCI_ATTR_ENQ_TIME                /* enqueue time (only OCIAttrGet) */
	attrMsgState                      ociAttrType = C.OCI_ATTR_MSG_STATE               /* message state (only OCIAttrGet) */
	attrAgentName                     ociAttrType = C.OCI_ATTR_AGENT_NAME              /* agent name */
	attrAgentAddress                  ociAttrType = C.OCI_ATTR_AGENT_ADDRESS           /* agent address */
	attrAgentProtocol                 ociAttrType = C.OCI_ATTR_AGENT_PROTOCOL          /* agent protocol */
	attrUserProperty                  ociAttrType = C.OCI_ATTR_USER_PROPERTY           /* user property */
	attrSenderID                      ociAttrType = C.OCI_ATTR_SENDER_ID               /* sender id */
	attrOriginalMsgID                 ociAttrType = C.OCI_ATTR_ORIGINAL_MSGID          /* original message id */
	attrQueueName                     ociAttrType = C.OCI_ATTR_QUEUE_NAME              /* queue name */
	attrNfyMsgID                      ociAttrType = C.OCI_ATTR_NFY_MSGID               /* message id */
	attrMsgProp                       ociAttrType = C.OCI_ATTR_MSG_PROP                /* message properties */
	attrNumDmlErrors                  ociAttrType = C.OCI_ATTR_NUM_DML_ERRORS          /* num of errs in array DML */
	attrDmlRowOffset                  ociAttrType = C.OCI_ATTR_DML_ROW_OFFSET          /* row offset in the array */
	attrAQNumErrors                   ociAttrType = C.OCI_ATTR_AQ_NUM_ERRORS
	attrAQErrorIndex                  ociAttrType = C.OCI_ATTR_AQ_ERROR_INDEX
	attrDateFormat                    ociAttrType = C.OCI_ATTR_DATEFORMAT                        /* default date format string */
	attrBufAddr                       ociAttrType = C.OCI_ATTR_BUF_ADDR                          /* buffer address */
	attrBufSize                       ociAttrType = C.OCI_ATTR_BUF_SIZE                          /* buffer size */
	attrNumRows                       ociAttrType = C.OCI_ATTR_NUM_ROWS                          /* number of rows in column array */
	attrColCount                      ociAttrType = C.OCI_ATTR_COL_COUNT                         /* columns of column array processed so far.       */
	attrStreamOffset                  ociAttrType = C.OCI_ATTR_STREAM_OFFSET                     /* str off of last row processed */
	attrSharedHeapalloc               ociAttrType = C.OCI_ATTR_SHARED_HEAPALLOC                  /* Shared Heap Allocation Size */
	attrServerGroup                   ociAttrType = C.OCI_ATTR_SERVER_GROUP                      /* server group name */
	attrMigratableSession             ociAttrType = C.OCI_ATTR_MIGSESSION                        /* migratable session attribute */
	attrNoCache                       ociAttrType = C.OCI_ATTR_NOCACHE                           /* Temporary LOBs */
	attrMempoolSize                   ociAttrType = C.OCI_ATTR_MEMPOOL_SIZE                      /* Pool Size */
	attrMempoolInstname               ociAttrType = C.OCI_ATTR_MEMPOOL_INSTNAME                  /* Instance name */
	attrMempoolAppname                ociAttrType = C.OCI_ATTR_MEMPOOL_APPNAME                   /* Application name */
	attrMempoolHomename               ociAttrType = C.OCI_ATTR_MEMPOOL_HOMENAME                  /* Home Directory name */
	attrMempoolModel                  ociAttrType = C.OCI_ATTR_MEMPOOL_MODEL                     /* Pool Model (proc,thrd,both)*/
	attrModes                         ociAttrType = C.OCI_ATTR_MODES                             /* Modes */
	attrSubscrName                    ociAttrType = C.OCI_ATTR_SUBSCR_NAME                       /* name of subscription */
	attrSubscrCallback                ociAttrType = C.OCI_ATTR_SUBSCR_CALLBACK                   /* associated callback */
	attrSubscrCtx                     ociAttrType = C.OCI_ATTR_SUBSCR_CTX                        /* associated callback context */
	attrSubscrPayload                 ociAttrType = C.OCI_ATTR_SUBSCR_PAYLOAD                    /* associated payload */
	attrSubscrNamespace               ociAttrType = C.OCI_ATTR_SUBSCR_NAMESPACE                  /* associated namespace */
	attrProxyCredentials              ociAttrType = C.OCI_ATTR_PROXY_CREDENTIALS                 /* Proxy user credentials */
	attrInitialClientRoles            ociAttrType = C.OCI_ATTR_INITIAL_CLIENT_ROLES              /* Initial client role list */
	attrUnk                           ociAttrType = C.OCI_ATTR_UNK                               /* unknown attribute */
	attrNumCols                       ociAttrType = C.OCI_ATTR_NUM_COLS                          /* number of columns */
	attrListColumns                   ociAttrType = C.OCI_ATTR_LIST_COLUMNS                      /* parameter of the column list */
	attrRdba                          ociAttrType = C.OCI_ATTR_RDBA                              /* DBA of the segment header */
	attrClustered                     ociAttrType = C.OCI_ATTR_CLUSTERED                         /* whether the table is clustered */
	attrPartitioned                   ociAttrType = C.OCI_ATTR_PARTITIONED                       /* whether the table is partitioned */
	attrIndexOnly                     ociAttrType = C.OCI_ATTR_INDEX_ONLY                        /* whether the table is index only */
	attrListArguments                 ociAttrType = C.OCI_ATTR_LIST_ARGUMENTS                    /* parameter of the argument list */
	attrListSubprograms               ociAttrType = C.OCI_ATTR_LIST_SUBPROGRAMS                  /* parameter of the subprogram list */
	attrRefTypeDescriptor             ociAttrType = C.OCI_ATTR_REF_TDO                           /* REF to the type descriptor */
	attrLink                          ociAttrType = C.OCI_ATTR_LINK                              /* the database link name */
	attrMin                           ociAttrType = C.OCI_ATTR_MIN                               /* minimum value */
	attrMax                           ociAttrType = C.OCI_ATTR_MAX                               /* maximum value */
	attrIncr                          ociAttrType = C.OCI_ATTR_INCR                              /* increment value */
	attrCache                         ociAttrType = C.OCI_ATTR_CACHE                             /* number of sequence numbers cached */
	attrOrder                         ociAttrType = C.OCI_ATTR_ORDER                             /* whether the sequence is ordered */
	attrHighWaterMark                 ociAttrType = C.OCI_ATTR_HW_MARK                           /* high-water mark */
	attrTypeSchema                    ociAttrType = C.OCI_ATTR_TYPE_SCHEMA                       /* type's schema name */
	attrTimestamp                     ociAttrType = C.OCI_ATTR_TIMESTAMP                         /* timestamp of the object */
	attrNumAttrs                      ociAttrType = C.OCI_ATTR_NUM_ATTRS                         /* number of sttributes */
	attrNumParams                     ociAttrType = C.OCI_ATTR_NUM_PARAMS                        /* number of parameters */
	attrTableObjID                    ociAttrType = C.OCI_ATTR_OBJID                             /* object id for a table or view */
	attrPtype                         ociAttrType = C.OCI_ATTR_PTYPE                             /* type of info described by */
	attrParam                         ociAttrType = C.OCI_ATTR_PARAM                             /* parameter descriptor */
	attrOverloadID                    ociAttrType = C.OCI_ATTR_OVERLOAD_ID                       /* overload ID for funcs and procs */
	attrTablespace                    ociAttrType = C.OCI_ATTR_TABLESPACE                        /* table name space */
	attrTDO                           ociAttrType = C.OCI_ATTR_TDO                               /* TDO of a type */
	attrListType                      ociAttrType = C.OCI_ATTR_LTYPE                             /* list type */
	attrParseErrorOffset              ociAttrType = C.OCI_ATTR_PARSE_ERROR_OFFSET                /* Parse Error offset */
	attrIsTemporary                   ociAttrType = C.OCI_ATTR_IS_TEMPORARY                      /* whether table is temporary */
	attrIsTyped                       ociAttrType = C.OCI_ATTR_IS_TYPED                          /* whether table is typed */
	attrDuration                      ociAttrType = C.OCI_ATTR_DURATION                          /* duration of temporary table */
	attrIsInvokerRights               ociAttrType = C.OCI_ATTR_IS_INVOKER_RIGHTS                 /* is invoker rights */
	attrObjName                       ociAttrType = C.OCI_ATTR_OBJ_NAME                          /* top level schema obj name */
	attrObjSchema                     ociAttrType = C.OCI_ATTR_OBJ_SCHEMA                        /* schema name */
	attrSchemaObjID                   ociAttrType = C.OCI_ATTR_OBJ_ID                            /* top level schema object id */
	attrTransTimeout                  ociAttrType = C.OCI_ATTR_TRANS_TIMEOUT                     /* transaction timeout */
	attrServerStatus                  ociAttrType = C.OCI_ATTR_SERVER_STATUS                     /* state of the server handle */
	attrStatement                     ociAttrType = C.OCI_ATTR_STATEMENT                         /* statement txt in stmt hdl */
	attrDeqCondition                  ociAttrType = C.OCI_ATTR_DEQCOND                           /* dequeue condition */
	attrSubscrRecpt                   ociAttrType = C.OCI_ATTR_SUBSCR_RECPT                      /* recepient of subscription */
	attrSubscrRecptProtocol           ociAttrType = C.OCI_ATTR_SUBSCR_RECPTPROTO                 /* protocol for recepient */
	attrLdapHost                      ociAttrType = C.OCI_ATTR_LDAP_HOST                         /* LDAP host to connect to */
	attrLdapPort                      ociAttrType = C.OCI_ATTR_LDAP_PORT                         /* LDAP port to connect to */
	attrBindDN                        ociAttrType = C.OCI_ATTR_BIND_DN                           /* bind DN */
	attrLdapCred                      ociAttrType = C.OCI_ATTR_LDAP_CRED                         /* credentials to connect to LDAP */
	attrWalletLocation                ociAttrType = C.OCI_ATTR_WALL_LOC                          /* client wallet location */
	attrLdapAuth                      ociAttrType = C.OCI_ATTR_LDAP_AUTH                         /* LDAP authentication method */
	attrLdapContext                   ociAttrType = C.OCI_ATTR_LDAP_CTX                          /* LDAP adminstration context DN */
	attrServerDNS                     ociAttrType = C.OCI_ATTR_SERVER_DNS                        /* list of registration server DNs */
	attrDNCount                       ociAttrType = C.OCI_ATTR_DN_COUNT                          /* the number of server DNs */
	attrServerDN                      ociAttrType = C.OCI_ATTR_SERVER_DN                         /* server DN attribute */
	attrMaxcharSize                   ociAttrType = C.OCI_ATTR_MAXCHAR_SIZE                      /* max char size of data */
	attrCurrentPosition               ociAttrType = C.OCI_ATTR_CURRENT_POSITION                  /* for scrollable result sets*/
	attrDigestAlgo                    ociAttrType = C.OCI_ATTR_DIGEST_ALGO                       /* digest algorithm */
	attrCertificate                   ociAttrType = C.OCI_ATTR_CERTIFICATE                       /* certificate */
	attrSignatureAlgo                 ociAttrType = C.OCI_ATTR_SIGNATURE_ALGO                    /* signature algorithm */
	attrCanonicalAlgo                 ociAttrType = C.OCI_ATTR_CANONICAL_ALGO                    /* canonicalization algo. */
	attrPrivateKey                    ociAttrType = C.OCI_ATTR_PRIVATE_KEY                       /* private key */
	attrDigestValue                   ociAttrType = C.OCI_ATTR_DIGEST_VALUE                      /* digest value */
	attrSignatureVal                  ociAttrType = C.OCI_ATTR_SIGNATURE_VAL                     /* signature value */
	attrSignature                     ociAttrType = C.OCI_ATTR_SIGNATURE                         /* signature */
	attrStmtcachesize                 ociAttrType = C.OCI_ATTR_STMTCACHESIZE                     /* size of the stm cache */
	attrConnNoWait                    ociAttrType = C.OCI_ATTR_CONN_NOWAIT                       /* connection pool attributes */
	attrConnBusyCount                 ociAttrType = C.OCI_ATTR_CONN_BUSY_COUNT                   /* connection pool attributes */
	attrConnOpenCount                 ociAttrType = C.OCI_ATTR_CONN_OPEN_COUNT                   /* connection pool attributes */
	attrConnTimeout                   ociAttrType = C.OCI_ATTR_CONN_TIMEOUT                      /* connection pool attributes */
	attrStmtState                     ociAttrType = C.OCI_ATTR_STMT_STATE                        /* connection pool attributes */
	attrConnMin                       ociAttrType = C.OCI_ATTR_CONN_MIN                          /* connection pool attributes */
	attrConnMax                       ociAttrType = C.OCI_ATTR_CONN_MAX                          /* connection pool attributes */
	attrConnIncr                      ociAttrType = C.OCI_ATTR_CONN_INCR                         /* connection pool attributes */
	attrNumOpenStmts                  ociAttrType = C.OCI_ATTR_NUM_OPEN_STMTS                    /* open stmts in session */
	attrDescribeNative                ociAttrType = C.OCI_ATTR_DESCRIBE_NATIVE                   /* get native info via desc */
	attrBindCount                     ociAttrType = C.OCI_ATTR_BIND_COUNT                        /* number of bind postions */
	attrHandlePosition                ociAttrType = C.OCI_ATTR_HANDLE_POSITION                   /* pos of bind/define handle */
	attrServerBusy                    ociAttrType = C.OCI_ATTR_SERVER_BUSY                       /* call in progress on server*/
	attrSubscrRecptpres               ociAttrType = C.OCI_ATTR_SUBSCR_RECPTPRES                  /* notification presentation for recipient */
	attrTransformation                ociAttrType = C.OCI_ATTR_TRANSFORMATION                    /* AQ message transformation */
	attrRowsFetched                   ociAttrType = C.OCI_ATTR_ROWS_FETCHED                      /* rows fetched in last call */
	attrScnBase                       ociAttrType = C.OCI_ATTR_SCN_BASE                          /* snapshot base */
	attrScnWrap                       ociAttrType = C.OCI_ATTR_SCN_WRAP                          /* snapshot wrap */
	attrReadonlyTxn                   ociAttrType = C.OCI_ATTR_READONLY_TXN                      /* txn is readonly */
	attrErroneousColumn               ociAttrType = C.OCI_ATTR_ERRONEOUS_COLUMN                  /* position of erroneous col */
	attrAsmVolSprt                    ociAttrType = C.OCI_ATTR_ASM_VOL_SPRT                      /* ASM volume supported? */
	attrInstType                      ociAttrType = C.OCI_ATTR_INST_TYPE                         /* oracle instance type */
	attrEnvUtf16                      ociAttrType = C.OCI_ATTR_ENV_UTF16                         /* is env in utf16 mode? */
	attrIsExternal                    ociAttrType = C.OCI_ATTR_IS_EXTERNAL                       /* whether table is external */
	attrStmtIsReturning               ociAttrType = C.OCI_ATTR_STMT_IS_RETURNING                 /* stmt has returning clause */
	attrCurrentSchema                 ociAttrType = C.OCI_ATTR_CURRENT_SCHEMA                    /* Current Schema */
	attrSubscrQosFlags                ociAttrType = C.OCI_ATTR_SUBSCR_QOSFLAGS                   /* QOS flags */
	attrSubscrPayloadCallback         ociAttrType = C.OCI_ATTR_SUBSCR_PAYLOADCBK                 /* Payload callback */
	attrSubscrTimeout                 ociAttrType = C.OCI_ATTR_SUBSCR_TIMEOUT                    /* Timeout */
	attrSubscrNamespaceCtx            ociAttrType = C.OCI_ATTR_SUBSCR_NAMESPACE_CTX              /* Namespace context */
	attrSubscrCqQosFlags              ociAttrType = C.OCI_ATTR_SUBSCR_CQ_QOSFLAGS                /* change notification (CQ) specific QOS flags */
	attrSubscrCqRegID                 ociAttrType = C.OCI_ATTR_SUBSCR_CQ_REGID                   /* change notification registration id */
	attrSubscrNtfnGroupingClass       ociAttrType = C.OCI_ATTR_SUBSCR_NTFN_GROUPING_CLASS        /* ntfn grouping class */
	attrSubscrNtfnGroupingValue       ociAttrType = C.OCI_ATTR_SUBSCR_NTFN_GROUPING_VALUE        /* ntfn grouping value */
	attrSubscrNtfnGroupingType        ociAttrType = C.OCI_ATTR_SUBSCR_NTFN_GROUPING_TYPE         /* ntfn grouping type */
	attrSubscrNtfnGroupingStartTime   ociAttrType = C.OCI_ATTR_SUBSCR_NTFN_GROUPING_START_TIME   /* ntfn grp start time */
	attrSubscrNtfnGroupingRepeatCount ociAttrType = C.OCI_ATTR_SUBSCR_NTFN_GROUPING_REPEAT_COUNT /* ntfn grp rep count */
	attrAQNtfnGroupingMsgIDArray      ociAttrType = C.OCI_ATTR_AQ_NTFN_GROUPING_MSGID_ARRAY      /* aq grp msgid array */
	attrAqNtfnGroupingCount           ociAttrType = C.OCI_ATTR_AQ_NTFN_GROUPING_COUNT            /* ntfns recd in grp */
	attrBindRowCallback               ociAttrType = C.OCI_ATTR_BIND_ROWCBK                       /* bind row callback */
	attrBindRowCtx                    ociAttrType = C.OCI_ATTR_BIND_ROWCTX                       /* ctx for bind row callback */
	attrSkipBuffer                    ociAttrType = C.OCI_ATTR_SKIP_BUFFER                       /* skip buffer in array ops */
	attrXStreamAckInterval            ociAttrType = C.OCI_ATTR_XSTREAM_ACK_INTERVAL              /* XStream ack interval */
	attrXStreamIdleTimeout            ociAttrType = C.OCI_ATTR_XSTREAM_IDLE_TIMEOUT              /* XStream idle timeout */
	attrCqQueryID                     ociAttrType = C.OCI_ATTR_CQ_QUERYID                        /**/
	attrChnfTablenames                ociAttrType = C.OCI_ATTR_CHNF_TABLENAMES                   /* out: array of table names   */
	attrChnfRowIDs                    ociAttrType = C.OCI_ATTR_CHNF_ROWIDS                       /* in: rowids needed */
	attrChnfOperations                ociAttrType = C.OCI_ATTR_CHNF_OPERATIONS                   /* in: notification operation filter*/
	attrChnfChangeLag                 ociAttrType = C.OCI_ATTR_CHNF_CHANGELAG                    /* txn lag between notifications  */
	attrChdesDBName                   ociAttrType = C.OCI_ATTR_CHDES_DBNAME                      /* source database    */
	attrChdesNotifyType               ociAttrType = C.OCI_ATTR_CHDES_NFYTYPE                     /* notification type flags */
	attrChdesXID                      ociAttrType = C.OCI_ATTR_CHDES_XID                         /* XID  of the transaction */
	attrChdesTableChanges             ociAttrType = C.OCI_ATTR_CHDES_TABLE_CHANGES               /* array of table chg descriptors*/
	attrChdesTableName                ociAttrType = C.OCI_ATTR_CHDES_TABLE_NAME                  /* table name */
	attrChdesTableOpflags             ociAttrType = C.OCI_ATTR_CHDES_TABLE_OPFLAGS               /* table operation flags */
	attrChdesTableRowChanges          ociAttrType = C.OCI_ATTR_CHDES_TABLE_ROW_CHANGES           /* array of changed rows   */
	attrChdesRowRowid                 ociAttrType = C.OCI_ATTR_CHDES_ROW_ROWID                   /* rowid of changed row    */
	attrChdesRowOpflags               ociAttrType = C.OCI_ATTR_CHDES_ROW_OPFLAGS                 /* row operation flags     */
	attrChnfReghandle                 ociAttrType = C.OCI_ATTR_CHNF_REGHANDLE                    /* IN: subscription handle  */
	attrNetworkFileDesc               ociAttrType = C.OCI_ATTR_NETWORK_FILE_DESC                 /* network file descriptor */
	attrProxyClient                   ociAttrType = C.OCI_ATTR_PROXY_CLIENT                      /**/
	attrTableEnc                      ociAttrType = C.OCI_ATTR_TABLE_ENC                         /* does table have any encrypt columns */
	attrTableEncAlg                   ociAttrType = C.OCI_ATTR_TABLE_ENC_ALG                     /* Table encryption Algorithm */
	attrTableEncAlgID                 ociAttrType = C.OCI_ATTR_TABLE_ENC_ALG_ID                  /* Internal Id of encryption Algorithm*/
	attrStmtcacheCbkctx               ociAttrType = C.OCI_ATTR_STMTCACHE_CBKCTX                  /* opaque context on stmt */
	attrStmtcacheCbk                  ociAttrType = C.OCI_ATTR_STMTCACHE_CBK                     /* callback fn for stmtcache */
	attrCqdesOperation                ociAttrType = C.OCI_ATTR_CQDES_OPERATION                   /**/
	attrCqdesTableChanges             ociAttrType = C.OCI_ATTR_CQDES_TABLE_CHANGES               /**/
	attrCqdesQueryID                  ociAttrType = C.OCI_ATTR_CQDES_QUERYID                     /**/
	attrChdesQueries                  ociAttrType = C.OCI_ATTR_CHDES_QUERIES                     /* Top level change desc array of queries */
	attrConnectionClass               ociAttrType = C.OCI_ATTR_CONNECTION_CLASS                  /* server-side session pool support */
	attrPurity                        ociAttrType = C.OCI_ATTR_PURITY                            /* server-side session pool support */
	attrPurityDefault                 ociAttrType = C.OCI_ATTR_PURITY_DEFAULT                    /* purity support */
	attrPurityNew                     ociAttrType = C.OCI_ATTR_PURITY_NEW                        /* purity support */
	attrPuritySelf                    ociAttrType = C.OCI_ATTR_PURITY_SELF                       /* purity support */
	attrSendTimeout                   ociAttrType = C.OCI_ATTR_SEND_TIMEOUT                      /* NS send timeout */
	attrReceiveTimeout                ociAttrType = C.OCI_ATTR_RECEIVE_TIMEOUT                   /* NS receive timeout */
	attrDefaultLobPrefetchSize        ociAttrType = C.OCI_ATTR_DEFAULT_LOBPREFETCH_SIZE          /* default prefetch size */
	attrLobPrefetchSize               ociAttrType = C.OCI_ATTR_LOBPREFETCH_SIZE                  /* prefetch size */
	attrLobPrefetchLength             ociAttrType = C.OCI_ATTR_LOBPREFETCH_LENGTH                /* prefetch length & chunk */
	attrLobRegionPrimary              ociAttrType = C.OCI_ATTR_LOB_REGION_PRIMARY                /* Primary LOB Locator */
	attrLobRegionPrimoff              ociAttrType = C.OCI_ATTR_LOB_REGION_PRIMOFF                /* Offset into Primary LOB */
	attrLobRegionOffset               ociAttrType = C.OCI_ATTR_LOB_REGION_OFFSET                 /* Region Offset */
	attrLobRegionLength               ociAttrType = C.OCI_ATTR_LOB_REGION_LENGTH                 /* Region Length Bytes/Chars */
	attrLobRegionMime                 ociAttrType = C.OCI_ATTR_LOB_REGION_MIME                   /* Region mime type */
	attrFetchRowID                    ociAttrType = C.OCI_ATTR_FETCH_ROWID                       /* fetch rowid */
	attrSubscrIpAddr                  ociAttrType = C.OCI_ATTR_SUBSCR_IPADDR                     /* ip address to listen on  */
	attrEnvCharsetID                  ociAttrType = C.OCI_ATTR_ENV_CHARSET_ID                    /* charset id in env */
	attrEnvNcharsetID                 ociAttrType = C.OCI_ATTR_ENV_NCHARSET_ID                   /* ncharset id in env */
	attrEvtCallback                   ociAttrType = C.OCI_ATTR_EVTCBK                            /* ha callback */
	attrEvtCtx                        ociAttrType = C.OCI_ATTR_EVTCTX                            /* ctx for ha callback */
	attrUserMemory                    ociAttrType = C.OCI_ATTR_USER_MEMORY                       /* pointer to user memory */
	attrAccessBanner                  ociAttrType = C.OCI_ATTR_ACCESS_BANNER                     /* access banner */
	attrAuditBanner                   ociAttrType = C.OCI_ATTR_AUDIT_BANNER                      /* audit banner */
	attrSubscrPortno                  ociAttrType = C.OCI_ATTR_SUBSCR_PORTNO                     /* port no to listen        */
	attrSessPoolTimeout               ociAttrType = C.OCI_ATTR_SPOOL_TIMEOUT                     /* session timeout */
	attrSessPoolGetMode               ociAttrType = C.OCI_ATTR_SPOOL_GETMODE                     /* session get mode */
	attrSessPoolBusyCount             ociAttrType = C.OCI_ATTR_SPOOL_BUSY_COUNT                  /* busy session count */
	attrSessPoolOpenCount             ociAttrType = C.OCI_ATTR_SPOOL_OPEN_COUNT                  /* open session count */
	attrSessPoolMin                   ociAttrType = C.OCI_ATTR_SPOOL_MIN                         /* min session count */
	attrSessPoolMax                   ociAttrType = C.OCI_ATTR_SPOOL_MAX                         /* max session count */
	attrSessPoolIncr                  ociAttrType = C.OCI_ATTR_SPOOL_INCR                        /* session increment count */
	attrSessPoolStmtCacheSize         ociAttrType = C.OCI_ATTR_SPOOL_STMTCACHESIZE               /*Stmt cache size of pool  */
	attrSessPoolAuth                  ociAttrType = C.OCI_ATTR_SPOOL_AUTH                        /* Auth handle on pool handle*/
	attrDataSize                      ociAttrType = C.OCI_ATTR_DATA_SIZE                         /* maximum size of the data */
	attrDataType                      ociAttrType = C.OCI_ATTR_DATA_TYPE                         /* the SQL type of the column/argument */
	attrDispSize                      ociAttrType = C.OCI_ATTR_DISP_SIZE                         /* the display size */
	attrName                          ociAttrType = C.OCI_ATTR_NAME                              /* the name of the column/argument */
	attrPrecision                     ociAttrType = C.OCI_ATTR_PRECISION                         /* precision if number type */
	attrScale                         ociAttrType = C.OCI_ATTR_SCALE                             /* scale if number type */
	attrIsNull                        ociAttrType = C.OCI_ATTR_IS_NULL                           /* is it null ? */
	attrTypeName                      ociAttrType = C.OCI_ATTR_TYPE_NAME                         /* name of the named data type or a package name for package private types */
	attrSchemaName                    ociAttrType = C.OCI_ATTR_SCHEMA_NAME                       /* the schema name */
	attrSubName                       ociAttrType = C.OCI_ATTR_SUB_NAME                          /* type name if package private type */
	attrPosition                      ociAttrType = C.OCI_ATTR_POSITION                          /* relative position of col/arg in the list of cols/args */
	attrComplexobjectcompType         ociAttrType = C.OCI_ATTR_COMPLEXOBJECTCOMP_TYPE
	attrComplexobjectcompTypeLevel    ociAttrType = C.OCI_ATTR_COMPLEXOBJECTCOMP_TYPE_LEVEL
	attrComplexobjectLevel            ociAttrType = C.OCI_ATTR_COMPLEXOBJECT_LEVEL
	attrComplexobjectCollOutofline    ociAttrType = C.OCI_ATTR_COMPLEXOBJECT_COLL_OUTOFLINE
	attrDispName                      ociAttrType = C.OCI_ATTR_DISP_NAME                /* the display name */
	attrEnccSize                      ociAttrType = C.OCI_ATTR_ENCC_SIZE                /* encrypted data size */
	attrColEnc                        ociAttrType = C.OCI_ATTR_COL_ENC                  /* column is encrypted ? */
	attrColEncSalt                    ociAttrType = C.OCI_ATTR_COL_ENC_SALT             /* is encrypted column salted ? */
	attrOverload                      ociAttrType = C.OCI_ATTR_OVERLOAD                 /* is this position overloaded */
	attrLevel                         ociAttrType = C.OCI_ATTR_LEVEL                    /* level for structured types */
	attrHasDefault                    ociAttrType = C.OCI_ATTR_HAS_DEFAULT              /* has a default value */
	attrIomode                        ociAttrType = C.OCI_ATTR_IOMODE                   /* in, out inout */
	attrRadix                         ociAttrType = C.OCI_ATTR_RADIX                    /* returns a radix */
	attrNumArgs                       ociAttrType = C.OCI_ATTR_NUM_ARGS                 /* total number of arguments */
	attrTypecode                      ociAttrType = C.OCI_ATTR_TYPECODE                 /* object or collection */
	attrCollectionTypecode            ociAttrType = C.OCI_ATTR_COLLECTION_TYPECODE      /* varray or nested table */
	attrVersion                       ociAttrType = C.OCI_ATTR_VERSION                  /* user assigned version */
	attrIsIncompleteType              ociAttrType = C.OCI_ATTR_IS_INCOMPLETE_TYPE       /* is this an incomplete type */
	attrIsSystemType                  ociAttrType = C.OCI_ATTR_IS_SYSTEM_TYPE           /* a system type */
	attrIsPredefinedType              ociAttrType = C.OCI_ATTR_IS_PREDEFINED_TYPE       /* a predefined type */
	attrIsTransientType               ociAttrType = C.OCI_ATTR_IS_TRANSIENT_TYPE        /* a transient type */
	attrIsSystemGeneratedType         ociAttrType = C.OCI_ATTR_IS_SYSTEM_GENERATED_TYPE /* system generated type */
	attrHasNestedTable                ociAttrType = C.OCI_ATTR_HAS_NESTED_TABLE         /* contains nested table attr */
	attrHasLob                        ociAttrType = C.OCI_ATTR_HAS_LOB                  /* has a lob attribute */
	attrHasFile                       ociAttrType = C.OCI_ATTR_HAS_FILE                 /* has a file attribute */
	attrCollectionElement             ociAttrType = C.OCI_ATTR_COLLECTION_ELEMENT       /* has a collection attribute */
	attrNumTypeAttrs                  ociAttrType = C.OCI_ATTR_NUM_TYPE_ATTRS           /* number of attribute types */
	attrListTypeAttrs                 ociAttrType = C.OCI_ATTR_LIST_TYPE_ATTRS          /* list of type attributes */
	attrNumTypeMethods                ociAttrType = C.OCI_ATTR_NUM_TYPE_METHODS         /* number of type methods */
	attrListTypeMethods               ociAttrType = C.OCI_ATTR_LIST_TYPE_METHODS        /* list of type methods */
	attrMapMethod                     ociAttrType = C.OCI_ATTR_MAP_METHOD               /* map method of type */
	attrOrderMethod                   ociAttrType = C.OCI_ATTR_ORDER_METHOD             /* order method of type */
	attrNumElems                      ociAttrType = C.OCI_ATTR_NUM_ELEMS                /* number of elements */
	attrEncapsulation                 ociAttrType = C.OCI_ATTR_ENCAPSULATION            /* encapsulation level */
	attrIsSelfish                     ociAttrType = C.OCI_ATTR_IS_SELFISH               /* method selfish */
	attrIsVirtual                     ociAttrType = C.OCI_ATTR_IS_VIRTUAL               /* virtual */
	attrIsInline                      ociAttrType = C.OCI_ATTR_IS_INLINE                /* inline */
	attrIsConstant                    ociAttrType = C.OCI_ATTR_IS_CONSTANT              /* constant */
	attrHasResult                     ociAttrType = C.OCI_ATTR_HAS_RESULT               /* has result */
	attrIsConstructor                 ociAttrType = C.OCI_ATTR_IS_CONSTRUCTOR           /* constructor */
	attrIsDestructor                  ociAttrType = C.OCI_ATTR_IS_DESTRUCTOR            /* destructor */
	attrIsOperator                    ociAttrType = C.OCI_ATTR_IS_OPERATOR              /* operator */
	attrIsMap                         ociAttrType = C.OCI_ATTR_IS_MAP                   /* a map method */
	attrIsOrder                       ociAttrType = C.OCI_ATTR_IS_ORDER                 /* order method */
	attrIsRnds                        ociAttrType = C.OCI_ATTR_IS_RNDS                  /* read no data state method */
	attrIsRnps                        ociAttrType = C.OCI_ATTR_IS_RNPS                  /* read no process state */
	attrIsWnds                        ociAttrType = C.OCI_ATTR_IS_WNDS                  /* write no data state method */
	attrIsWnps                        ociAttrType = C.OCI_ATTR_IS_WNPS                  /* write no process state */
	attrDescPublic                    ociAttrType = C.OCI_ATTR_DESC_PUBLIC              /* public object */
	attrCacheClientContext            ociAttrType = C.OCI_ATTR_CACHE_CLIENT_CONTEXT
	attrUciConstruct                  ociAttrType = C.OCI_ATTR_UCI_CONSTRUCT
	attrUciDestruct                   ociAttrType = C.OCI_ATTR_UCI_DESTRUCT
	attrUciCopy                       ociAttrType = C.OCI_ATTR_UCI_COPY
	attrUciPickle                     ociAttrType = C.OCI_ATTR_UCI_PICKLE
	attrUciUnpickle                   ociAttrType = C.OCI_ATTR_UCI_UNPICKLE
	attrUciRefresh                    ociAttrType = C.OCI_ATTR_UCI_REFRESH
	attrIsSubtype                     ociAttrType = C.OCI_ATTR_IS_SUBTYPE
	attrSupertypeSchemaName           ociAttrType = C.OCI_ATTR_SUPERTYPE_SCHEMA_NAME
	attrSupertypeName                 ociAttrType = C.OCI_ATTR_SUPERTYPE_NAME
	attrListObjects                   ociAttrType = C.OCI_ATTR_LIST_OBJECTS           /* list of objects in schema */
	attrNcharsetId                    ociAttrType = C.OCI_ATTR_NCHARSET_ID            /* char set id */
	attrListSchemas                   ociAttrType = C.OCI_ATTR_LIST_SCHEMAS           /* list of schemas */
	attrMaxProcLen                    ociAttrType = C.OCI_ATTR_MAX_PROC_LEN           /* max procedure length */
	attrMaxColumnLen                  ociAttrType = C.OCI_ATTR_MAX_COLUMN_LEN         /* max column name length */
	attrCursorCommitBehavior          ociAttrType = C.OCI_ATTR_CURSOR_COMMIT_BEHAVIOR /* cursor commit behavior */
	attrMaxCatalogNamelen             ociAttrType = C.OCI_ATTR_MAX_CATALOG_NAMELEN    /* catalog namelength */
	attrCatalogLocation               ociAttrType = C.OCI_ATTR_CATALOG_LOCATION       /* catalog location */
	attrSavepointSupport              ociAttrType = C.OCI_ATTR_SAVEPOINT_SUPPORT      /* savepoint support */
	attrNowaitSupport                 ociAttrType = C.OCI_ATTR_NOWAIT_SUPPORT         /* nowait support */
	attrAutocommitDdl                 ociAttrType = C.OCI_ATTR_AUTOCOMMIT_DDL         /* autocommit DDL */
	attrLockingMode                   ociAttrType = C.OCI_ATTR_LOCKING_MODE           /* locking mode */
	attrAppctxSize                    ociAttrType = C.OCI_ATTR_APPCTX_SIZE            /* count of context to be init*/
	attrAppctxList                    ociAttrType = C.OCI_ATTR_APPCTX_LIST            /* count of context to be init*/
	attrAppctxName                    ociAttrType = C.OCI_ATTR_APPCTX_NAME            /* name  of context to be init*/
	attrAppctxAttr                    ociAttrType = C.OCI_ATTR_APPCTX_ATTR            /* attr  of context to be init*/
	attrAppctxValue                   ociAttrType = C.OCI_ATTR_APPCTX_VALUE           /* value of context to be init*/
	attrClientIdentifier              ociAttrType = C.OCI_ATTR_CLIENT_IDENTIFIER      /* value of client id to set*/
	attrIsFinalType                   ociAttrType = C.OCI_ATTR_IS_FINAL_TYPE          /* is final type ? */
	attrIsInstantiableType            ociAttrType = C.OCI_ATTR_IS_INSTANTIABLE_TYPE   /* is instantiable type ? */
	attrIsFinalMethod                 ociAttrType = C.OCI_ATTR_IS_FINAL_METHOD        /* is final method ? */
	attrIsInstantiableMethod          ociAttrType = C.OCI_ATTR_IS_INSTANTIABLE_METHOD /* is instantiable method ? */
	attrIsOverridingMethod            ociAttrType = C.OCI_ATTR_IS_OVERRIDING_METHOD   /* is overriding method ? */
	attrDescSynbase                   ociAttrType = C.OCI_ATTR_DESC_SYNBASE           /* Describe the base object */
	attrCharUsed                      ociAttrType = C.OCI_ATTR_CHAR_USED              /* char length semantics */
	attrCharSize                      ociAttrType = C.OCI_ATTR_CHAR_SIZE              /* char length */
)

func ociAttrGetString(handle unsafe.Pointer, htype ociHandleType, attrType ociAttrType, errHandle *C.OCIError) (rslt string, err *OciError) {
	var pstr *C.char
	var sz1 C.ub4

	err = checkError(C.OCIAttrGet(handle, (C.ub4)(htype), (unsafe.Pointer)(&pstr), &sz1, (C.ub4)(attrType), errHandle), errHandle)

	if err != nil {
		return "", err
	}

	rslt = C.GoStringN(pstr, (C.int)(sz1))

	return

}

func ociAttrGet(handle unsafe.Pointer, htype ociHandleType, attrPtr unsafe.Pointer, attrSize *C.ub4, attrType ociAttrType, errHandle *C.OCIError) *OciError {
	return checkError(C.OCIAttrGet(handle, (C.ub4)(htype), attrPtr, attrSize, (C.ub4)(attrType), errHandle), errHandle)
}

func ociAttrSet(handle unsafe.Pointer, htype ociHandleType, attrPtr unsafe.Pointer, attrSize C.ub4, attrType ociAttrType, errHandle *C.OCIError) *OciError {
	return checkError(C.OCIAttrSet(handle, (C.ub4)(htype), attrPtr, attrSize, (C.ub4)(attrType), errHandle), errHandle)
}
