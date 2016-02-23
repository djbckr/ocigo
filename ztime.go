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

type TimestampType C.ub4

// timestamp types
const (
	TypeTimestamp    TimestampType = (TimestampType)(dtypeTimestamp)
	TypeTimestampTZ  TimestampType = (TimestampType)(dtypeTimestampTZ)
	TypeTimestampLTZ TimestampType = (TimestampType)(dtypeTimestampLTZ)
)

type IntervalType C.ub4

// interval types
const (
	TypeIntervalYM IntervalType = (IntervalType)(dtypeIntervalYM)
	TypeIntervalDS IntervalType = (IntervalType)(dtypeIntervalDS)
)

type TimeStamp struct {
	ses      *Session
	err      *C.OCIError
	datetime *C.OCIDateTime
	tstype   TimestampType
}

type Interval struct {
	ses      *Session
	err      *C.OCIError
	interval *C.OCIInterval
	intype   IntervalType
}

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

/* Get System Time Stamp based on Database Session settings */
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

func (session *Session) TimeStampFromText(tstype TimestampType, params ...string) (*TimeStamp, error) {
	if len(params) == 0 {
		return nil, nil
	}

	date_str := []byte(params[0])
	dstr_length := len(params[0])

	var fmt []byte
	var fmtp unsafe.Pointer
	var fmt_length C.ub1 = 0
	var lang_name []byte
	var lang_namep unsafe.Pointer
	var lang_length C.size_t = 0

	if len(params) > 1 && len(params[1]) > 0 {
		fmt = []byte(params[1])
		fmt_length = (C.ub1)(len(params[1]))
		fmtp = (unsafe.Pointer)(&fmt[0])
	}

	if len(params) > 2 && len(params[2]) > 0 {
		lang_name = []byte(params[2])
		lang_length = (C.size_t)(len(params[2]))
		lang_namep = (unsafe.Pointer)(&lang_name[0])
	}

	rslt := makeTimestampInstance(session, tstype)

	err := checkError(
		C.OCIDateTimeFromText(
			unsafe.Pointer(session.ses),
			rslt.err,
			(*C.OraText)(unsafe.Pointer(&date_str[0])),
			(C.size_t)(dstr_length),
			(*C.OraText)(fmtp),
			fmt_length,
			(*C.OraText)(lang_namep),
			lang_length,
			rslt.datetime), rslt.err)

	return rslt, processError(err)

}

/*****************************************************************************/

