// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file

package xopconsole

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/xoplog/xop-go/internal/util/version"
	"github.com/xoplog/xop-go/xopbase"
	"github.com/xoplog/xop-go/xopnum"
	"github.com/xoplog/xop-go/xopproto"
	"github.com/xoplog/xop-go/xoptrace"
	"github.com/xoplog/xop-go/xoputil"
	"github.com/xoplog/xop-go/xoputil/replayutil"

	"github.com/pkg/errors"
)

type replayData struct {
	lineCount   int
	currentLine string
	errors      []error
	spans       map[xoptrace.HexBytes8]*replaySpan
	requests    map[xoptrace.HexBytes8]*replayRequest
	dest        xopbase.Logger
	attributes  *replayutil.GlobalAttributeDefinitions
}

type replayRequest struct {
	replayData
	ts                  time.Time
	trace               xoptrace.Trace
	version             int64
	name                string
	sourceAndVersion    string
	namespaceAndVersion string
	request             xopbase.Request // XXX
	requestAttributes   *replayutil.RequestAttributeDefinitons
}

type replaySpan struct {
	replayData
	request *replayRequest
	span    xopbase.Span
}

type replayLine struct {
	replayData
	ts         time.Time
	spanID     xoptrace.HexBytes8
	level      xopnum.Level
	message    string
	stack      []runtime.Frame
	line       xopbase.Line
	attributes []func(xopbase.Line)
}

