// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file

package xop

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/xoplog/xop-go/internal/util/version"
	"github.com/xoplog/xop-go/xopat"
	"github.com/xoplog/xop-go/xopbase"
	"github.com/xoplog/xop-go/xoptrace"
)

// Seed is used to create a Log.
type Seed struct {
	spanSeed
	settings LogSettings
}

func (s Seed) String() string {
	return "SEED:" + s.spanSeed.String() + s.settings.String()
}

// Copy makes a deep copy of a seed and also randomizes
// the SpanID.
func (seed Seed) Copy(mods ...SeedModifier) Seed {
	n := Seed{
		spanSeed: seed.spanSeed.copy(true),
		settings: seed.settings.Copy(),
	}
	if !seed.spanSet {
		n.traceBundle.Trace.SpanID().SetRandom()
	}
	n = n.applyMods(mods)
	return n
}

// Seed provides a copy of the current span's seed, but the
// spanID is randomized.
func (span *Span) Seed(mods ...SeedModifier) Seed {
	n := Seed{
		spanSeed: span.seed.copy(false),
		settings: span.log.settings.Copy(),
	}
	if !span.seed.spanSet {
		n.traceBundle.Trace.SpanID().SetRandom()
	}
	n = n.applyMods(mods)
	return n
}

// SubSeed provides a copy of the current span's seed, but the
// spanID is randomized and the Parent set to the now prior
// Trace
func (span *Span) SubSeed(mods ...SeedModifier) Seed {
	n := Seed{
		spanSeed: span.seed.copy(false),
		settings: span.log.settings.Copy(),
	}
	n.traceBundle.Parent = n.traceBundle.Trace
	if !span.seed.spanSet {
		n.traceBundle.Trace.SpanID().SetRandom()
	}
	n = n.applyMods(mods)
	return n
}

// SeedReactiveCallback is used to modify seeds as they are just sprouting.
type SeedReactiveCallback func(ctx context.Context, seed Seed, nameOrDescription string, isChildSpan bool, now time.Time) []SeedModifier

type seedReactiveCallbacks []SeedReactiveCallback

func (cbs seedReactiveCallbacks) Copy() seedReactiveCallbacks {
	n := make(seedReactiveCallbacks, len(cbs))
	copy(n, cbs)
	return n
}

type spanSeed struct {
	traceBundle          xoptrace.Bundle
	spanSequenceCode     string
	description          string
	loggers              loggers
	config               Config
	flushDelay           time.Duration
	reactive             seedReactiveCallbacks
	ctx                  context.Context
	currentReactiveIndex int
	traceSet             bool
	spanSet              bool
	sourceInfo           xopbase.SourceInfo
}

func (s spanSeed) copy(withHistory bool) spanSeed {
	n := s
	n.loggers = s.loggers.Copy(withHistory)
	n.traceBundle = s.traceBundle.Copy()
	n.reactive = s.reactive.Copy()
	return n
}

// String is purely meant for debugging purposes and is not performant
func (s spanSeed) String() string {
	var str string
	if !s.traceBundle.Parent.IsZero() {
		str += " parent:" + s.traceBundle.Parent.String()
	}
	if s.traceSet {
		str += " trace:" + s.traceBundle.Trace.TraceID().String()
	}
	if s.spanSet {
		str += " span:" + s.traceBundle.Trace.SpanID().String()
	}
	if !s.traceBundle.Baggage.IsZero() {
		str += " baggage:" + s.traceBundle.Baggage.String()
	}
	if !s.traceBundle.State.IsZero() {
		str += " state:" + s.traceBundle.State.String()
	}
	if s.sourceInfo.Source != "" {
		str += " source:" + s.sourceInfo.Source
	}
	if s.sourceInfo.Namespace != "" {
		str += " namespace:" + s.sourceInfo.Namespace
	}
	return str
}

type SeedModifier func(*Seed)

func NewSeed(mods ...SeedModifier) Seed {
	seed := &Seed{
		spanSeed: spanSeed{
			config:      DefaultConfig,
			traceBundle: xoptrace.NewBundle(),
		},
		settings: DefaultSettings,
	}
	return seed.applyMods(mods)
}

func (seed Seed) applyMods(mods []SeedModifier) Seed {
	for _, mod := range mods {
		mod(&seed)
	}
	if seed.sourceInfo.Source == "" {
		seed.sourceInfo.Source, seed.sourceInfo.SourceVersion = version.SplitVersion(filepath.Base(os.Args[0]))
	}
	if seed.sourceInfo.SourceVersion == nil {
		seed.sourceInfo.SourceVersion = version.ZeroVersion
	}
	if seed.sourceInfo.Namespace == "" {
		seed.sourceInfo.Namespace, seed.sourceInfo.NamespaceVersion = version.SplitVersion(xopat.DefaultNamespace)
	}
	if seed.sourceInfo.NamespaceVersion == nil {
		seed.sourceInfo.NamespaceVersion = version.ZeroVersion
	}
	return seed
}

