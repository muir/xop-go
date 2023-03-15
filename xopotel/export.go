// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file

package xopotel

import (
	"context"

	"github.com/muir/list"
	"github.com/xoplog/xop-go/xopbase"

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
	// XXX
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
	subLinks := make([][]sdktrace.Link, len(spans))
	id2Index := make(map[oteltrace.SpanID]int)
	for i, span := range spans {
		spanContext := span.SpanContext()
		if spanContext.HasSpanID() {
			id2Index[spanContext.SpanID()] = i
		}
	}
	for i, span := range spans {
		parent := span.Parent()
		if !parent.HasSpanID() {
			continue
		}
		parentIndex, ok := id2Index[parent.SpanID()]
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