// xop alert 2023-05-31T22:20:09.200456-07:00 72b09846e8ed0099 "like a rock\"\\<'\n\r\t\b\x00" frightening=stuff STACK: /Users/sharnoff/src/github.com/muir/xop-go/xoptest/xoptestutil/cases.go:39 /Users/sharnoff/src/github.com/muir/xop-go/xopconsole/replay_test.go:43 /usr/local/Cellar/go/1.20.1/libexec/src/testing/testing.go:1576
func (x replayLine) replayLine(ctx context.Context, t string) error {
	var err error
	x.ts, t, err = oneTime(t)
	if err != nil {
		return err
	}
	spanIDString, _, t := oneWord(t, " ")
	if spanIDString == "" {
		return fmt.Errorf("missing idString")
	}
	spanID := xoptrace.NewHexBytes8FromString(spanIDString)
	spanData, ok := x.spans[spanID]
	if !ok {
		return fmt.Errorf("missing span %s", spanIDString)
	}
	message, t := oneStringAndSpace(t)
	for {
		key, sep, t := oneWord(t, "=:")
		switch sep {
		case ':':
			if key != "STACK" {
				return fmt.Errorf("invalid stack indicator")
			}
			for {
				file, _, t := oneWord(t, ":")
				if file == "" {
					return fmt.Errorf("invalid stack frame")
				}
				lineNum, sep, t := oneWord(t, " ")
				if lineNum == "" {
					return fmt.Errorf("invalid stack frame, line")
				}
				num, err := strconv.ParseInt(lineNum, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid stack frame, line num: %w", err)
				}
				x.stack = append(x.stack, runtime.Frame{
					File: file,
					Line: int(num),
				})
				if sep == '\000' {
					break
				}
			}
			break
		case '=':
			if len(t) == 0 {
				return fmt.Errorf("empty value")
			}
			if t[0] == '(' {
				var lengthString string
				lengthString, _, t = oneWord(t, ")")
				length, err := strconv.ParseUint(lengthString, 10, 64)
				if err != nil {
					return fmt.Errorf("parse model length: %w", err)
				}
				if len(t) < int(length)+2 {
					return fmt.Errorf("expected remaining string to be at least %d bytes", length+2)
				}
				encoded := t[:length]
				if t[length] != '/' {
					return fmt.Errorf("malformed model")
				}
				t = t[length+1:]
				var typ string
				var sep byte
				var encoding string
				encoding, sep, t = oneWord(t, "/")
				if typ == "" {
					return fmt.Errorf("missing model type")
				}
				ma := xopbase.ModelArg{
					Encoded: []byte(encoded),
				}
				if en, ok := xopproto.Encoding_value[encoding]; ok {
					ma.Encoding = xopproto.Encoding(en)
				} else {
					return errors.Errorf("invalid encoding (%s) when decoding attribute", encoding)
				}
				ma.ModelType, sep, t = oneWord(t, " ")
				if ma.ModelType == "" {
					return errors.Errorf("empty model type")
				}
				x.attributes = append(x.attributes, func(line xopbase.Line) { line.Any(key, ma) })
				if sep == '\000' {
					break
				}
				continue
			}
			var value string
			var sep byte
			value, sep, t = oneWord(t, " (/") // )
			switch sep {
			case '(':
				i := strings.IndexByte(t, ')')
				if i == -1 {
					return fmt.Errorf("invalid type specifier")
				}
				typ := t[:i]
				t = t[i+1:]
				switch typ {
				case "dur":
					dur, err := time.ParseDuration(value)
					if err != nil {
						return fmt.Errorf("invalid duration: %w", err)
					}
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.Duration(key, dur) })
				case "f32":
					f, err := strconv.ParseFloat(value, 32)
					if err != nil {
						return fmt.Errorf("invalid float: %w", err)
					}
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.Float64(key, f, xopbase.Float32DataType) })
				case "f64":
					f, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return fmt.Errorf("invalid float: %w", err)
					}
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.Float64(key, f, xopbase.Float64DataType) })
				case "string":
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.String(key, value, xopbase.StringDataType) })
				case "stringer":
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.String(key, value, xopbase.StringerDataType) })
				case "i8", "i16", "i32", "i64", "int":
					i, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid int: %w", err)
					}
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.Int64(key, i, xopbase.StringToDataType[typ]) })
				case "u8", "u16", "u32", "u64", "uint", "uintptr":
					i, err := strconv.ParseUint(value, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid uint: %w", err)
					}
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.Uint64(key, i, xopbase.StringToDataType[typ]) })
				case "time":
					ts, err := time.Parse(time.RFC3339Nano, value)
					if err != nil {
						return fmt.Errorf("invalid time: %w", err)
					}
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.Time(key, ts) })
				default:
					return fmt.Errorf("invalid type: %s", typ)
				}
			case ' ', '\000':
				if value == "" {
					return errors.Errorf("invalid value")
				}
				switch value[0] {
				case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					i, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid int: %w", err)
					}
					x.attributes = append(x.attributes, func(line xopbase.Line) { line.Int64(key, i, xopbase.IntDataType) })
				default:
					switch value {
					case "t":
						x.attributes = append(x.attributes, func(line xopbase.Line) { line.Bool(key, true) })
					case "f":
						x.attributes = append(x.attributes, func(line xopbase.Line) { line.Bool(key, false) })
					default:
						x.attributes = append(x.attributes, func(line xopbase.Line) { line.String(key, value, xopbase.StringDataType) })
					}
				}
				if sep == '\000' {
					break
				}
				continue
			case '/':
				// XXX enum: int/text
			default:
				// error
			}
		default:
			return fmt.Errorf("invalid input")
		}
	}
	line := spanData.span.NoPrefill().Line(x.level, x.ts, x.stack)
	for _, af := range x.attributes {
		af(line)
	}
	// XXX Model
	// XXX Link
	// XXX Template
	line.Msg(message)
	return nil
}

