// Code generated by "enumer -type=EventType -linecomment -json -sql"; DO NOT EDIT.

package xoprecorder

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _EventTypeName = "linerequestStartrequestDonespanStartspanStartflushmetadatacustom"

var _EventTypeIndex = [...]uint8{0, 4, 16, 27, 36, 45, 50, 58, 64}

const _EventTypeLowerName = "linerequeststartrequestdonespanstartspanstartflushmetadatacustom"

func (i EventType) String() string {
	if i < 0 || i >= EventType(len(_EventTypeIndex)-1) {
		return fmt.Sprintf("EventType(%d)", i)
	}
	return _EventTypeName[_EventTypeIndex[i]:_EventTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _EventTypeNoOp() {
	var x [1]struct{}
	_ = x[LineEvent-(0)]
	_ = x[RequestStart-(1)]
	_ = x[RequestDone-(2)]
	_ = x[SpanStart-(3)]
	_ = x[SpanDone-(4)]
	_ = x[FlushEvent-(5)]
	_ = x[MetadataSet-(6)]
	_ = x[CustomEvent-(7)]
}

var _EventTypeValues = []EventType{LineEvent, RequestStart, RequestDone, SpanStart, SpanDone, FlushEvent, MetadataSet, CustomEvent}

var _EventTypeNameToValueMap = map[string]EventType{
	_EventTypeName[0:4]:        LineEvent,
	_EventTypeLowerName[0:4]:   LineEvent,
	_EventTypeName[4:16]:       RequestStart,
	_EventTypeLowerName[4:16]:  RequestStart,
	_EventTypeName[16:27]:      RequestDone,
	_EventTypeLowerName[16:27]: RequestDone,
	_EventTypeName[27:36]:      SpanStart,
	_EventTypeLowerName[27:36]: SpanStart,
	_EventTypeName[36:45]:      SpanDone,
	_EventTypeLowerName[36:45]: SpanDone,
	_EventTypeName[45:50]:      FlushEvent,
	_EventTypeLowerName[45:50]: FlushEvent,
	_EventTypeName[50:58]:      MetadataSet,
	_EventTypeLowerName[50:58]: MetadataSet,
	_EventTypeName[58:64]:      CustomEvent,
	_EventTypeLowerName[58:64]: CustomEvent,
}

var _EventTypeNames = []string{
	_EventTypeName[0:4],
	_EventTypeName[4:16],
	_EventTypeName[16:27],
	_EventTypeName[27:36],
	_EventTypeName[36:45],
	_EventTypeName[45:50],
	_EventTypeName[50:58],
	_EventTypeName[58:64],
}

// EventTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func EventTypeString(s string) (EventType, error) {
	if val, ok := _EventTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _EventTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to EventType values", s)
}

// EventTypeValues returns all values of the enum
func EventTypeValues() []EventType {
	return _EventTypeValues
}

// EventTypeStrings returns a slice of all String values of the enum
func EventTypeStrings() []string {
	strs := make([]string, len(_EventTypeNames))
	copy(strs, _EventTypeNames)
	return strs
}

// IsAEventType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i EventType) IsAEventType() bool {
	for _, v := range _EventTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for EventType
func (i EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for EventType
func (i *EventType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("EventType should be a string, got %s", data)
	}

	var err error
	*i, err = EventTypeString(s)
	return err
}

func (i EventType) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *EventType) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of EventType: %[1]T(%[1]v)", value)
	}

	val, err := EventTypeString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}
