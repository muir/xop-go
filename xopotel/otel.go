// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file

package xopotel

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/xoplog/xop-go"
	"github.com/xoplog/xop-go/xopat"
	"github.com/xoplog/xop-go/xopbase"
	"github.com/xoplog/xop-go/xopnum"
	"github.com/xoplog/xop-go/xoptrace"
	"github.com/xoplog/xop-go/xoputil"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// SpanLog allows xop to add logs to an existing OTEL span.  log.Done() will be
// ignored for this span.
func SpanLog(ctx context.Context, name string, extraModifiers ...xop.SeedModifier) *xop.Log {
	span := oteltrace.SpanFromContext(ctx)
	var xoptrace xoptrace.Trace
	xoptrace.TraceID().Set(span.SpanContext().TraceID())
	xoptrace.SpanID().Set(span.SpanContext().SpanID())
	xoptrace.Flags().Set([1]byte{byte(span.SpanContext().TraceFlags())})
	xoptrace.Version().Set([1]byte{1})
	log := xop.NewSeed(
		xop.CombineSeedModifiers(extraModifiers...),
		xop.WithContext(ctx),
		xop.WithTrace(xoptrace),
		xop.WithBase(&logger{
			id:         "otel-" + uuid.New().String(),
			doLogging:  true,
			ignoreDone: span,
			tracer:     span.TracerProvider().Tracer(""),
		}),
		// The first time through, we do not want to change the spanID,
		// but on subsequent calls, we do so the outer reactive function
		// just sets the future function.
		xop.WithReactive(func(ctx context.Context, seed xop.Seed, nameOrDescription string, isChildSpan bool) []xop.SeedModifier {
			return []xop.SeedModifier{
				xop.WithTrace(xoptrace),
				xop.WithReactiveReplaced(
					func(ctx context.Context, seed xop.Seed, nameOrDescription string, isChildSpan bool) []xop.SeedModifier {
						var newSpan oteltrace.Span
						if isChildSpan {
							ctx, newSpan = span.TracerProvider().Tracer("").Start(ctx, nameOrDescription, oteltrace.WithSpanKind(oteltrace.SpanKindInternal))
						} else {
							ctx, newSpan = span.TracerProvider().Tracer("").Start(ctx, nameOrDescription)
						}
						return []xop.SeedModifier{
							xop.WithContext(ctx),
							xop.WithSpan(newSpan.SpanContext().SpanID()),
						}
					}),
			}
		}),
	).SubSpan(name)
	go func() {
		<-ctx.Done()
		log.Done()
	}()
	return log
}

// BaseLogger provides SeedModifiers to set up an OTEL Tracer as a xopbase.Logger
// so that xop logs are output through the OTEL Tracer.
func BaseLogger(ctx context.Context, tracer oteltrace.Tracer, doLogging bool) xop.SeedModifier {
	return xop.CombineSeedModifiers(
		xop.WithBase(&logger{
			id:        "otel-" + uuid.New().String(),
			doLogging: doLogging,
			tracer:    tracer,
		}),
		xop.WithContext(ctx),
		xop.WithReactive(func(ctx context.Context, seed xop.Seed, nameOrDescription string, isChildSpan bool) []xop.SeedModifier {
			if isChildSpan {
				ctx, span := tracer.Start(ctx, nameOrDescription, oteltrace.WithSpanKind(oteltrace.SpanKindInternal))
				return []xop.SeedModifier{
					xop.WithContext(ctx),
					xop.WithSpan(span.SpanContext().SpanID()),
				}
			}
			ctx, span := tracer.Start(ctx, nameOrDescription, oteltrace.WithNewRoot())
			bundle := seed.Bundle()
			if bundle.Parent.IsZero() {
				bundle.State.SetString(span.SpanContext().TraceState().String())
				bundle.Trace.Flags().Set([1]byte{byte(span.SpanContext().TraceFlags())})
				bundle.Trace.Version().Set([1]byte{1})
				bundle.Trace.TraceID().Set(span.SpanContext().TraceID())
			}
			bundle.Trace.SpanID().Set(span.SpanContext().SpanID())
			return []xop.SeedModifier{
				xop.WithContext(ctx),
				xop.WithBundle(bundle),
			}
		}),
	)
}