// XXX
// Example: xop Span 2023-06-03T12:22:47.699766-07:00 Start c253fd02cd66f874 5f23a4838a2c7205 "a fork one span" T1.1.A
func (x replayData) replaySpan(ctx context.Context, t string) error {
	var err error
	var ts time.Time
	ts, t, err = oneTime(t)
	if err != nil {
		return err
	}
	var n string
	n, _, t = oneWord(t, " ")
	if n == "" {
		return errors.Errorf("invalid span start")
	}
	var spanIDString string
	var sep byte
	spanIDString, sep, t = oneWord(t, " ")
	if spanIDString == "" {
		return errors.Errorf("invalid span spanID")
	}
	spanID := xoptrace.NewHexBytes8FromString(spanIDString)
	if n == "Start" {
		var parentIDString string
		parentIDString, _, t = oneWord(t, " ")
		if parentIDString == "" {
			return errors.Errorf("invalid span parentID")
		}
		parentID := xoptrace.NewHexBytes8FromString(parentIDString)
		parentSpan, ok := x.spans[parentID]
		if !ok {
			return errors.Errorf("%s span %s is missing parent %s", n, spanIDString, parentIDString)
		}
		request, err := x.getRequest(parentSpan)
		if err != nil {
			return err
		}
		bundle := request.bundle.Copy()
		bundle.Trace.SpanID().Set(spanID)
		var name string
		name, _, t = oneWord(t, " ")
		var spanSeqCode string
		spanSeqCode, _, t = oneWord(t, " ")
		if spanSeqCode == "" {
			return errors.Errorf("invalid span sequence code")
		}
		span := parentSpan.Span(ctx, ts, bundle, name, spanSeqCode)
		x.spans[spanID] = &spanData{
			replayData: x,
			span:       span,
			request:    request,
		}
		return nil
	}
	if n[0] != 'v' {
		return errors.Errorf("invalid span numbering")
	}
	v, err := strconv.ParseUint(n[1:], 10, 64)
	if err != nil {
		return errors.Errorf("invalid span numbering: %w", err)
	}
	spanData, ok := x.spans[spanID]
	if !ok {
		return errors.Errof("span id %s not found", spanIDString)
	}
	err = x.collectMetadata(spanData, sep, t)
	if err != nil {
		return err
	}
	spanData.span.Done(ts, false)
	return nil
}

func (x *requestData) collectMetdata(spanData *spanData, sep string, t string) error {
	var key string
	for ; sep != '\000'; key, sep, t = oneWord(t, "=") {
		var string value
		value, sep, t = oneWord(t, " ")
		def, ok := spanData.request.definitions[key]
		if !ok {
			return errors.Errorf("missing definition for %s", key)
		}
		switch def.Type {
		}
	}
}

func (x replayData) replayDef(ctx context.Context, t string) error {
	return nil
}

// so far: xop Request
// this func: timestamp "Start1" or "vNNN"
func (x replayData) replayRequest(ctx context.Context, t string) error {
	ts, t, err := oneTime(t)
	if err != nil {
		return err
	}
	var n string
	n, _, t = oneWord(t, " ")
	switch n {
	case "":
		return errors.Errorf("invalid request")
	case "Start1":
		return replayRequest{
			replayData: x,
			ts:         ts,
		}.replayRequestStart(ctx, t)
	default:
		if !strings.HasPrefix(n, "v") {
			return errors.Errorf("invalid request with %s", t)
		}
		v, err := strconv.ParseInt(n[1:], 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid request, invalid version number")
		}
		var requestIDString string
		requestIDString, sep, t = oneWord(t, " ")
		requestID := xoptrace.NewHexBytes8FromString(parentIDString)
		y, ok := x.requests[requestID]
		if !ok {
			return errors.Errorf("update to request %s that doesn't exist", requestIDString)
		}
		// XXX
	}
}

// so far: xop Request timestamp Start1
// this func: trace-headder request-name source+version namespace+version
// xop Request 2023-06-02T22:35:26.81344-07:00 Start1 00-d456604ffc88ac5f4f971afbfce39cda-8fd4b01b0c7684a5-01 TestReplayConsole/one-span xopconsole.test-0.0.0 xopconsole.test-0.0.0
func (x replayRequest) replayRequestStart(ctx context.Context, t string) error {
	th, _, t := oneWord(t, " ")
	if th == "" {
		return errors.Errorf("missing trace header")
	}
	var ok bool
	x.trace, ok = xoptrace.TraceFromString(th)
	if !ok {
		return errors.Errorf("invalid trace header")
	}
	x.name, t = oneStringAndSpace(t)
	if x.name == "" {
		return errors.Errorf("missing request name")
	}
	x.sourceAndVersion, t = oneStringAndSpace(t)
	if x.sourceAndVersion == "" {
		return errors.Errorf("missing source+version, trace is %s/%s, name is %s, remaining is %s", th, x.trace, x.name, t)
	}
	x.namespaceAndVersion, t = oneStringAndSpace(t)
	if x.namespaceAndVersion == "" {
		return errors.Errorf("missing namespace+version, remaining is %s", t)
	}
	// XXX baggage
	// XXX span
	// XXX parent
	bundle := xoptrace.Bundle{
		Trace: x.trace,
	}
	ns, nsVers := version.SplitVersion(x.namespaceAndVersion)
	so, soVers := version.SplitVersion(x.sourceAndVersion)
	sourceInfo := xopbase.SourceInfo{
		Source:           so,
		SourceVersion:    soVers,
		Namespace:        ns,
		NamespaceVersion: nsVers,
	}
	request := x.dest.Request(ctx, x.ts, bundle, x.name, sourceInfo)
	x.request = request
	x.requestAttributes = x.attributes.NewRequestAttributeDefinitions(bundle.Trace.SpanID().String())
	x.spans[bundle.Trace.GetSpanID()] = &spanData{
		replayData: x.replayData,
		span:       request,
		request:    &x,
	}
	x.requests[bundle.Trace.GetSpanID()] = &x
	return nil
}

