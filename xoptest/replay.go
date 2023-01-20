package xoptest

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/xoplog/xop-go/xopbase"
	"github.com/xoplog/xop-go/xoptrace"
)

func (log *TestLogger) Replay(ctx context.Context, input any, logger xopbase.Logger) error {
	return log.LosslessReplay(ctx, input, logger)
}

func (_ *TestLogger) LosslessReplay(ctx context.Context, input any, logger xopbase.Logger) error {
	log, ok := input.(*TestLogger)
	if !ok {
		return errors.Errorf("xoptest Replay only supports *TestLogger")
	}
	requests := make(map[xoptrace.HexBytes8]xopbase.Request)
	spans := make(map[xoptrace.HexBytes8]xopbase.Span)
	for _, event := range log.Events {
		switch event.Type {
		case CustomEvent:
			// ignore
		case RequestStart:
			request := logger.Request(ctx, event.Span.StartTime, event.Span.Bundle, event.Span.Name)
			id := event.Span.Bundle.Trace.GetSpanID()
			requests[id] = request
			spans[id] = request
		case RequestDone:
			if req, ok := requests[event.Span.Bundle.Trace.GetSpanID()]; ok {
				req.Done(time.Unix(0, event.Span.EndTime), event.Done)
			} else {
				return errors.Errorf("RequestDone event without corresponding RequestStart for %s", event.Span.Bundle.Trace)
			}
		case SpanDone:
			if span, ok := spans[event.Span.Bundle.Trace.GetSpanID()]; ok {
				req.Done(time.Unix(0, event.Span.EndTime), event.Done)
			} else {
				return errors.Errorf("SpanDone event without corresponding SpanStart for %s", event.Span.Bundle.Trace)
			}
		case FlushEvent:
			id := event.Span.Bundle.Trace.GetSpanID()
			if span, ok := spans[id]; ok {
				span.Flush()
			} else {
				return errors.Errorf("Flush for unknown span %s", event.Span.Bundle.Trace)
			}
		case SpanStart:
			if event.Span.Parent == nil {
				return errors.Errorf("Span w/o parent, %s", event.Span.Bundle.Trace)
			}
			if parent, ok := spans[event.Span.Parent.Bundle.Trace.GetSpanID()]; ok {
				span := parent.Span(replaying.ctx, span.StartTime, span.Bundle, span.Name, span.SequenceCode)
				spans[event.Span.Bundle.Trace.GetSpanID()] = span
			}
		case LineEvent:
			span, ok := spans[event.Line.Span.Bundle.Trace.GetSpanID()]
			if !ok {
				return errors.Errorf("missing span %s for line", event.Line.Span.Bundle.Trace)
			}
			line := span.NoPrefill().Line(event.Line.Level, event.Line.Timestamp, nil /* XXX TODO */)
			for k, v := range event.Line.Data {
				dataType := event.Line.DataType[k]
				switch dataType {
				//MACRO BaseDataWithoutType
				case xopbase.ZZZDataType:
					line.ZZZ(k, v)
				//END
				//MACRO BaseDataWithType
				case xopbase.ZZZDataType:
					line.ZZZ(k, v, dataType)
				//END
				case xopbase.EnumDataType:
					// XXX TODO
				default:
					return errors.Errorf("unexpected data type %s in line", dataType)
				}
			}
		case MetadataSet:
			// XXX
		default:
		}
	}
	return nil
}