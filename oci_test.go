package oci_test

import (
	"fmt"
	"github.com/djbckr/ocigo"
	"os"
	"testing"
	"time"
)

func TestOCI(t *testing.T) {

	connstring := os.Getenv("CONNECT_STRING")

	if connstring == "" {
		t.Fatal("To run a test, create an environment variable called \n\"CONNECT_STRING\" formatted like \"username/password@host/dbname\"")
	}

	fmt.Println("Creating Pool")
	pool, err := oci.CreatePool(connstring, 1, 5, 1)
	checkerr(t, err)
	defer pool.Destroy()

	fmt.Println("Getting Session")
	ses, err := pool.Acquire()
	checkerr(t, err)
	defer ses.Release()

	fmt.Println("Dropping table...")
	execSql(ses, t, "drop table foo cascade constraints")

	fmt.Println("Creating table...")

	execSql(ses, t, `
create table foo (
  "AnyData"      ANYDATA,
  "BinaryDbl"    BINARY_DOUBLE,
  "BinaryFlt"    BINARY_FLOAT,
  "Blob"         BLOB,
  "CharC"        CHAR(10 char),
  "CharB"        CHAR(10 byte),
  "Clob"         CLOB,
  "Date"         DATE,
  "Float"        FLOAT,
  "IntervalDS"   INTERVAL DAY TO SECOND,
  "IntervalYM"   INTERVAL YEAR TO MONTH,
  "Long"         LONG,
  "NChar"        NCHAR(100),
  "NClob"        NCLOB,
  "Number"       NUMBER,
  "NVarchar2"    NVARCHAR2(100),
  "Raw"          RAW(100),
  "RowID"        ROWID,
  "TimeStamp"    TIMESTAMP,
  "TimeStampTZ"  TIMESTAMP WITH TIME ZONE,
  "TimeStampLTZ" TIMESTAMP WITH LOCAL TIME ZONE,
  "URowID"       UROWID,
  "VarChar2C"    VARCHAR2(100 char),
  "VarChar2B"    VARCHAR2(100 byte),
  "XmlType"      XMLTYPE
)`)

	execSql(ses, t, "alter session set nls_timestamp_format='YYYY-MM-DD HH24:MI:SS.FF4'")
	execSql(ses, t, "alter session set nls_timestamp_tz_format='YYYY-MM-DD HH24:MI:SSXFF TZR'")

	execSql(ses, t, `insert into foo values (
		anydata.convertVarchar2('hello world'), 1.13, 1.13, hextoraw('deadbeef'), 'こんにちは', 'helloworld', 'Mary had a little lamb...',
		sysdate, 1.13, numtodsinterval(1.27777777, 'day'), numtoyminterval(1.234, 'year'), 'Mary had a little lamb...', 'Mary had a little lamb...',
		'Mary had a little lamb...', 1.13, 'Mary had a little lamb...', hextoraw('deadbeef'), 'ABGVw8AHoAAAAFNAAA',
		systimestamp, systimestamp, systimestamp, 'ABGVw8AHoAAAAFNAAA', 'メリーさんの羊...', 'Mary had a little lamb...', xmltype('<rt><dat>hello world</dat></rt>')
		)`)

	execSql(ses, t, `insert into foo values (
		anydata.convertVarchar2('hello world'), 1.13, 1.13, hextoraw('deadbeef'), 'こんにちは', 'helloworld', 'Mary had a little lamb...',
		sysdate, 1.13, numtodsinterval(1.27777777, 'day'), numtoyminterval(1.234, 'year'), 'Mary had a little lamb...', 'Mary had a little lamb...',
		'Mary had a little lamb...', 1.13, 'Mary had a little lamb...', hextoraw('deadbeef'), 'ABGVw8AHoAAAAFNAAA',
		systimestamp, systimestamp, systimestamp, 'ABGVw8AHoAAAAFNAAA', 'メリーさんの羊...', 'Mary had a little lamb...', xmltype('<rt><dat>hello world</dat></rt>')
		)`)

	ses.Commit()

	fmt.Println("Running query...")
	querySql(ses, t, "select f.*, 'literal' as \"SomethingLiteral\" from foo f")

	n1, err := oci.NumberFromInt(2)
	checkerr(t, err)

	n2, err := oci.NumberFromInt(7)
	checkerr(t, err)

	n3, err := n1.Div(n2)
	checkerr(t, err)

	n4, err := n3.Round(2)
	checkerr(t, err)

	n5, err := n3.Trunc(2)
	checkerr(t, err)

	fmt.Println(n1)
	fmt.Println(n2)
	fmt.Println(n3)
	fmt.Println(n4)
	fmt.Println(n5)

	ts, err := ses.SysTimeStamp(oci.TypeTimestampTZ)
	checkerr(t, err)
	fmt.Println(ts)

	ts, err = ses.TimeStampFromGoTime(oci.TypeTimestampTZ, time.Now())
	checkerr(t, err)
	fmt.Println(ts)

	ses.Commit()

}

func checkerr(t *testing.T, e error) {
	if e != nil {
		panic(e)
	}
}

func logerr(t *testing.T, e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

func execSql(ses *oci.Session, t *testing.T, sql string) {
	stmt, err := ses.Prepare(sql)
	checkerr(t, err)
	defer stmt.Release(false)
	logerr(t, stmt.Execute())
}

func querySql(ses *oci.Session, t *testing.T, sql string) {
	stmt, err := ses.Prepare(sql)
	checkerr(t, err)
	defer stmt.Release(false)
	_, err = stmt.Query()
	checkerr(t, err)
}