func oneStringAndSpace(t string) (string, string) {
	a, b := oneString(t)
	if a == "" {
		return a, b
	}
	if len(b) > 0 && b[0] == ' ' {
		return a, b[1:]
	}
	return a, b
}

// oneString reads a possibly-quoted string
func oneString(t string) (string, string) {
	if len(t) == 0 {
		return "", ""
	}
	if t[0] == '"' {
		for i := 1; i < len(t); i++ {
			switch t[i] {
			case '\\':
				if i < len(t) {
					i++
				}
			case '"':
				one, err := strconv.Unquote(t[0 : i+1])
				if err != nil {
					return "", t
				}
				return one, t[i+1:]
			}
		}
	}
	one := xoputil.UnquotedConsoleStringRE.FindString(t)
	if one != "" {
		return one, t[len(one):]
	}
	return "", t
}

func oneTime(t string) (time.Time, string, error) {
	w, _, t := oneWord(t, " ")
	ts, err := time.Parse(time.RFC3339, w)
	return ts, t, err
}

func oneWord(t string, boundary string) (string, byte, string) {
	i := strings.IndexAny(t, boundary)
	switch i {
	case -1:
		return "", '\000', t
	case 0:
		return "", t[0], t[1:]
	}
	return t[:i], t[i], t[i+1:]
}

func Replay(ctx context.Context, inputStream io.Reader, dest xopbase.Logger) error {
	scanner := bufio.NewScanner(inputStream)
	x := replayData{
		dest:       dest,
		spans:      make(map[xoptrace.HexBytes8]xopbase.Span),
		requests:   make(map[xoptrace.HexBytes8]*replayRequest),
		attributes: replayutil.NewGlobalAttributeDefinitions(),
	}
	for scanner.Scan() {
		x.lineCount++
		t := scanner.Text()
		if !strings.HasPrefix(t, "xop ") {
			continue
		}
		x.currentLine = t
		t = t[len("xop "):]
		kind, _, t := oneWord(t, " ")
		var err error
		switch kind {
		case "Request":
			err = x.replayRequest(ctx, t)
		case "Span":
			err = x.replaySpan(ctx, t)
		case "Def":
			err = x.replayDef(ctx, t)
		case "alert":
			err = replayLine{
				replayData: x,
				level:      xopnum.AlertLevel,
			}.replayLine(ctx, t)
		case "debug":
			err = replayLine{
				replayData: x,
				level:      xopnum.DebugLevel,
			}.replayLine(ctx, t)
		case "error":
			err = replayLine{
				replayData: x,
				level:      xopnum.ErrorLevel,
			}.replayLine(ctx, t)
		case "info":
			err = replayLine{
				replayData: x,
				level:      xopnum.InfoLevel,
			}.replayLine(ctx, t)
		case "trace":
			err = replayLine{
				replayData: x,
				level:      xopnum.TraceLevel,
			}.replayLine(ctx, t)
		case "warn":
			err = replayLine{
				replayData: x,
				level:      xopnum.WarnLevel,
			}.replayLine(ctx, t)

			// prior line must be blank
		default:
			err = fmt.Errorf("invalid kind designator '%s'", kind)
		}
		if err != nil {
			x.errors = append(x.errors, errors.Wrapf(err, "line %d: %s", x.lineCount, x.currentLine))
		}
	}
	if len(x.errors) != 0 {
		// TODO: use a multi-error
		return x.errors[0]
	}
	return nil
}