func (logger *logger) ID() string           { return logger.id }
func (logger *logger) ReferencesKept() bool { return true }
func (logger *logger) Buffered() bool       { return false }

func (logger *logger) Request(ctx context.Context, ts time.Time, _ xoptrace.Bundle, description string) xopbase.Request {
	return logger.span(ctx, ts, description, "")
}

func (span *span) Flush()                         {}
func (span *span) Final()                         {}
func (span *span) SetErrorReporter(f func(error)) {}
func (span *span) Boring(_ bool)                  {}
func (span *span) ID() string                     { return span.logger.id }
func (span *span) Done(endTime time.Time, final bool) {
	if !final {
		return
	}
	if span.logger.ignoreDone == span.span {
		// skip Done for spans passed in to SpanLog()
		return
	}
	span.span.End()
}

func (span *span) Span(ctx context.Context, ts time.Time, bundle xoptrace.Bundle, description string, spanSequenceCode string) xopbase.Span {
	return span.logger.span(ctx, ts, description, spanSequenceCode)
}

func (logger *logger) span(ctx context.Context, ts time.Time, description string, spanSequence string) xopbase.Request {
	otelspan := oteltrace.SpanFromContext(ctx)
	if spanSequence != "" {
		otelspan.SetAttributes(logSpanSequence.String(spanSequence))
	}
	return &span{
		logger: logger,
		span:   otelspan,
		ctx:    ctx,
	}
}

func (span *span) NoPrefill() xopbase.Prefilled {
	return &prefilled{
		builder: builder{
			span: span,
		},
	}
}

func (span *span) StartPrefill() xopbase.Prefilling {
	return &prefilling{
		builder: builder{
			span: span,
		},
	}
}

func (prefill *prefilling) PrefillComplete(msg string) xopbase.Prefilled {
	prefill.builder.prefillMsg = msg
	return &prefilled{
		builder: prefill.builder,
	}
}

func (prefilled *prefilled) Line(level xopnum.Level, _ time.Time, pc []uintptr) xopbase.Line {
	if !prefilled.span.logger.doLogging || !prefilled.span.span.IsRecording() {
		return xoputil.SkipLine
	}
	// PERFORMANCE: get line from a pool
	line := &line{}
	line.level = level
	line.span = prefilled.span
	line.attributes = line.prealloc[:0]
	line.attributes = append(line.attributes, prefilled.span.spanPrefill...)
	line.attributes = append(line.attributes, prefilled.attributes...)
	line.prefillMsg = prefilled.prefillMsg
	line.linkKey = prefilled.linkKey
	line.linkValue = prefilled.linkValue
	if len(pc) > 0 {
		var b strings.Builder
		frames := runtime.CallersFrames(pc)
		for {
			frame, more := frames.Next()
			if strings.Contains(frame.File, "runtime/") {
				break
			}
			b.WriteString(frame.File)
			b.WriteByte(':')
			b.WriteString(strconv.Itoa(frame.Line))
			b.WriteByte('\n')
			if !more {
				break
			}
		}
		line.attributes = append(line.attributes, semconv.ExceptionStacktraceKey.String(b.String()))
	}
	return line
}

