package oci

/*
#cgo pkg-config: oci
#include <oci.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

/*
   Support for Oracle Timestamp (with Time Zone/Local Time Zone)
   Support for Oracle Interval (Year to Month/Day to Second)

   The OCI structures are more complicated and require
   allocating/deallocating descriptors when used.
*/

// TimeStamp is an opaque structure that represents an Oracle TIMESTAMP [WITH [LOCAL] TIMEZONE]
type TimeStamp struct {
	ses      *Session
	err      *C.OCIError
	datetime *C.OCIDateTime
	tstype   TimestampType
}

// Interval is an opaque structure that represents an Oracle INTERVAL [YEAR TO MONTH]|[DAY TO SECOND]
type Interval struct {
	ses      *Session
	err      *C.OCIError
	interval *C.OCIInterval
	intype   IntervalType
}

// TimestampType clarifies the actual TimeStamp struct.
type TimestampType C.ub4

// timestamp types
const (
	TypeTimestamp    TimestampType = (TimestampType)(dtypeTimestamp)
	TypeTimestampTZ  TimestampType = (TimestampType)(dtypeTimestampTZ)
	TypeTimestampLTZ TimestampType = (TimestampType)(dtypeTimestampLTZ)
)

// IntervalType clarifies the actual Interval struct.
type IntervalType C.ub4

// interval types
const (
	TypeIntervalYM IntervalType = (IntervalType)(dtypeIntervalYM)
	TypeIntervalDS IntervalType = (IntervalType)(dtypeIntervalDS)
)

func finalizer(hndl unsafe.Pointer, errhndl unsafe.Pointer, typ C.ub4) {
	ociHandleFree(errhndl, htypeError)
	ociDescriptorFree(hndl, (ociDescriptorType)(typ))
}

func finalizerTimestamp(ts *TimeStamp) {
	finalizer((unsafe.Pointer)(ts.datetime), (unsafe.Pointer)(ts.err), (C.ub4)(ts.tstype))
}

func finalizerInterval(iv *Interval) {
	finalizer((unsafe.Pointer)(iv.interval), (unsafe.Pointer)(iv.err), (C.ub4)(iv.intype))
}

func makeTimestampInstance(s *Session, typ TimestampType) (rslt *TimeStamp) {
	rslt = &TimeStamp{ses: s, tstype: typ}
	ociDescriptorAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.datetime)), (ociDescriptorType)(typ))
	ociHandleAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.err)), htypeError)
	runtime.SetFinalizer(rslt, finalizerTimestamp)
	return
}

func makeIntervalInstance(s *Session, typ IntervalType) (rslt *Interval) {
	rslt = &Interval{ses: s, intype: typ}
	ociDescriptorAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.interval)), (ociDescriptorType)(typ))
	ociHandleAlloc((unsafe.Pointer)(genv), (*unsafe.Pointer)(unsafe.Pointer(&rslt.err)), htypeError)
	runtime.SetFinalizer(rslt, finalizerInterval)
	return
}

/*****************************************************************************/

// SysTimeStamp gets System Time Stamp based on Database Session settings
func (session *Session) SysTimeStamp(tstype TimestampType) (*TimeStamp, error) {

	rslt := makeTimestampInstance(session, tstype)

	err := checkError(
		C.OCIDateTimeSysTimeStamp(
			unsafe.Pointer(session.ses),
			rslt.err,
			rslt.datetime), rslt.err)

	return rslt, processError(err)

}

/*****************************************************************************/

