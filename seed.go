// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file

package xop

import (
	"context"
	"time"

	"github.com/muir/xop-go/trace"
)

// Seed is used to create a Log.
type Seed struct {
	spanSeed
	settings LogSettings
}

func (seed Seed) Copy(mods ...SeedModifier) Seed {
	return Seed{
		spanSeed: seed.spanSeed.copy(true),
		settings: seed.settings.Copy(),
	}
}

// SeedReactiveCallback is used to modify seeds as they are just sprouting
// The selfIndex parameter can be used with WithReactiveReplaced or
// WithReactiveRemoved.
type SeedReactiveCallback func(ctx context.Context, seed Seed, selfIndex int, nameOrDescription string, isChildSpan bool) Seed

type seedReactiveCallbacks []SeedReactiveCallback

func (cbs seedReactiveCallbacks) Copy() seedReactiveCallbacks {
	n := make(seedReactiveCallbacks, len(cbs))
	copy(n, cbs)
	return n
}

type spanSeed struct {
	traceBundle          trace.Bundle
	spanSequenceCode     string
	description          string
	loggers              loggers
	config               Config
	flushDelay           time.Duration
	reactive             seedReactiveCallbacks
	ctx                  context.Context
	currentReactiveIndex int
	reactiveReplaced     bool
}

func (s spanSeed) copy(withHistory bool) spanSeed {
	n := s
	n.loggers = s.loggers.Copy(withHistory)
	n.traceBundle = s.traceBundle.Copy()
	n.reactive = s.reactive.Copy()
	return n
}

type SeedModifier func(*Seed)

func NewSeed(mods ...SeedModifier) Seed {
	seed := &Seed{
		spanSeed: spanSeed{
			config:      DefaultConfig,
			traceBundle: trace.NewBundle(),
		},
		settings: DefaultSettings,
	}
	return seed.applyMods(mods)
}

// Seed provides a copy of the current span's seed, but the
// spanID is randomized.
func (span *Span) Seed(mods ...SeedModifier) Seed {
	seed := Seed{
		spanSeed: span.seed.copy(false),
		settings: span.log.settings.Copy(),
	}
	seed.spanSeed.traceBundle.Trace.RandomizeSpanID()
	return seed.applyMods(mods)
}

func (seed Seed) applyMods(mods []SeedModifier) Seed {
	for _, mod := range mods {
		mod(&seed)
	}
	return seed
}

// WithReactive provides a callback that is used to modify seeds as they
// are in the process of sprouting.  Just as a seed is being used to create
// a request of sub-span, all reactive functions will be called.  Such
// functions must return a valid seed.  The seed they start with will be
// valid, so they can simply return that seed.
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
		s.reactiveReplaced = true
	}
}

// WithContext puts a context into the seed.  That context will be
// passed through to the base layer Request and Seed functions.
func WithContext(ctx context.Context) SeedModifier {
	return func(s *Seed) {
		s.ctx = ctx
	}
}

func WithBundle(bundle trace.Bundle) SeedModifier {
	return func(s *Seed) {
		s.traceBundle = bundle
	}
}

func WithSpan(spanID [8]byte) SeedModifier {
	return func(s *Seed) {
		s.traceBundle.Trace.SpanID().Set(spanID)
	}
}

func WithTrace(trace trace.Trace) SeedModifier {
	return func(s *Seed) {
		s.traceBundle.Trace = trace
	}
}

func WithSettings(f func(*LogSettings)) SeedModifier {
	return func(s *Seed) {
		f(&s.settings)
	}
}

func CombineSeedModfiers(mods ...SeedModifier) SeedModifier {
	return func(s *Seed) {
		for _, f := range mods {
			f(s)
		}
	}
}

func (seed Seed) Bundle() trace.Bundle {
	return seed.traceBundle
}

func (seed Seed) react(isRequest bool, description string) Seed {
	if isRequest {
		seed.traceBundle.Trace.RebuildSetNonZero()
	} else {
		seed.traceBundle = seed.traceBundle.Copy()
		seed.traceBundle.Trace.RandomizeSpanID()
	}
	if len(seed.reactive) == 0 {
		return seed
	}
	var nilSeen bool
	for i := 0; i < len(seed.reactive); {
		f := seed.reactive[i]
		if f == nil {
			nilSeen = true
			i++
			continue
		}
		seed.currentReactiveIndex = i
		seed.reactiveReplaced = false
		seed = f(seed.ctx, seed, i, description, !isRequest)
		if !seed.reactiveReplaced {
			i++
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