func (line *line) Link(k string, v xoptrace.Trace) {
	line.attributes = append(line.attributes,
		logMessageKey.String(line.prefillMsg+k),
		typeKey.String("link"),
		attribute.StringSlice("xop.link", []string{"link", v.TraceID().String(), v.SpanID().String()}),
	)
	_, tmpSpan := line.span.logger.tracer.Start(line.span.ctx, k, oteltrace.WithLinks(
		oteltrace.Link{
			SpanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
				TraceID:    v.TraceID().Array(),
				SpanID:     v.SpanID().Array(),
				TraceFlags: oteltrace.TraceFlags(v.Flags().Array()[0]),
				TraceState: emptyTraceState, // TODO: is this right?
				Remote:     true,            // information not available
			}),
		},
	))
	tmpSpan.AddEvent(line.level.String(), oteltrace.WithAttributes(line.attributes...))
	tmpSpan.SetAttributes(typeKey.String("link-event"))
	tmpSpan.End()
}

func (line *line) Model(msg string, modelArg xopbase.ModelArg) {}
func (line *line) Msg(msg string) {
	line.attributes = append(line.attributes, logMessageKey.String(line.prefillMsg+msg), typeKey.String("line"))
	if line.linkKey == "" {
		line.span.span.AddEvent(line.level.String(), oteltrace.WithAttributes(line.attributes...))
		return
		// PERFORMANCE: return line to pool
	}
}

var templateRE = regexp.MustCompile(`\{.+?\}`)

func (line *line) Template(template string) {
	kv := make(map[string]int)
	for i, a := range line.attributes {
		kv[string(a.Key)] = i
	}
	msg := templateRE.ReplaceAllStringFunc(template, func(k string) string {
		k = k[1 : len(k)-1]
		if i, ok := kv[k]; ok {
			a := line.attributes[i]
			switch a.Value.Type() {
			case attribute.BOOL:
				return strconv.FormatBool(a.Value.AsBool())
			case attribute.INT64:
				return strconv.FormatInt(a.Value.AsInt64(), 10)
			case attribute.FLOAT64:
				return strconv.FormatFloat(a.Value.AsFloat64(), 'g', -1, 64)
			case attribute.STRING:
				return a.Value.AsString()
			case attribute.BOOLSLICE:
				return fmt.Sprint(a.Value.AsBoolSlice())
			case attribute.INT64SLICE:
				return fmt.Sprint(a.Value.AsInt64Slice())
			case attribute.FLOAT64SLICE:
				return fmt.Sprint(a.Value.AsFloat64Slice())
			case attribute.STRINGSLICE:
				return fmt.Sprint(a.Value.AsStringSlice())
			default:
				return "{" + k + "}"
			}
		}
		return "''"
	})
	line.Msg(msg)
}

func (builder *builder) Enum(k *xopat.EnumAttribute, v xopat.Enum) {
	builder.attributes = append(builder.attributes, attribute.Stringer(k.Key(), v))
}

func (builder *builder) Any(k string, v xopbase.ModelArg) {
	switch typed := v.Model.(type) {
	case bool:
		builder.attributes = append(builder.attributes, attribute.Bool(k, typed))
	case []bool:
		builder.attributes = append(builder.attributes, attribute.BoolSlice(k, typed))
	case float64:
		builder.attributes = append(builder.attributes, attribute.Float64(k, typed))
	case []float64:
		builder.attributes = append(builder.attributes, attribute.Float64Slice(k, typed))
	case int64:
		builder.attributes = append(builder.attributes, attribute.Int64(k, typed))
	case []int64:
		builder.attributes = append(builder.attributes, attribute.Int64Slice(k, typed))
	case string:
		builder.attributes = append(builder.attributes, attribute.String(k, typed))
	case []string:
		builder.attributes = append(builder.attributes, attribute.StringSlice(k, typed))
	case fmt.Stringer:
		builder.attributes = append(builder.attributes, attribute.Stringer(k, typed))

	default:
		enc, err := json.Marshal(v)
		if err != nil {
			builder.attributes = append(builder.attributes, attribute.String(k+"-error", err.Error()))
		} else {
			builder.attributes = append(builder.attributes, attribute.String(k, string(enc)))
		}
	}
}

func (builder *builder) Time(k string, v time.Time) {
	builder.attributes = append(builder.attributes, attribute.String(k, v.Format(time.RFC3339Nano)))
}

