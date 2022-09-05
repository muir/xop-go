// Code generated by "enumer -type=DataType -linecomment -json -sql"; DO NOT EDIT.

package xopbase

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

const _DataTypeName = "EnumDataTypeAnyDataTypeBoolDataTypeDurationDataTypeErrorDataTypeFloat32DataTypeFloat64DataTypeIntDataTypeInt16DataTypeInt32DataTypeInt64DataTypeInt8DataTypeLinkDataTypeStringDataTypeTimeDataTypeUintDataTypeUint16DataTypeUint32DataTypeUint64DataTypeUint8DataType"

var _DataTypeIndex = [...]uint16{0, 12, 23, 35, 51, 64, 79, 94, 105, 118, 131, 144, 156, 168, 182, 194, 206, 220, 234, 248, 261}

const _DataTypeLowerName = "enumdatatypeanydatatypebooldatatypedurationdatatypeerrordatatypefloat32datatypefloat64datatypeintdatatypeint16datatypeint32datatypeint64datatypeint8datatypelinkdatatypestringdatatypetimedatatypeuintdatatypeuint16datatypeuint32datatypeuint64datatypeuint8datatype"

func (i DataType) String() string {
	if i < 0 || i >= DataType(len(_DataTypeIndex)-1) {
		return fmt.Sprintf("DataType(%d)", i)
	}
	return _DataTypeName[_DataTypeIndex[i]:_DataTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _DataTypeNoOp() {
	var x [1]struct{}
	_ = x[EnumDataType-(0)]
	_ = x[AnyDataType-(1)]
	_ = x[BoolDataType-(2)]
	_ = x[DurationDataType-(3)]
	_ = x[ErrorDataType-(4)]
	_ = x[Float32DataType-(5)]
	_ = x[Float64DataType-(6)]
	_ = x[IntDataType-(7)]
	_ = x[Int16DataType-(8)]
	_ = x[Int32DataType-(9)]
	_ = x[Int64DataType-(10)]
	_ = x[Int8DataType-(11)]
	_ = x[LinkDataType-(12)]
	_ = x[StringDataType-(13)]
	_ = x[TimeDataType-(14)]
	_ = x[UintDataType-(15)]
	_ = x[Uint16DataType-(16)]
	_ = x[Uint32DataType-(17)]
	_ = x[Uint64DataType-(18)]
	_ = x[Uint8DataType-(19)]
}

var _DataTypeValues = []DataType{EnumDataType, AnyDataType, BoolDataType, DurationDataType, ErrorDataType, Float32DataType, Float64DataType, IntDataType, Int16DataType, Int32DataType, Int64DataType, Int8DataType, LinkDataType, StringDataType, TimeDataType, UintDataType, Uint16DataType, Uint32DataType, Uint64DataType, Uint8DataType}

var _DataTypeNameToValueMap = map[string]DataType{
	_DataTypeName[0:12]:         EnumDataType,
	_DataTypeLowerName[0:12]:    EnumDataType,
	_DataTypeName[12:23]:        AnyDataType,
	_DataTypeLowerName[12:23]:   AnyDataType,
	_DataTypeName[23:35]:        BoolDataType,
	_DataTypeLowerName[23:35]:   BoolDataType,
	_DataTypeName[35:51]:        DurationDataType,
	_DataTypeLowerName[35:51]:   DurationDataType,
	_DataTypeName[51:64]:        ErrorDataType,
	_DataTypeLowerName[51:64]:   ErrorDataType,
	_DataTypeName[64:79]:        Float32DataType,
	_DataTypeLowerName[64:79]:   Float32DataType,
	_DataTypeName[79:94]:        Float64DataType,
	_DataTypeLowerName[79:94]:   Float64DataType,
	_DataTypeName[94:105]:       IntDataType,
	_DataTypeLowerName[94:105]:  IntDataType,
	_DataTypeName[105:118]:      Int16DataType,
	_DataTypeLowerName[105:118]: Int16DataType,
	_DataTypeName[118:131]:      Int32DataType,
	_DataTypeLowerName[118:131]: Int32DataType,
	_DataTypeName[131:144]:      Int64DataType,
	_DataTypeLowerName[131:144]: Int64DataType,
	_DataTypeName[144:156]:      Int8DataType,
	_DataTypeLowerName[144:156]: Int8DataType,
	_DataTypeName[156:168]:      LinkDataType,
	_DataTypeLowerName[156:168]: LinkDataType,
	_DataTypeName[168:182]:      StringDataType,
	_DataTypeLowerName[168:182]: StringDataType,
	_DataTypeName[182:194]:      TimeDataType,
	_DataTypeLowerName[182:194]: TimeDataType,
	_DataTypeName[194:206]:      UintDataType,
	_DataTypeLowerName[194:206]: UintDataType,
	_DataTypeName[206:220]:      Uint16DataType,
	_DataTypeLowerName[206:220]: Uint16DataType,
	_DataTypeName[220:234]:      Uint32DataType,
	_DataTypeLowerName[220:234]: Uint32DataType,
	_DataTypeName[234:248]:      Uint64DataType,
	_DataTypeLowerName[234:248]: Uint64DataType,
	_DataTypeName[248:261]:      Uint8DataType,
	_DataTypeLowerName[248:261]: Uint8DataType,
}

var _DataTypeNames = []string{
	_DataTypeName[0:12],
	_DataTypeName[12:23],
	_DataTypeName[23:35],
	_DataTypeName[35:51],
	_DataTypeName[51:64],
	_DataTypeName[64:79],
	_DataTypeName[79:94],
	_DataTypeName[94:105],
	_DataTypeName[105:118],
	_DataTypeName[118:131],
	_DataTypeName[131:144],
	_DataTypeName[144:156],
	_DataTypeName[156:168],
	_DataTypeName[168:182],
	_DataTypeName[182:194],
	_DataTypeName[194:206],
	_DataTypeName[206:220],
	_DataTypeName[220:234],
	_DataTypeName[234:248],
	_DataTypeName[248:261],
}

// DataTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DataTypeString(s string) (DataType, error) {
	if val, ok := _DataTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _DataTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DataType values", s)
}

// DataTypeValues returns all values of the enum
func DataTypeValues() []DataType {
	return _DataTypeValues
}

// DataTypeStrings returns a slice of all String values of the enum
func DataTypeStrings() []string {
	strs := make([]string, len(_DataTypeNames))
	copy(strs, _DataTypeNames)
	return strs
}

// IsADataType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i DataType) IsADataType() bool {
	for _, v := range _DataTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for DataType
func (i DataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for DataType
func (i *DataType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("DataType should be a string, got %s", data)
	}

	var err error
	*i, err = DataTypeString(s)
	return err
}

func (i DataType) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *DataType) Scan(value interface{}) error {
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
		return fmt.Errorf("invalid value of DataType: %[1]T(%[1]v)", value)
	}

	val, err := DataTypeString(str)
	if err != nil {
		return err
	}

	*i = val
	return nil
}