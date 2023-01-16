// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file

/*
Package xopbase defines the base-level loggers for xop.

In xop, the logger is divided into a top-level logger and a lower-level
logger.  The top-level logger is what users use to add logs to their
programs.  The lower-level logger is used to send those logs somewhere
useful.  Xopjson, xoptest, and xopotel are all examples.

In OpenTelemetry, these are called "exporters".
In logr, they are called "LogSinks".
*/
package xopbase

import (
	"context"
	"time"

	"github.com/xoplog/xop-go/xopat"
	"github.com/xoplog/xop-go/xopnum"
	"github.com/xoplog/xop-go/xoptrace"
)

// Logger is the bottom half of a logger -- the part that actually
// outputs data somewhere.  There can be many Logger implementations.
type Logger interface {
	// Request beings a new span that represents the start of an
	// operation: either a call to a server, a cron-job, or an event
	// being processed.  The provided Context is a pass-through from
	// the Seed and if the seed does not provide a context, the context
	// can be nil.
	Request(ctx context.Context, ts time.Time, span xoptrace.Bundle, description string) Request

	// ID returns a unique id for this instance of a logger.  This
	// is used to prevent duplicate Requets from being created when
	// additional loggers to added to an ongoing Request.
	ID() string

	// ReferencesKept should return true if Any() objects are not immediately
	// serialized (the object is kept around and serilized later).  If copies
	// are kept, then xop.Log will make copies.
	ReferencesKept() bool

	// Buffered should return true if the logger buffers output and sends
	// it when Flush() is called. Even if Buffered() returns false,
	// Flush() may still be invoked but it doesn't have to do anything.
	Buffered() bool

	// Replay is how a Logger will feed it's output to another Logger.  The
	// input param should be the collected output from the first Logger.  If
	// it isn't in the right format, or type, it should throw an error.  This
	// capability can be used to transform from one format to another, it can
	// also be used to during testing to make sure that a Logger can round-trip
	// without data loss.  Loggers are not required to round-trip without
	// data loss.
	Replay(ctx context.Context, input any, logger Logger) error
}

type RoundTripLogger interface {
	Logger

	// LosslessReplay is just like Replay but guarantees that no data
	// is lost.
	LosslessReplay(ctx context.Context, input any, logger Logger) error
}

type Request interface {
	Span

	// Flush calls are single-threaded.  Flush can be triggered explicitly by
	// users and it can be triggered because all parts of a request have had
	// Done() called on them.
	Flush()

	// SetErrorReported will always be called before any other method on the
	// Request.
	//
	// If a base logger encounters an error, it may use the provided function to
	// report it.  The base logger cannot assume that execution will stop.
	// The base logger may not panic.
	SetErrorReporter(func(error))

	// Final is called when there is no possibility of any further calls to
	// this Request of any sort.  There is no guarantee that Final will be
	// called in a timely fashion or even at all before program exit.
	Final()
}

type Span interface {
	// Span creates a new Span that should inherit prefil but not data
	Span(ctx context.Context, ts time.Time, span xoptrace.Bundle, descriptionOrName string, spanSequenceCode string) Span

	// MetadataAny adds a key/value pair to describe the span.
	MetadataAny(*xopat.AnyAttribute, interface{})
	// MetadataBool adds a key/value pair to describe the span.
	MetadataBool(*xopat.BoolAttribute, bool)
	// MetadataEnum adds a key/value pair to describe the span.
	MetadataEnum(*xopat.EnumAttribute, xopat.Enum)
	// MetadataFloat64 adds a key/value pair to describe the span.
	MetadataFloat64(*xopat.Float64Attribute, float64)
	// MetadataInt64 adds a key/value pair to describe the span.
	MetadataInt64(*xopat.Int64Attribute, int64)
	// MetadataString adds a key/value pair to describe the span.
	MetadataString(*xopat.StringAttribute, string)
	// MetadataTime adds a key/value pair to describe the span.
	MetadataTime(*xopat.TimeAttribute, time.Time)

	// Boring true indicates that a span (or request) is boring.  The
	// suggested meaning for this is that a boring request that is buffered
	// can ignore Flush() and never get sent to output.  A boring span
	// can be un-indexed. Boring requests that do get sent to output can
	// be marked as boring so that they're dropped at the indexing stage.
	Boring(bool)

	// ID must return the same string as the Logger it came from
	ID() string

	// TODO: Gauge()
	// TODO: Event()

	NoPrefill() Prefilled

	StartPrefill() Prefilling

	// Done is called when (1) log.Done is called on the log corresponding
	// to this span; (2) log.Done is called on a parent log of the log
	// corresponding to this span, and the log is not Detach()ed; or
	// (3) preceding Flush() if there has been logging activity since the
	// last call to Flush(), Done(), or the start of the span.
	//
	// final is true when the log is done, it is false when Done is called
	// prior to a Flush().  Just because Done was called with final does not
	// mean that Done won't be called again.  Any further calls would only
	// happen due a bug in the application: logging things after calling
	// log.Done.
	//
	// If the application never calls log.Done(), then final will never
	// be true.
	Done(endTime time.Time, final bool)
}