// Request creates a new top-level span (a request).  Use when
// starting something new, like receiving an http request or
// starting a cron job.
func (seed Seed) Request(descriptionOrName string) *Log {
	now := time.Now()
	seed = seed.react(true, descriptionOrName, now)
	return seed.request(descriptionOrName, now)
}

// SubSpan creates a new top-level span (a request) but
// initializes the span/trace data as if it were a subspan.
// The traceID must already be set.  Use this with caution,
// it is meant for handing off from spans created elsewhere.
func (seed Seed) SubSpan(descriptionOrName string) *Log {
	now := time.Now()
	seed = seed.react(false, descriptionOrName, now)
	return seed.request(descriptionOrName, now)
}

// WithReactive provides a callback that is used to modify seeds as they
// are in the process of sprouting.  Just as a seed is being used to create
// a request of sub-span, all reactive functions will be called.  Such
// functions must return a valid seed.  The seed they start with will be
// valid, so they can simply return that seed.
//
// When WithReactive is called from a SeedReactiveCallback, the new reactive
// function is only evaluated on descendent seeds.
func WithReactive(f SeedReactiveCallback) SeedModifier {
	return func(s *Seed) {
		s.reactive = append(s.reactive, f)
	}
}

// WithReactiveReplaced may only be used from within a call to a reactive
// function.  The current reactive function is the one that will be replaced.
// To remove a reactive function, call WithReactiveReplaced with nil.
func WithReactiveReplaced(f SeedReactiveCallback) SeedModifier {
	return func(s *Seed) {
		s.reactive[s.currentReactiveIndex] = f
	}
}

// WithContext puts a context into the seed.  That context will be
// passed through to the base layer Request and Seed functions.
func WithContext(ctx context.Context) SeedModifier {
	return func(s *Seed) {
		s.ctx = ctx
	}
}

func WithBundle(bundle xoptrace.Bundle) SeedModifier {
	return func(s *Seed) {
		s.traceBundle = bundle
	}
}

func WithSpan(spanID [8]byte) SeedModifier {
	return func(s *Seed) {
		s.traceBundle.Trace.SpanID().SetArray(spanID)
		s.spanSet = true
	}
}

func WithTrace(trace xoptrace.Trace) SeedModifier {
	return func(s *Seed) {
		s.traceBundle.Trace = trace
		s.traceSet = true
		s.spanSet = true
	}
}

func WithSettings(f func(*LogSettings)) SeedModifier {
	return func(s *Seed) {
		f(&s.settings)
	}
}

// WithNamespace serves to provide informationa about where logs
// will come from. The best value for the namespace is either the
// go module name for the overall project or the repository path.
// If Xop was logging, it should use "github.com/xoplog/xop-go"
// as the source. The source can cal be versioned. Split on space
// or dash (eg. "github.com/xoplog/xop-go v1.3.10")
//
// The namespace exists to establish the which semantic conventions are
// in use for whatever logging follows.
//
// If not specified, namespace will default to xopat.DefaultNamespace
func WithNamespace(namespace string) SeedModifier {
	return func(s *Seed) {
		s.sourceInfo.Namespace, s.sourceInfo.NamespaceVersion = version.SplitVersion(namespace)
	}
}

// WithSource specifies the source of the the logs. This should be the program
// name or the program family if there are a group of related programs.
//
// Versioning is supported.
//
// If not specified, filepath.Base(os.Argv[0]) will be used.
func WithSource(source string) SeedModifier {
	return func(s *Seed) {
		s.sourceInfo.Source, s.sourceInfo.SourceVersion = version.SplitVersion(source)
	}
}

func CombineSeedModifiers(mods ...SeedModifier) SeedModifier {
	return func(s *Seed) {
		for _, f := range mods {
			f(s)
		}
	}
}

func (seed Seed) Bundle() xoptrace.Bundle {
	return seed.traceBundle
}

func (seed Seed) SourceInfo() xopbase.SourceInfo {
	return seed.sourceInfo
}

func (seed Seed) react(isRequest bool, description string, now time.Time) Seed {
	if isRequest {
		if !seed.traceSet {
			seed.traceBundle.Trace.RebuildSetNonZero()
		}
	}
	seed.traceSet = false
	seed.spanSet = false
	if len(seed.reactive) == 0 {
		return seed
	}
	var nilSeen bool
	initialCount := len(seed.reactive)
	for i := 0; i < initialCount; i++ {
		f := seed.reactive[i]
		if f == nil {
			nilSeen = true
			i++
			continue
		}
		seed.currentReactiveIndex = i
		seed = seed.applyMods(f(seed.ctx, seed, description, !isRequest, now))
		if seed.reactive[i] == nil {
			nilSeen = true
		}
	}
	if nilSeen {
		reactive := make(seedReactiveCallbacks, 0, len(seed.reactive))
		for _, f := range seed.reactive {
			if f != nil {
				reactive = append(reactive, f)
			}
		}
		seed.reactive = reactive
	}
	return seed
}