// TimeStampFromText takes one or more strings and converts to a TimeStamp struct.
func (session *Session) TimeStampFromText(tstype TimestampType, params ...string) (*TimeStamp, error) {
	if len(params) == 0 {
		return nil, nil
	}

	dateStr := []byte(params[0])
	dstrLength := len(params[0])

	var fmt []byte
	var fmtp unsafe.Pointer
	var fmtLength C.ub1
	var langName []byte
	var langNameP unsafe.Pointer
	var langLength C.size_t

	if len(params) > 1 && len(params[1]) > 0 {
		fmt = []byte(params[1])
		fmtLength = (C.ub1)(len(params[1]))
		fmtp = (unsafe.Pointer)(&fmt[0])
	}

	if len(params) > 2 && len(params[2]) > 0 {
		langName = []byte(params[2])
		langLength = (C.size_t)(len(params[2]))
		langNameP = (unsafe.Pointer)(&langName[0])
	}

	rslt := makeTimestampInstance(session, tstype)

	err := checkError(
		C.OCIDateTimeFromText(
			unsafe.Pointer(session.ses),
			rslt.err,
			(*C.OraText)(unsafe.Pointer(&dateStr[0])),
			(C.size_t)(dstrLength),
			(*C.OraText)(fmtp),
			fmtLength,
			(*C.OraText)(langNameP),
			langLength,
			rslt.datetime), rslt.err)

	return rslt, processError(err)

}

/*****************************************************************************/

// TimeStampFromGoTime converts a go time to an Oracle TimeStamp
func (session *Session) TimeStampFromGoTime(tstype TimestampType, t time.Time) (*TimeStamp, error) {

	var year = int16(t.Year())
	var month = uint8(t.Month())
	var day = uint8(t.Day())
	var hour = uint8(t.Hour())
	var min = uint8(t.Minute())
	var sec = uint8(t.Second())
	var fsec = uint32(t.Nanosecond() / 1000000)
	_, offset := t.Zone()

	// offset is number of seconds; convert to a string that Oracle can interpret
	z1 := offset / 60 / 60
	z2 := (offset - (z1 * 60 * 60)) / 60
	if z2 < 0 {
		z2 = -z2
	}

	timezone := []byte(fmt.Sprintf("%+03d:%02d\n", z1, z2))

	rslt := makeTimestampInstance(session, tstype)

	err := checkError(
		C.OCIDateTimeConstruct(
			unsafe.Pointer(session.ses),
			rslt.err,
			rslt.datetime,
			(C.sb2)(year),
			(C.ub1)(month),
			(C.ub1)(day),
			(C.ub1)(hour),
			(C.ub1)(min),
			(C.ub1)(sec),
			(C.ub4)(fsec),
			(*C.OraText)(unsafe.Pointer(&timezone[0])), 6), rslt.err)

	return rslt, processError(err)

}

// GetDate extracts the year, month, and day values from a TimeStamp
func (ts *TimeStamp) GetDate() (year int16, month, day uint8) {

	err := checkError(
		C.OCIDateTimeGetDate(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			(*C.sb2)(unsafe.Pointer(&year)),
			(*C.ub1)(unsafe.Pointer(&month)),
			(*C.ub1)(unsafe.Pointer(&day))), ts.err)

	if err != nil {
		panic(err)
	}

	return
}

// GetTime extracts the hour, minute, second, and fracsecond values from a TimeStamp
func (ts *TimeStamp) GetTime() (hour, minute, second uint8, fracsecond uint32) {

	err := checkError(
		C.OCIDateTimeGetTime(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			(*C.ub1)(unsafe.Pointer(&hour)),
			(*C.ub1)(unsafe.Pointer(&minute)),
			(*C.ub1)(unsafe.Pointer(&second)),
			(*C.ub4)(unsafe.Pointer(&fracsecond))), ts.err)

	if err != nil {
		panic(err)
	}

	return
}

// GetTimeZoneName returns the name of the timezone (if there is one) from this TimeStamp
func (ts *TimeStamp) GetTimeZoneName() string {

	tzname := make([]byte, 128)
	var tznamelen uint32 = 128

	err := checkError(
		C.OCIDateTimeGetTimeZoneName(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			(*C.ub1)(unsafe.Pointer(&tzname[0])),
			(*C.ub4)(unsafe.Pointer(&tznamelen))), ts.err)

	if err != nil {
		panic(err)
	}

	return string(tzname[:tznamelen])
}