func (session *Session) TimeStampFromGoTime(tstype TimestampType, t time.Time) (*TimeStamp, error) {

	var year int16 = int16(t.Year())
	var month uint8 = uint8(t.Month())
	var day uint8 = uint8(t.Day())
	var hour uint8 = uint8(t.Hour())
	var min uint8 = uint8(t.Minute())
	var sec uint8 = uint8(t.Second())
	var fsec uint32 = uint32(t.Nanosecond() / 1000000)
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

func (dt *TimeStamp) GetDate() (year int16, month, day uint8) {

	err := checkError(
		C.OCIDateTimeGetDate(
			unsafe.Pointer(dt.ses.ses),
			dt.err,
			dt.datetime,
			(*C.sb2)(unsafe.Pointer(&year)),
			(*C.ub1)(unsafe.Pointer(&month)),
			(*C.ub1)(unsafe.Pointer(&day))), dt.err)

	if err != nil {
		panic(err)
	}

	return
}

func (dt *TimeStamp) GetTime() (hour, minute, second uint8, fracsecond uint32) {

	err := checkError(
		C.OCIDateTimeGetTime(
			unsafe.Pointer(dt.ses.ses),
			dt.err,
			dt.datetime,
			(*C.ub1)(unsafe.Pointer(&hour)),
			(*C.ub1)(unsafe.Pointer(&minute)),
			(*C.ub1)(unsafe.Pointer(&second)),
			(*C.ub4)(unsafe.Pointer(&fracsecond))), dt.err)

	if err != nil {
		panic(err)
	}

	return
}

func (dt *TimeStamp) GetTimeZoneName() string {

	tzname := make([]byte, 128)
	var tznamelen uint32 = 128

	err := checkError(
		C.OCIDateTimeGetTimeZoneName(
			unsafe.Pointer(dt.ses.ses),
			dt.err,
			dt.datetime,
			(*C.ub1)(unsafe.Pointer(&tzname[0])),
			(*C.ub4)(unsafe.Pointer(&tznamelen))), dt.err)

	if err != nil {
		panic(err)
	}

	return string(tzname[:tznamelen])
}

func (dt *TimeStamp) GetTimeZoneOffset() (hourOffset, minuteOffset int8) {

	err := checkError(
		C.OCIDateTimeGetTimeZoneOffset(
			unsafe.Pointer(dt.ses.ses),
			dt.err,
			dt.datetime,
			(*C.sb1)(unsafe.Pointer(&hourOffset)),
			(*C.sb1)(unsafe.Pointer(&minuteOffset))), dt.err)

	if err != nil {
		panic(err)
	}

	return
}

func (dt *TimeStamp) ToGoTime() time.Time {

	year, month, day := dt.GetDate()
	hour, min, sec, fsec := dt.GetTime()
	timezone := dt.GetTimeZoneName()
	hroffs, mnoffs := dt.GetTimeZoneOffset()

	var loc *time.Location
	var locerr error

	loc, locerr = time.LoadLocation(timezone)
	if locerr != nil {
		loc = time.FixedZone("", int((hroffs*60*60))+int((mnoffs*60)))
	}

	return time.Date(int(year), (time.Month)(month), int(day), int(hour), int(min), int(sec), int(fsec), loc)

}

type DateInvalidFlags uint32

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

func (dt *TimeStamp) Check() (DateInvalidFlags, error) {

	var rslt DateInvalidFlags

	err := checkError(
		C.OCIDateTimeCheck(
			unsafe.Pointer(dt.ses.ses),
			dt.err,
			dt.datetime,
			(*C.ub4)(unsafe.Pointer(&rslt))), dt.err)

	return rslt, processError(err)

}

func (dt *TimeStamp) Compare(d2 *TimeStamp) (int, error) {

	var rslt int16

	err := checkError(
		C.OCIDateTimeCompare(
			unsafe.Pointer(dt.ses.ses),
			dt.err,
			dt.datetime,
			d2.datetime,
			(*C.sword)(unsafe.Pointer(&rslt))), dt.err)

	return int(rslt), processError(err)
}

func (dt *TimeStamp) IntervalAdd(intvl *Interval) (*TimeStamp, error) {

	rslt := makeTimestampInstance(dt.ses, dt.tstype)

	err := checkError(
		C.OCIDateTimeIntervalAdd(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			dt.datetime,
			intvl.interval,
			rslt.datetime), rslt.err)

	return rslt, processError(err)

}

func (dt *TimeStamp) IntervalSub(intvl *Interval) (*TimeStamp, error) {

	rslt := makeTimestampInstance(dt.ses, dt.tstype)

	err := checkError(
		C.OCIDateTimeIntervalSub(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			dt.datetime,
			intvl.interval,
			rslt.datetime), rslt.err)

	return rslt, processError(err)

}

func (dt *TimeStamp) Subtract(d2 *TimeStamp) (*Interval, error) {

	rslt := makeIntervalInstance(dt.ses, TypeIntervalDS)

	err := checkError(
		C.OCIDateTimeSubtract(
			unsafe.Pointer(rslt.ses.ses),
			rslt.err,
			d2.datetime,
			dt.datetime,
			rslt.interval), rslt.err)

	return rslt, processError(err)

}

func (ts *TimeStamp) ToText(params ...string) (string, error) {

	buffer := make([]byte, 128)
	var buflen C.ub4 = 128

	var format []byte
	var formatp unsafe.Pointer
	var fmt_length C.ub1 = 0

	var lang_name []byte
	var lang_namep unsafe.Pointer
	var lang_length C.size_t = 0

	if len(params) > 0 {
		format = []byte(params[0])
		fmt_length = (C.ub1)(len(params[0]))
		formatp = (unsafe.Pointer)(&format[0])
	}

	if len(params) > 1 {
		lang_name = []byte(params[1])
		lang_length = (C.size_t)(len(params[1]))
		lang_namep = (unsafe.Pointer)(&lang_name[0])
	}

	err := checkError(
		C.OCIDateTimeToText(
			unsafe.Pointer(ts.ses.ses),
			ts.err,
			ts.datetime,
			(*C.OraText)(formatp),
			fmt_length,
			9,
			(*C.OraText)(lang_namep),
			lang_length,
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

type IntervalInvalidFlags uint32

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
