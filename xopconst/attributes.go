// This file is generated, DO NOT EDIT.  It comes from the corresponding .zzzgo file
package xopconst

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/muir/xoplog/trace"
)

// TODO: PERFORMANCE: pre-allocate blocks of 128 Attributes to provide better locality of reference when using these
// TODO: UintAttribute?
// TODO: TableAttribute?
// TODO: URLAttribute?

// Attribute represents an "any" attribute for a span.
type Attribute struct {
	properties   Make
	number       int
	jsonKey      string
	exampleValue interface{}
	reflectType  reflect.Type
	typeName     string
}

// DefaultNamespace sets the namespace for attribute names
// used to describe spans.  Within a namespace, the use of
// attribute names should be consistent.  If not specified,
// the name of the running program (os.Args[0]) is used.
// A better value is to set the namespace to be name of the
// code repository.  DefaultNamespace can be overridden in
// an init() function.
var DefaultNamespace = os.Args[0]

type Make struct {
	Key          string // the attribute name
	Description  string // the attribute description
	Namespace    string // the namespace for this attribute (otherwise DefaultNamespace is used)
	Indexed      bool   // hint: this attribute should be indexed
	Prominence   int    // hint: how important is this attribute (lower is more important)
	Multiple     bool   // keep all values if the attribute is given multiple times
	Distinct     bool   // when keeping all values, only keep distinct values (not supported for interface{})
	Ranged       bool   // hint: comparisons between values are meaningful (eg: time, integers)
	Locked       bool   // only keep the first value
	AppenderFunc func() ArrayAppender
}

type ArrayAppender interface {
	AppendAny(interface{})
}

var (
	lock            sync.RWMutex
	registeredNames = make(map[string]*Attribute)
	allAttributes   []*Attribute
)

// Can't use MACRO for these since default values are needed

func (s Make) DurationAttribute() *DurationAttribute {
	return &DurationAttribute{Attribute: s.attribute(time.Duration(0), nil)}
}

func (s Make) TryDurationAttribute() (_ *DurationAttribute, err error) {
	return &DurationAttribute{Attribute: s.attribute(time.Duration(0), &err)}, err
}

func (s Make) LinkAttribute() *LinkAttribute {
	return &LinkAttribute{Attribute: s.attribute(trace.Trace{}, nil)}
}

func (s Make) TryLinkAttribute() (_ *LinkAttribute, err error) {
	return &LinkAttribute{Attribute: s.attribute(trace.Trace{}, &err)}, err
}

func (s Make) StrAttribute() *StrAttribute {
	return &StrAttribute{Attribute: s.attribute("", nil)}
}

func (s Make) TryStrAttribute() (_ *StrAttribute, err error) {
	return &StrAttribute{Attribute: s.attribute("", &err)}, err
}

func (s Make) IntAttribute() *IntAttribute {
	return &IntAttribute{Attribute: s.attribute(int(0), nil)}
}

func (s Make) TryIntAttribute() (_ *IntAttribute, err error) {
	return &IntAttribute{Attribute: s.attribute(int(0), &err)}, err
}

func (s Make) BoolAttribute() *BoolAttribute {
	return &BoolAttribute{Attribute: s.attribute(false, nil)}
}

func (s Make) TryBoolAttribute() (_ *BoolAttribute, err error) {
	return &BoolAttribute{Attribute: s.attribute(false, &err)}, err
}

func (s Make) TimeAttribute() *TimeAttribute {
	return &TimeAttribute{Attribute: s.attribute(time.Time{}, nil)}
}

func (s Make) TryTimeAttribute() (_ *TimeAttribute, err error) {
	return &TimeAttribute{Attribute: s.attribute(time.Time{}, &err)}, err
}

func (s Make) AnyAttribute(exampleValue interface{}) *AnyAttribute {
	return &AnyAttribute{Attribute: s.attribute(exampleValue, nil)}
}

func (s Make) TryAnyAttribute(exampleValue interface{}) (_ *AnyAttribute, err error) {
	return &AnyAttribute{Attribute: s.attribute(exampleValue, &err)}, err
}

func (s Make) attribute(exampleValue interface{}, ep *error) Attribute {
	a, err := s.make(exampleValue)
	if err != nil {
		if ep == nil {
			panic(err)
		}
		*ep = err
	}
	return a
}

func (s Make) make(exampleValue interface{}) (Attribute, error) {
	lock.Lock()
	defer lock.Unlock()
	if prior, ok := registeredNames[s.Key]; ok {
		return *prior, fmt.Errorf("duplicate attribute registration for '%s'", s.Key)
	}
	if s.Namespace == "" {
		s.Namespace = DefaultNamespace
	}
	ra := Attribute{
		properties:   s,
		exampleValue: exampleValue,
		reflectType:  reflect.TypeOf(exampleValue),
		typeName:     fmt.Sprintf("%T", exampleValue),
	}
	jsonKey, err := json.Marshal(s.Key)
	if err != nil {
		return ra, fmt.Errorf("cannot marshal attribute name '%s': %w", s.Key, err)
	}
	ra.jsonKey = string(jsonKey)
	registeredNames[s.Key] = &ra
	allAttributes = append(allAttributes, &ra)
	return ra, nil
}

// JSONKey returns a JSON-quoted string that can be used as the key to the
// name of the attribute.  It is a string because []byte is mutable
func (r Attribute) JSONKey() string { return r.jsonKey }

// ReflectType can be nil if the example value was nil
func (r Attribute) ReflectType() reflect.Type { return r.reflectType }

func (r Attribute) Key() string               { return r.properties.Key }
func (r Attribute) Description() string       { return r.properties.Description }
func (r Attribute) Namespace() string         { return r.properties.Namespace }
func (r Attribute) Indexed() bool             { return r.properties.Indexed }
func (r Attribute) Multiple() bool            { return r.properties.Multiple }
func (r Attribute) Ranged() bool              { return r.properties.Ranged }
func (r Attribute) Locked() bool              { return r.properties.Locked }
func (r Attribute) Distinct() bool            { return r.properties.Distinct }
func (r Attribute) Prominence() int           { return r.properties.Prominence }
func (r Attribute) RegistrationNumber() int   { return r.number }
func (r Attribute) ExampleValue() interface{} { return r.exampleValue }
func (r Attribute) TypeName() string          { return r.typeName }

type (
	AnyAttribute      struct{ Attribute }
	BoolAttribute     struct{ Attribute }
	DurationAttribute struct{ Attribute }
	IntAttribute      struct{ Attribute }
	LinkAttribute     struct{ Attribute }
	StrAttribute      struct{ Attribute }
	TimeAttribute     struct{ Attribute }
)