// GetTimeZoneOffset returns the hour/minute offset from this TimeStamp
func (ts *TimeStamp) GetTimeZoneOffset() (hourOffset, minuteOffset int8) {

	err := checkError(
		C.OCIDateTimeGetTimeZoneOffset(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			(*C.sb1)(unsafe.Pointer(&hourOffset)),
			(*C.sb1)(unsafe.Pointer(&minuteOffset))), ts.err)

	if err != nil {
		panic(err)
	}

	return
}

// ToGoTime converts an Oracle TimeStamp to a go time value
func (ts *TimeStamp) ToGoTime() time.Time {

	year, month, day := ts.GetDate()
	hour, min, sec, fsec := ts.GetTime()
	timezone := ts.GetTimeZoneName()
	hroffs, mnoffs := ts.GetTimeZoneOffset()

	var loc *time.Location
	var locerr error

	loc, locerr = time.LoadLocation(timezone)
	if locerr != nil {
		loc = time.FixedZone("", int((hroffs*60*60))+int((mnoffs*60)))
	}

	return time.Date(int(year), (time.Month)(month), int(day), int(hour), int(min), int(sec), int(fsec), loc)

}

// DateInvalidFlags type to represent invalid date fields/flags.
type DateInvalidFlags uint32

// The possible flags that could be an error in a timestamp.
const (
	DateInvalidDay         DateInvalidFlags = C.OCI_DT_INVALID_DAY
	DateDayBelowValid      DateInvalidFlags = C.OCI_DT_DAY_BELOW_VALID
	DateInvalidMonth       DateInvalidFlags = C.OCI_DT_INVALID_MONTH
	DateMonthBelowValid    DateInvalidFlags = C.OCI_DT_MONTH_BELOW_VALID
	DateInvalidYear        DateInvalidFlags = C.OCI_DT_INVALID_YEAR
	DateYearBelowBalid     DateInvalidFlags = C.OCI_DT_YEAR_BELOW_VALID
	DateInvalidHour        DateInvalidFlags = C.OCI_DT_INVALID_HOUR
	DateHourBelowValid     DateInvalidFlags = C.OCI_DT_HOUR_BELOW_VALID
	DateInvalidMinute      DateInvalidFlags = C.OCI_DT_INVALID_MINUTE
	DateMinuteBelowValid   DateInvalidFlags = C.OCI_DT_MINUTE_BELOW_VALID
	DateInvalidSecond      DateInvalidFlags = C.OCI_DT_INVALID_SECOND
	DateSecondBelowValid   DateInvalidFlags = C.OCI_DT_SECOND_BELOW_VALID
	DateDayMissingFrom1582 DateInvalidFlags = C.OCI_DT_DAY_MISSING_FROM_1582
	DateYearZero           DateInvalidFlags = C.OCI_DT_YEAR_ZERO
	DateInvalidTimezone    DateInvalidFlags = C.OCI_DT_INVALID_TIMEZONE
	DateInvalidFormat      DateInvalidFlags = C.OCI_DT_INVALID_FORMAT
)

// Check validates an Oracle TimeStamp for errors.
func (ts *TimeStamp) Check() (DateInvalidFlags, error) {

	var rslt DateInvalidFlags

	err := checkError(
		C.OCIDateTimeCheck(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			(*C.ub4)(unsafe.Pointer(&rslt))), ts.err)

	return rslt, processError(err)

}

// Compare this TimeStamp with another TimeStamp.
func (ts *TimeStamp) Compare(d2 *TimeStamp) (int, error) {

	var rslt int16

	err := checkError(
		C.OCIDateTimeCompare(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			d2.datetime,
			(*C.sword)(unsafe.Pointer(&rslt))), ts.err)

	return int(rslt), processError(err)
}