type Prefilling interface {
	AttributeParts

	PrefillComplete(msg string) Prefilled
}

type Prefilled interface {
	// Line starts another line of log output.  Span implementations
	// can expect multiple calls simultaneously and even during a call
	// to SpanInfo() or Flush().  The []uintptr slice are stack frames.
	Line(xopnum.Level, time.Time, []uintptr) Line
}

type Line interface {
	AttributeParts

	LineDone
}

// LineDone are methods that complete the line.  No additional methods may
// be invoked on the line after one of these is called.
type LineDone interface {
	Msg(string)
	Template(string)
	// Static(string) XXX remove implemenations
	// Object may change in the future to also take an un-redaction string,
	Model(string, ModelArg)
	Link(string, xoptrace.Trace)
}

// ModelArg may be expanded in the future to supply: an Encoder; redaction
// information.
type ModelArg struct {
	// If specified, overrides what would be provided by reflect.TypeOf(obj).Name()
	TypeName string
	Model    interface{}
}

type AttributeParts interface {
	Enum(*xopat.EnumAttribute, xopat.Enum)
	Any(key string, typeName string, unredaction string, obj interface{})

	// MACRO BaseDataWithType Skip:Any
	ZZZ(string, zzz, DataType)

	Any(string, interface{})
	Bool(string, bool)
	Duration(string, time.Duration)
	Time(string, time.Time)

	// TODO: split the above off as "BasicAttributeParts"
	// TODO: Table(string, table)
	// TODO: Encoded(name string, elementName string, encoder Encoder, data interface{})
	// TODO: PreEncodedBytes(name string, elementName string, mimeType string, data []byte)
	// TODO: PreEncodedText(name string, elementName string, mimeType string, data string)
	// TODO: ExternalReference(name string, itemID string, storageID string)
	// TODO: RedactedString(name string, value string, unredaction string)
}

// TODO
type Encoder interface {
	MimeType() string
	ProducesText() bool
	Encode(elementName string, data interface{}) ([]byte, error)
}

//go:generate enumer -type=DataType -linecomment -json -sql

type DataType int

const (
	EnumDataType          DataType = iota
	EnumArrayDataType     DataType = iota
	AnyDataType           DataType = iota
	BoolDataType          DataType = iota
	DurationDataType      DataType = iota
	ErrorDataType         DataType = iota
	Float32DataType       DataType = iota
	Float64DataType       DataType = iota
	IntDataType           DataType = iota
	Int16DataType         DataType = iota
	Int32DataType         DataType = iota
	Int64DataType         DataType = iota
	Int8DataType          DataType = iota
	LinkDataType          DataType = iota
	StringDataType        DataType = iota
	StringerDataType      DataType = iota
	TimeDataType          DataType = iota
	UintDataType          DataType = iota
	Uint16DataType        DataType = iota
	Uint32DataType        DataType = iota
	Uint64DataType        DataType = iota
	Uint8DataType         DataType = iota
	AnyArrayDataType      DataType = iota
	BoolArrayDataType     DataType = iota
	DurationArrayDataType DataType = iota
	ErrorArrayDataType    DataType = iota
	Float32ArrayDataType  DataType = iota
	Float64ArrayDataType  DataType = iota
	IntArrayDataType      DataType = iota
	Int16ArrayDataType    DataType = iota
	Int32ArrayDataType    DataType = iota
	Int64ArrayDataType    DataType = iota
	Int8ArrayDataType     DataType = iota
	LinkArrayDataType     DataType = iota
	StringArrayDataType   DataType = iota
	StringerArrayDataType DataType = iota
	TimeArrayDataType     DataType = iota
	UintArrayDataType     DataType = iota
	Uint16ArrayDataType   DataType = iota
	Uint32ArrayDataType   DataType = iota
	Uint64ArrayDataType   DataType = iota
	Uint8ArrayDataType    DataType = iota
)