func (builder *builder) Duration(k string, v time.Duration) {
	builder.attributes = append(builder.attributes, attribute.Stringer(k, v))
}

func (span *span) MetadataLink(k *xopat.LinkAttribute, v xoptrace.Trace) {
	_, tmpSpan := span.logger.tracer.Start(span.ctx, k.Key(), oteltrace.WithLinks(
		oteltrace.Link{
			SpanContext: oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
				TraceID:    v.TraceID().Array(),
				SpanID:     v.SpanID().Array(),
				TraceFlags: oteltrace.TraceFlags(v.Flags().Array()[0]),
				TraceState: emptyTraceState, // TODO: is this right?
				Remote:     true,            // information not available
			}),
		},
	))
	tmpSpan.SetAttributes(spanIsLinkAttributeKey.Bool(true))
	tmpSpan.End()
}

func (builder *builder) Uint64(k string, v uint64, dt xopbase.DataType) {
	if dt == xopbase.Uint64DataType {
		builder.attributes = append(builder.attributes, attribute.String(k, strconv.FormatUint(v, 10)))
	} else {
		builder.attributes = append(builder.attributes, attribute.Int64(k, int64(v)))
	}
}

func (builder *builder) Bool(k string, v bool) {
	builder.attributes = append(builder.attributes, attribute.Bool(k, v))
}

func (builder *builder) Float64(k string, v float64, _ xopbase.DataType) {
	builder.attributes = append(builder.attributes, attribute.Float64(k, v))
}

func (builder *builder) Int64(k string, v int64, _ xopbase.DataType) {
	builder.attributes = append(builder.attributes, attribute.Int64(k, v))
}

func (builder *builder) String(k string, v string, _ xopbase.DataType) {
	builder.attributes = append(builder.attributes, attribute.String(k, v))
}

func (span *span) MetadataAny(k *xopat.AnyAttribute, v interface{}) {
	key := k.Key()
	enc, err := json.Marshal(v)
	var value string
	if err != nil {
		value = fmt.Sprintf("[zopotel] could not marshal %T value: %s", v, err)
	} else {
		value = string(enc)
	}
	if !k.Multiple() {
		if k.Locked() {
			span.lock.Lock()
			defer span.lock.Unlock()
			if span.hasPrior == nil {
				span.hasPrior = make(map[string]struct{})
			}
			if _, ok := span.hasPrior[key]; ok {
				return
			}
			span.hasPrior[key] = struct{}{}
		}
		span.span.SetAttributes(attribute.String(key, value))
		return
	}
	span.lock.Lock()
	defer span.lock.Unlock()
	if k.Distinct() {
		if span.metadataSeen == nil {
			span.metadataSeen = make(map[string]interface{})
		}
		seenRaw, ok := span.metadataSeen[key]
		if !ok {
			seen := make(map[string]struct{})
			span.metadataSeen[key] = seen
			seen[value] = struct{}{}
		} else {
			seen := seenRaw.(map[string]struct{})
			if _, ok := seen[value]; ok {
				return
			}
			seen[value] = struct{}{}
		}
	}
	if span.priorStringSlices == nil {
		span.priorStringSlices = make(map[string][]string)
	}
	s := span.priorStringSlices[key]
	s = append(s, value)
	span.priorStringSlices[key] = s
	span.span.SetAttributes(attribute.StringSlice(key, s))
}