// IntervalAdd adds some time to an existing TimeStamp.
func (ts *TimeStamp) IntervalAdd(intvl *Interval) (*TimeStamp, error) {

	rslt := makeTimestampInstance(ts.ses, ts.tstype)

	err := checkError(
		C.OCIDateTimeIntervalAdd(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			ts.datetime,
			intvl.interval,
			rslt.datetime), rslt.err)

	return rslt, processError(err)

}

// IntervalSub subtracts some time from an existing TimeStamp.
func (ts *TimeStamp) IntervalSub(intvl *Interval) (*TimeStamp, error) {

	rslt := makeTimestampInstance(ts.ses, ts.tstype)

	err := checkError(
		C.OCIDateTimeIntervalSub(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			ts.datetime,
			intvl.interval,
			rslt.datetime), rslt.err)

	return rslt, processError(err)

}

// Subtract provided TimeStamp from this TimeStamp. Returns an Interval.
func (ts *TimeStamp) Subtract(d2 *TimeStamp) (*Interval, error) {

	rslt := makeIntervalInstance(ts.ses, TypeIntervalDS)

	err := checkError(
		C.OCIDateTimeSubtract(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			d2.datetime,
			ts.datetime,
			rslt.interval), rslt.err)

	return rslt, processError(err)

}

// ToText returns a TimeStamp as a string; you can provide a format and NLS Lang if you don't want to use the default session settings.
func (ts *TimeStamp) ToText(params ...string) (string, error) {

	buffer := make([]byte, 128)
	var buflen C.ub4 = 128

	var format []byte
	var formatp unsafe.Pointer
	var fmtLength C.ub1

	var langName []byte
	var langNameP unsafe.Pointer
	var langLength C.size_t

	if len(params) > 0 {
		format = []byte(params[0])
		fmtLength = (C.ub1)(len(params[0]))
		formatp = (unsafe.Pointer)(&format[0])
	}

	if len(params) > 1 {
		langName = []byte(params[1])
		langLength = (C.size_t)(len(params[1]))
		langNameP = (unsafe.Pointer)(&langName[0])
	}

	err := checkError(
		C.OCIDateTimeToText(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			(*C.OraText)(formatp),
			fmtLength,
			9,
			(*C.OraText)(langNameP),
			langLength,
			(*C.ub4)(unsafe.Pointer(&buflen)),
			(*C.OraText)(unsafe.Pointer(&buffer[0]))), ts.err)

	return string(buffer[:buflen]), processError(err)

}

func (ts *TimeStamp) String() string {
	rslt, err := ts.ToText()
	if err != nil {
		return err.Error()
	}
	return rslt
}

// IntervalInvalidFlags type to represent invalid interval fields/flags.
type IntervalInvalidFlags uint32

// The possible flags that could be an error in an interval.
const (
	IntervalInvalidDay        IntervalInvalidFlags = C.OCI_INTER_INVALID_DAY
	IntervalDayBelowValid     IntervalInvalidFlags = C.OCI_INTER_DAY_BELOW_VALID
	IntervalInvalidMonth      IntervalInvalidFlags = C.OCI_INTER_INVALID_MONTH
	IntervalMonthBelowValid   IntervalInvalidFlags = C.OCI_INTER_MONTH_BELOW_VALID
	IntervalInvalidYear       IntervalInvalidFlags = C.OCI_INTER_INVALID_YEAR
	IntervalYearBelowValid    IntervalInvalidFlags = C.OCI_INTER_YEAR_BELOW_VALID
	IntervalInvalidHour       IntervalInvalidFlags = C.OCI_INTER_INVALID_HOUR
	IntervalHourBelowValid    IntervalInvalidFlags = C.OCI_INTER_HOUR_BELOW_VALID
	IntervalInvalidMinute     IntervalInvalidFlags = C.OCI_INTER_INVALID_MINUTE
	IntervalMinuteBelowValid  IntervalInvalidFlags = C.OCI_INTER_MINUTE_BELOW_VALID
	IntervalInvalidSecond     IntervalInvalidFlags = C.OCI_INTER_INVALID_SECOND
	IntervalSecondBelowValid  IntervalInvalidFlags = C.OCI_INTER_SECOND_BELOW_VALID
	IntervalInvalidFracsec    IntervalInvalidFlags = C.OCI_INTER_INVALID_FRACSEC
	IntervalFracsecBelowValid IntervalInvalidFlags = C.OCI_INTER_FRACSEC_BELOW_VALID
)

