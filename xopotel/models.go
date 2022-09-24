package xopotel

import (
	"context"
	"sync"

	"github.com/muir/xop-go/trace"
	"github.com/muir/xop-go/xopbase"
	"github.com/muir/xop-go/xopnum"

	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type logger struct {
	tracer     oteltrace.Tracer
	id         string
	doLogging  bool
	ignoreDone oteltrace.Span
}

type span struct {
	span               oteltrace.Span
	logger             *logger
	ctx                context.Context
	lock               sync.Mutex
	priorBoolSlices    map[string][]bool
	priorFloat64Slices map[string][]float64
	priorStringSlices  map[string][]string
	priorInt64Slices   map[string][]int64
	hasPrior           map[string]struct{}
	metadataSeen       map[string]interface{}
	spanPrefill        []attribute.KeyValue // holds spanID & traceID
}

type prefilling struct {
	builder
}

type prefilled struct {
	builder
}

type line struct {
	builder
	prealloc [15]attribute.KeyValue
	level    xopnum.Level
}

type builder struct {
	attributes []attribute.KeyValue
	span       *span
	prefillMsg string
	linkKey    string
	linkValue  trace.Trace
}

var _ xopbase.Logger = &logger{}
var _ xopbase.Request = &span{}
var _ xopbase.Span = &span{}
var _ xopbase.Line = &line{}
var _ xopbase.Prefilling = &prefilling{}
var _ xopbase.Prefilled = &prefilled{}

var logMessageKey = attribute.Key("log.message")
var logSpanSequence = attribute.Key("log.xopSpanSequence")
var spanIsLinkAttributeKey = attribute.Key("span.is-link-attribute")
var spanIsLinkEventKey = attribute.Key("span.is-link-event")

var emptyTraceState oteltrace.TraceState
