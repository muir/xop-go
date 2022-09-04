// Code generated by "enumer -type=Level -linecomment -json -sql"; DO NOT EDIT.

package xopnum

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	_LevelName_0      = "debug"
	_LevelLowerName_0 = "debug"
	_LevelName_1      = "traceinfo"
	_LevelLowerName_1 = "traceinfo"
	_LevelName_2      = "warn"
	_LevelLowerName_2 = "warn"
	_LevelName_3      = "error"
	_LevelLowerName_3 = "error"
	_LevelName_4      = "alert"
	_LevelLowerName_4 = "alert"
)

var (
	_LevelIndex_0 = [...]uint8{0, 5}
	_LevelIndex_1 = [...]uint8{0, 5, 9}
	_LevelIndex_2 = [...]uint8{0, 4}
	_LevelIndex_3 = [...]uint8{0, 5}
	_LevelIndex_4 = [...]uint8{0, 5}
)

func (i Level) String() string {
	switch {
	case i == 5:
		return _LevelName_0
	case 8 <= i && i <= 9:
		i -= 8
		return _LevelName_1[_LevelIndex_1[i]:_LevelIndex_1[i+1]]
	case i == 13:
		return _LevelName_2
	case i == 17:
		return _LevelName_3
	case i == 20:
		return _LevelName_4
	default:
		return fmt.Sprintf("Level(%d)", i)
	}
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _LevelNoOp() {
	var x [1]struct{}
	_ = x[DebugLevel-(5)]
	_ = x[TraceLevel-(8)]
	_ = x[InfoLevel-(9)]
	_ = x[WarnLevel-(13)]
	_ = x[ErrorLevel-(17)]
	_ = x[AlertLevel-(20)]
}

var _LevelValues = []Level{DebugLevel, TraceLevel, InfoLevel, WarnLevel, ErrorLevel, AlertLevel}

var _LevelNameToValueMap = map[string]Level{
	_LevelName_0[0:5]:      DebugLevel,
	_LevelLowerName_0[0:5]: DebugLevel,
	_LevelName_1[0:5]:      TraceLevel,
	_LevelLowerName_1[0:5]: TraceLevel,
	_LevelName_1[5:9]:      InfoLevel,
	_LevelLowerName_1[5:9]: InfoLevel,
	_LevelName_2[0:4]:      WarnLevel,
	_LevelLowerName_2[0:4]: WarnLevel,
	_LevelName_3[0:5]:      ErrorLevel,
	_LevelLowerName_3[0:5]: ErrorLevel,
	_LevelName_4[0:5]:      AlertLevel,
	_LevelLowerName_4[0:5]: AlertLevel,
}

var _LevelNames = []string{
	_LevelName_0[0:5],
	_LevelName_1[0:5],
	_LevelName_1[5:9],
	_LevelName_2[0:4],
	_LevelName_3[0:5],
	_LevelName_4[0:5],
}

// LevelString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func LevelString(s string) (Level, error) {
	if val, ok := _LevelNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _LevelNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Level values", s)
}

// LevelValues returns all values of the enum
func LevelValues() []Level {
	return _LevelValues
}

// LevelStrings returns a slice of all String values of the enum
func LevelStrings() []string {
	strs := make([]string, len(_LevelNames))
	copy(strs, _LevelNames)
	return strs
}

// IsALevel returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Level) IsALevel() bool {
	for _, v := range _LevelValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Level
func (i Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Level
func (i *Level) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Level should be a string, got %s", data)
	}

	var err error
	*i, err = LevelString(s)
	return err
}

func (i Level) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *Level) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of Level: %[1]T(%[1]v)", value)
	}

	val, err := LevelString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}