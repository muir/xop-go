// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file

package xopotel

import (
	"context"

	"github.com/muir/list"
	"github.com/xoplog/xop-go/internal/util/version"
	"github.com/xoplog/xop-go/xopbase"
	"github.com/xoplog/xop-go/xoptrace"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var (
	_ sdktrace.SpanExporter = &spanExporter{}
	_ sdktrace.SpanExporter = &unhack{}
)

type spanExporter struct {
	base xopbase.Logger
}

func NewExporter(base xopbase.Logger) sdktrace.SpanExporter {
	return &spanExporter{base: base}
}

func (e *spanExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	id2Index := makeIndex(spans)
	subSpans, todo := makeSubspans(id2Index, spans)
	_ = subSpans // XXX
	baseSpans := make([]xopbase.Span, len(spans))
	for _, i := range todo {
		span := spans[i]
		attributeMap := mapAttributes(span)
		var bundle xoptrace.Bundle
		spanContext := span.SpanContext()
		if spanContext.HasTraceID() {
			bundle.Trace.TraceID().SetArray(spanContext.TraceID())
		}
		if spanContext.HasSpanID() {
			bundle.Trace.SpanID().SetArray(spanContext.SpanID())
		}
		if spanContext.IsSampled() {
			bundle.Trace.Flags().SetArray([1]byte{1})
		}
		if spanContext.TraceState().Len() != 0 {
			bundle.State.SetString(spanContext.TraceState().String())
		}
		parentIndex, ok := lookupParent(id2Index, span)
		if ok {
			parentContext := spans[parentIndex].SpanContext()
			xopParent := baseSpans[parentIndex]
			if parentContext.HasTraceID() {
				bundle.Parent.TraceID().SetArray(parentContext.TraceID())
				if bundle.Trace.TraceID().IsZero() {
					bundle.Trace.TraceID().Set(bundle.Parent.GetTraceID())
				}
			}
			if parentContext.HasSpanID() {
				bundle.Parent.SpanID().SetArray(parentContext.SpanID())
			}
			if parentContext.IsSampled() {
				bundle.Parent.Flags().SetArray([1]byte{1})
			}
			bundle.Parent.Version().SetArray([1]byte{1})
		}
		bundle.Trace.Version().SetArray([1]byte{1})
		spanKind := span.SpanKind()
		if spanKind == oteltrace.SpanKindUnspecified {
			spanKind = oteltrace.SpanKind(defaulted(attributeMap.GetInt(otelSpanKind), int(oteltrace.SpanKindUnspecified)))
		}
		switch spanKind {
		case oteltrace.SpanKindUnspecified, oteltrace.SpanKindInternal:
			if ok {
				baseSpan := xopParent.Span(ctx, span.StartTime(), bundle, span.Name(), defaulted(attributeMap.GetString(logSpanSequence), ""))
				baseSpans[i] = baseSpan
			} else {
				// This is a difficult sitatuion. We have an internal/unspecified span
				// that does not have a parent present. There is no right answer for what
				// to do. In the Xop world, such a span isn't allowed to exist. We'll treat
				// this span as a request, but mark it as promoted.
				request := e.base.Request(ctx, span.StartTime(), bundle, span.Name(), buildSourceInfo(span, attributeMap))
				request.MetadataBool(xopPromotedMetadata, true)
				baseSpans[i] = request
			}
		default:
			request := e.base.Request(ctx, span.StartTime(), bundle, span.Name(), buildSourceInfo(span, attributeMap))
			baseSpans[i] = request
		}
		baseSpans[i].MetadataEnum(otelconst.SpanKind, otelconst.SpanKindEnum(spanKind))
	}
	return nil
}

func (e *spanExporter) Shutdown(ctx context.Context) error {
	// XXX
	return nil
}

type unhack struct {
	next sdktrace.SpanExporter
}

// NewUnhacker wraps a SpanExporter and if the input is from BaseLogger or SpanLog,
// then it "fixes" the data hack in the output that puts inter-span links in sub-spans
// rather than in the span that defined them.
func NewUnhacker(exporter sdktrace.SpanExporter) sdktrace.SpanExporter {
	return &unhack{next: exporter}
}

func (u *unhack) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	// TODO: fix up SpanKind if spanKind is one of the attributes
	id2Index := makeIndex(spans)
	subLinks := make([][]sdktrace.Link, len(spans))
	for i, span := range spans {
		parentIndex, ok := lookupParent(id2Index, span)
		if !ok {
			continue
		}
		var addToParent bool
		for _, attribute := range span.Attributes() {
			switch attribute.Key {
			case spanIsLinkAttributeKey, spanIsLinkEventKey:
				spans[i] = nil
				addToParent = true
			}
		}
		if !addToParent {
			continue
		}
		subLinks[parentIndex] = append(subLinks[parentIndex], span.Links()...)
	}
	n := make([]sdktrace.ReadOnlySpan, 0, len(spans))
	for i, span := range spans {
		span := span
		switch {
		case len(subLinks[i]) > 0:
			n = append(n, wrappedReadOnlySpan{
				ReadOnlySpan: span,
				links:        append(list.Copy(span.Links()), subLinks[i]...),
			})
		case span == nil:
			// skip
		default:
			n = append(n, span)
		}
	}
	return u.next.ExportSpans(ctx, n)
}

func (u *unhack) Shutdown(ctx context.Context) error {
	return u.next.Shutdown(ctx)
}

type wrappedReadOnlySpan struct {
	sdktrace.ReadOnlySpan
	links []sdktrace.Link
}

var _ sdktrace.ReadOnlySpan = wrappedReadOnlySpan{}

func (w wrappedReadOnlySpan) Links() []sdktrace.Link {
	return w.links
}

func makeIndex(spans []sdktrace.ReadOnlySpan) map[oteltrace.SpanID]int {
	id2Index := make(map[oteltrace.SpanID]int)
	for i, span := range spans {
		spanContext := span.SpanContext()
		if spanContext.HasSpanID() {
			id2Index[spanContext.SpanID()] = i
		}
	}
	return id2Index
}

func lookupParent(id2Index map[oteltrace.SpanID]int, span sdktrace.ReadOnlySpan) (int, bool) {
	parent := span.Parent()
	if !parent.HasSpanID() {
		return 0, false
	}
	parentIndex, ok := id2Index[parent.SpanID()]
	if !ok {
		return 0, false
	}
	return parentIndex, true
}

func makeSubspans(id2Index map[oteltrace.SpanID]int, spans []sdktrace.ReadOnlySpan) ([][]oteltrace.SpanID, []int) {
	ss := make([][]oteltrace.SpanID, len(spans))
	noParent := make([]int, 0, len(spans))
	for i, span := range spans {
		parentIndex, ok := lookupParent(id2Index, span)
		if !ok {
			noParent = append(noParent, i)
		}
		ss[parentIndex] = append(ss[parentIndex], i)
	}
	return ss, noParent
}

func buildSourceInfo(span sdktrace.ReadOnlySpan, attributeMap AttributeMap) {
	var si xopbase.SourceInfo
	var source string
	// XXX grab namespace from scope instead
	if s := attributeMap.GetString(xopSource); s != "" {
		source = s
	} else if n := span.InstrumentationScope().Name; n != "" {
		if v := span.InstrumentationScope().Version; v != "" {
			source = n + " " + v
		} else {
			source = n
		}
	} else {
		source = "OTEL"
	}
	namespace := defaulted(attributeMap.GetString(xopNamespace), source)
	si.Source, si.SourceVersion = version.SplitVersion(source)
	si.Namespace, si.NamespaceVersion = version.SplitVersion(namespace)
	return si
}
