// Code generated by "enumer -type=ProcessPriority -output enums.go"; DO NOT EDIT.

package engine

import (
	"fmt"
	"strings"
)

const _ProcessPriorityName = "BeforeMergeLowBeforeMergeBeforeMergeHighBeforeMergeFinalAfterMergeLowAfterMergeAfterMergeHighAfterMergeFinal"

var _ProcessPriorityIndex = [...]uint8{0, 14, 25, 40, 56, 69, 79, 93, 108}

const _ProcessPriorityLowerName = "beforemergelowbeforemergebeforemergehighbeforemergefinalaftermergelowaftermergeaftermergehighaftermergefinal"

func (i ProcessPriority) String() string {
	if i < 0 || i >= ProcessPriority(len(_ProcessPriorityIndex)-1) {
		return fmt.Sprintf("ProcessPriority(%d)", i)
	}
	return _ProcessPriorityName[_ProcessPriorityIndex[i]:_ProcessPriorityIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _ProcessPriorityNoOp() {
	var x [1]struct{}
	_ = x[BeforeMergeLow-(0)]
	_ = x[BeforeMerge-(1)]
	_ = x[BeforeMergeHigh-(2)]
	_ = x[BeforeMergeFinal-(3)]
	_ = x[AfterMergeLow-(4)]
	_ = x[AfterMerge-(5)]
	_ = x[AfterMergeHigh-(6)]
	_ = x[AfterMergeFinal-(7)]
}

var _ProcessPriorityValues = []ProcessPriority{BeforeMergeLow, BeforeMerge, BeforeMergeHigh, BeforeMergeFinal, AfterMergeLow, AfterMerge, AfterMergeHigh, AfterMergeFinal}

var _ProcessPriorityNameToValueMap = map[string]ProcessPriority{
	_ProcessPriorityName[0:14]:        BeforeMergeLow,
	_ProcessPriorityLowerName[0:14]:   BeforeMergeLow,
	_ProcessPriorityName[14:25]:       BeforeMerge,
	_ProcessPriorityLowerName[14:25]:  BeforeMerge,
	_ProcessPriorityName[25:40]:       BeforeMergeHigh,
	_ProcessPriorityLowerName[25:40]:  BeforeMergeHigh,
	_ProcessPriorityName[40:56]:       BeforeMergeFinal,
	_ProcessPriorityLowerName[40:56]:  BeforeMergeFinal,
	_ProcessPriorityName[56:69]:       AfterMergeLow,
	_ProcessPriorityLowerName[56:69]:  AfterMergeLow,
	_ProcessPriorityName[69:79]:       AfterMerge,
	_ProcessPriorityLowerName[69:79]:  AfterMerge,
	_ProcessPriorityName[79:93]:       AfterMergeHigh,
	_ProcessPriorityLowerName[79:93]:  AfterMergeHigh,
	_ProcessPriorityName[93:108]:      AfterMergeFinal,
	_ProcessPriorityLowerName[93:108]: AfterMergeFinal,
}

var _ProcessPriorityNames = []string{
	_ProcessPriorityName[0:14],
	_ProcessPriorityName[14:25],
	_ProcessPriorityName[25:40],
	_ProcessPriorityName[40:56],
	_ProcessPriorityName[56:69],
	_ProcessPriorityName[69:79],
	_ProcessPriorityName[79:93],
	_ProcessPriorityName[93:108],
}

// ProcessPriorityString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ProcessPriorityString(s string) (ProcessPriority, error) {
	if val, ok := _ProcessPriorityNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _ProcessPriorityNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to ProcessPriority values", s)
}

// ProcessPriorityValues returns all values of the enum
func ProcessPriorityValues() []ProcessPriority {
	return _ProcessPriorityValues
}

// ProcessPriorityStrings returns a slice of all String values of the enum
func ProcessPriorityStrings() []string {
	strs := make([]string, len(_ProcessPriorityNames))
	copy(strs, _ProcessPriorityNames)
	return strs
}

// IsAProcessPriority returns "true" if the value is listed in the enum definition. "false" otherwise
func (i ProcessPriority) IsAProcessPriority() bool {
	for _, v := range _ProcessPriorityValues {
		if i == v {
			return true
		}
	}
	return false
}