// Check the integrity of an interval.
func (intvl *Interval) Check() (IntervalInvalidFlags, error) {

	var rslt IntervalInvalidFlags

	err := checkError(
		C.OCIIntervalCheck(
			unsafe.Pointer(intvl.ses.ses),
			intvl.err,
			intvl.interval,
			(*C.ub4)(unsafe.Pointer(&rslt))), intvl.err)

	return rslt, processError(err)

}

// Compare this interval with a provided interval.
func (intvl *Interval) Compare(i2 *Interval) (int, error) {

	var rslt int16

	err := checkError(
		C.OCIIntervalCompare(
			unsafe.Pointer(intvl.ses.ses),
			intvl.err,
			intvl.interval,
			i2.interval,
			(*C.sword)(unsafe.Pointer(&rslt))), intvl.err)

	return int(rslt), processError(err)
}

// IntervalFromNumber converts an Oracle Number to an interval.
func (session *Session) IntervalFromNumber(intype IntervalType, num *Number) (*Interval, error) {

	rslt := makeIntervalInstance(session, intype)

	err := checkError(
		C.OCIIntervalFromNumber(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			rslt.interval,
			&num.number), rslt.err)

	return rslt, processError(err)
}

// IntervalFromString converts a string to an interval.
func (session *Session) IntervalFromString(intype IntervalType, intvl string) (*Interval, error) {

	rslt := makeIntervalInstance(session, intype)

	inpstring := []byte(intvl)

	err := checkError(
		C.OCIIntervalFromText(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			(*C.OraText)(unsafe.Pointer(&inpstring[0])),
			(C.size_t)(len(intvl)),
			rslt.interval), rslt.err)

	return rslt, processError(err)
}

// GetDaySecond extracts for a DAY TO SECOND interval
func (intvl *Interval) GetDaySecond() (day, hour, minute, second, fracsecond int32, e error) {

	err := checkError(
		C.OCIIntervalGetDaySecond(
			unsafe.Pointer(intvl.ses.ses),
			intvl.err,
			(*C.sb4)(unsafe.Pointer(&day)),
			(*C.sb4)(unsafe.Pointer(&hour)),
			(*C.sb4)(unsafe.Pointer(&minute)),
			(*C.sb4)(unsafe.Pointer(&second)),
			(*C.sb4)(unsafe.Pointer(&fracsecond)),
			intvl.interval), intvl.err)

	e = processError(err)

	return

}

// GetYearMonth extracts for a YEAR TO MONTH interval
func (intvl *Interval) GetYearMonth() (year, month int32, e error) {

	err := checkError(
		C.OCIIntervalGetYearMonth(
			unsafe.Pointer(intvl.ses.ses),
			intvl.err,
			(*C.sb4)(unsafe.Pointer(&year)),
			(*C.sb4)(unsafe.Pointer(&month)),
			intvl.interval), intvl.err)

	e = processError(err)

	return

}

// SetDaySecond makes a DAY TO SECOND interval
func (session *Session) SetDaySecond(day, hour, minute, second, fracsecond int32) (*Interval, error) {

	rslt := makeIntervalInstance(session, TypeIntervalDS)

	err := checkError(
		C.OCIIntervalSetDaySecond(
			unsafe.Pointer(session.ses),
			rslt.err,
			(C.sb4)(day),
			(C.sb4)(hour),
			(C.sb4)(minute),
			(C.sb4)(second),
			(C.sb4)(fracsecond),
			rslt.interval), rslt.err)

	return rslt, processError(err)

}