func (span *span) MetadataBool(k *xopat.BoolAttribute, v bool) {
	key := k.Key()
	value := v
	if !k.Multiple() {
		if k.Locked() {
			span.lock.Lock()
			defer span.lock.Unlock()
			if span.hasPrior == nil {
				span.hasPrior = make(map[string]struct{})
			}
			if _, ok := span.hasPrior[key]; ok {
				return
			}
			span.hasPrior[key] = struct{}{}
		}
		span.span.SetAttributes(attribute.Bool(key, value))
		return
	}
	span.lock.Lock()
	defer span.lock.Unlock()
	if k.Distinct() {
		if span.metadataSeen == nil {
			span.metadataSeen = make(map[string]interface{})
		}
		seenRaw, ok := span.metadataSeen[key]
		if !ok {
			seen := make(map[bool]struct{})
			span.metadataSeen[key] = seen
			seen[value] = struct{}{}
		} else {
			seen := seenRaw.(map[bool]struct{})
			if _, ok := seen[value]; ok {
				return
			}
			seen[value] = struct{}{}
		}
	}
	if span.priorBoolSlices == nil {
		span.priorBoolSlices = make(map[string][]bool)
	}
	s := span.priorBoolSlices[key]
	s = append(s, value)
	span.priorBoolSlices[key] = s
	span.span.SetAttributes(attribute.BoolSlice(key, s))
}

func (span *span) MetadataEnum(k *xopat.EnumAttribute, v xopat.Enum) {
	key := k.Key()
	value := v.String()
	if !k.Multiple() {
		if k.Locked() {
			span.lock.Lock()
			defer span.lock.Unlock()
			if span.hasPrior == nil {
				span.hasPrior = make(map[string]struct{})
			}
			if _, ok := span.hasPrior[key]; ok {
				return
			}
			span.hasPrior[key] = struct{}{}
		}
		span.span.SetAttributes(attribute.String(key, value))
		return
	}
	span.lock.Lock()
	defer span.lock.Unlock()
	if k.Distinct() {
		if span.metadataSeen == nil {
			span.metadataSeen = make(map[string]interface{})
		}
		seenRaw, ok := span.metadataSeen[key]
		if !ok {
			seen := make(map[string]struct{})
			span.metadataSeen[key] = seen
			seen[value] = struct{}{}
		} else {
			seen := seenRaw.(map[string]struct{})
			if _, ok := seen[value]; ok {
				return
			}
			seen[value] = struct{}{}
		}
	}
	if span.priorStringSlices == nil {
		span.priorStringSlices = make(map[string][]string)
	}
	s := span.priorStringSlices[key]
	s = append(s, value)
	span.priorStringSlices[key] = s
	span.span.SetAttributes(attribute.StringSlice(key, s))
}

func (span *span) MetadataFloat64(k *xopat.Float64Attribute, v float64) {
	key := k.Key()
	value := v
	if !k.Multiple() {
		if k.Locked() {
			span.lock.Lock()
			defer span.lock.Unlock()
			if span.hasPrior == nil {
				span.hasPrior = make(map[string]struct{})
			}
			if _, ok := span.hasPrior[key]; ok {
				return
			}
			span.hasPrior[key] = struct{}{}
		}
		span.span.SetAttributes(attribute.Float64(key, value))
		return
	}
	span.lock.Lock()
	defer span.lock.Unlock()
	if k.Distinct() {
		if span.metadataSeen == nil {
			span.metadataSeen = make(map[string]interface{})
		}
		seenRaw, ok := span.metadataSeen[key]
		if !ok {
			seen := make(map[float64]struct{})
			span.metadataSeen[key] = seen
			seen[value] = struct{}{}
		} else {
			seen := seenRaw.(map[float64]struct{})
			if _, ok := seen[value]; ok {
				return
			}
			seen[value] = struct{}{}
		}
	}
	if span.priorFloat64Slices == nil {
		span.priorFloat64Slices = make(map[string][]float64)
	}
	s := span.priorFloat64Slices[key]
	s = append(s, value)
	span.priorFloat64Slices[key] = s
	span.span.SetAttributes(attribute.Float64Slice(key, s))
}