// SetYearMonth makes a YEAR TO MONTH interval
func (session *Session) SetYearMonth(year, month int32) (*Interval, error) {

	rslt := makeIntervalInstance(session, TypeIntervalYM)

	err := checkError(
		C.OCIIntervalSetYearMonth(
			unsafe.Pointer(session.ses),
			rslt.err,
			(C.sb4)(year),
			(C.sb4)(month),
			rslt.interval), rslt.err)

	return rslt, processError(err)

}

// IntervalAdd adds provided Interval to this Interval, and returns a new Interval.
func (intvl *Interval) IntervalAdd(i2 *Interval) (*Interval, error) {

	rslt := makeIntervalInstance(intvl.ses, intvl.intype)

	err := checkError(
		C.OCIIntervalAdd(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			intvl.interval,
			i2.interval,
			rslt.interval), rslt.err)

	return rslt, processError(err)

}

// IntervalSubtract subtracts provided Interval from this Interval, and returns a new Interval.
func (intvl *Interval) IntervalSubtract(i2 *Interval) (*Interval, error) {

	rslt := makeIntervalInstance(intvl.ses, intvl.intype)

	err := checkError(
		C.OCIIntervalSubtract(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			intvl.interval,
			i2.interval,
			rslt.interval), rslt.err)

	return rslt, processError(err)

}

// ToNumber returns an Oracle Number from an Interval
func (intvl *Interval) ToNumber() (*Number, error) {

	rslt, e := NumberFromInt(0)

	if e != nil {
		return rslt, e
	}

	err := checkError(
		C.OCIIntervalToNumber(
			unsafe.Pointer(intvl.ses.ses),
			rslt.err,
			intvl.interval,
			&rslt.number), rslt.err)

	return rslt, processError(err)

}

// ToText returns the session formatted interval as a string.
func (intvl *Interval) ToText(params ...uint8) (string, error) {

	buffer := make([]byte, 128)
	var buflen C.size_t = 128
	var resultlen C.size_t

	var lfprec C.ub1 = 9
	var fsprec C.ub1 = 9

	if len(params) > 0 {
		lfprec = (C.ub1)(params[0])
	}

	if len(params) > 1 {
		fsprec = (C.ub1)(params[1])
	}

	err := checkError(
		C.OCIIntervalToText(
			unsafe.Pointer(intvl.ses.ses),
			intvl.err,
			intvl.interval,
			lfprec, fsprec,
			(*C.OraText)(unsafe.Pointer(&buffer[0])),
			buflen, (*C.size_t)(unsafe.Pointer(&resultlen))), intvl.err)

	return string(buffer[:resultlen]), processError(err)

}

func (intvl *Interval) getFloat() (float64, error) {

	num, e := intvl.ToNumber()

	if e != nil {
		return 0.0, e
	}

	return num.ToFloat()

}

// ToGoDuration converts an Oracle Interval to a Go duration type.
func (intvl *Interval) ToGoDuration() (time.Duration, error) {
	if intvl.intype == TypeIntervalDS {
		// unit returned is day
		// multiply by number of seconds in a day (does not count dst, leap seconds, etc)
		// multiply by time.Seconds to return duration
		d, e := intvl.getFloat()
		return time.Duration(d * float64(time.Hour) * 24 * float64(time.Second)), e
	}
	if intvl.intype == TypeIntervalYM {
		// unit returned is year
		// multiply by number of seconds in a year (approximate)
		// multiply time.Seconds to return duration
		d, e := intvl.getFloat()
		return time.Duration(d * 31557600 * float64(time.Second)), e
	}
	return time.Duration(0), errors.New("Unknown Interval Type")
}