func (span *span) MetadataInt64(k *xopat.Int64Attribute, v int64) {
	key := k.Key()
	value := v
	if !k.Multiple() {
		if k.Locked() {
			span.lock.Lock()
			defer span.lock.Unlock()
			if span.hasPrior == nil {
				span.hasPrior = make(map[string]struct{})
			}
			if _, ok := span.hasPrior[key]; ok {
				return
			}
			span.hasPrior[key] = struct{}{}
		}
		span.span.SetAttributes(attribute.Int64(key, value))
		return
	}
	span.lock.Lock()
	defer span.lock.Unlock()
	if k.Distinct() {
		if span.metadataSeen == nil {
			span.metadataSeen = make(map[string]interface{})
		}
		seenRaw, ok := span.metadataSeen[key]
		if !ok {
			seen := make(map[int64]struct{})
			span.metadataSeen[key] = seen
			seen[value] = struct{}{}
		} else {
			seen := seenRaw.(map[int64]struct{})
			if _, ok := seen[value]; ok {
				return
			}
			seen[value] = struct{}{}
		}
	}
	if span.priorInt64Slices == nil {
		span.priorInt64Slices = make(map[string][]int64)
	}
	s := span.priorInt64Slices[key]
	s = append(s, value)
	span.priorInt64Slices[key] = s
	span.span.SetAttributes(attribute.Int64Slice(key, s))
}

func (span *span) MetadataString(k *xopat.StringAttribute, v string) {
	key := k.Key()
	value := v
	if !k.Multiple() {
		if k.Locked() {
			span.lock.Lock()
			defer span.lock.Unlock()
			if span.hasPrior == nil {
				span.hasPrior = make(map[string]struct{})
			}
			if _, ok := span.hasPrior[key]; ok {
				return
			}
			span.hasPrior[key] = struct{}{}
		}
		span.span.SetAttributes(attribute.String(key, value))
		return
	}
	span.lock.Lock()
	defer span.lock.Unlock()
	if k.Distinct() {
		if span.metadataSeen == nil {
			span.metadataSeen = make(map[string]interface{})
		}
		seenRaw, ok := span.metadataSeen[key]
		if !ok {
			seen := make(map[string]struct{})
			span.metadataSeen[key] = seen
			seen[value] = struct{}{}
		} else {
			seen := seenRaw.(map[string]struct{})
			if _, ok := seen[value]; ok {
				return
			}
			seen[value] = struct{}{}
		}
	}
	if span.priorStringSlices == nil {
		span.priorStringSlices = make(map[string][]string)
	}
	s := span.priorStringSlices[key]
	s = append(s, value)
	span.priorStringSlices[key] = s
	span.span.SetAttributes(attribute.StringSlice(key, s))
}

func (span *span) MetadataTime(k *xopat.TimeAttribute, v time.Time) {
	key := k.Key()
	value := v.Format(time.RFC3339Nano)
	if !k.Multiple() {
		if k.Locked() {
			span.lock.Lock()
			defer span.lock.Unlock()
			if span.hasPrior == nil {
				span.hasPrior = make(map[string]struct{})
			}
			if _, ok := span.hasPrior[key]; ok {
				return
			}
			span.hasPrior[key] = struct{}{}
		}
		span.span.SetAttributes(attribute.String(key, value))
		return
	}
	span.lock.Lock()
	defer span.lock.Unlock()
	if k.Distinct() {
		if span.metadataSeen == nil {
			span.metadataSeen = make(map[string]interface{})
		}
		seenRaw, ok := span.metadataSeen[key]
		if !ok {
			seen := make(map[string]struct{})
			span.metadataSeen[key] = seen
			seen[value] = struct{}{}
		} else {
			seen := seenRaw.(map[string]struct{})
			if _, ok := seen[value]; ok {
				return
			}
			seen[value] = struct{}{}
		}
	}
	if span.priorStringSlices == nil {
		span.priorStringSlices = make(map[string][]string)
	}
	s := span.priorStringSlices[key]
	s = append(s, value)
	span.priorStringSlices[key] = s
	span.span.SetAttributes(attribute.StringSlice(key, s))
}
